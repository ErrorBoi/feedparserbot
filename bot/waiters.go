package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"

	"github.com/gocolly/colly"

	"github.com/ErrorBoi/feedparserbot/db"
	"github.com/ErrorBoi/feedparserbot/ent/post"
	"github.com/ErrorBoi/feedparserbot/ent/user"
)

var sources = []string{
	"https://vc.ru/rss/all",
	"https://rb.ru/feeds/all",
	"https://www.forbes.ru/newrss.xml",
	"https://www.fontanka.ru/fontanka.rss",
}

func (b *Bot) parseSources() {
	//TODO: fix time to 5 minutes after testing. Or remove hardcoded time and make it env variable
	now := time.Now().Add(-300 * time.Minute)

	for _, src := range sources {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(src)
		if err != nil {
			b.lg.Errorf("Error parsing URL: %v", err)
		}

		ctx := context.Background()

		items := feed.Items
		for _, item := range items {
			published, err := time.Parse(time.RFC1123Z, item.Published)
			if err != nil {
				b.lg.Errorf("Time parse error: %v", err)
			}
			if published.After(now) {
				collector := colly.NewCollector()

				var (
					h1       string
					contents []string
				)

				// Different parsing rules depending on source
				switch src {
				case "https://vc.ru/rss/all":
					collector.OnHTML("h1", func(el *colly.HTMLElement) {
						h1 = strings.TrimSpace(el.Text)
						h1 = strings.TrimSuffix(h1, "Материал редакции")
					})

					collector.OnHTML(".content--full", func(el *colly.HTMLElement) {
						el.ForEach(".layout--a", func(i int, innerEl *colly.HTMLElement) {
							contents = append(contents, strings.TrimSpace(innerEl.Text))
						})
					})
				case "https://rb.ru/feeds/all":
					collector.OnHTML("h1", func(el *colly.HTMLElement) {
						h1 = strings.TrimSpace(el.Text)
					})

					collector.OnHTML("div[itemtype]", func(el *colly.HTMLElement) {
						el.ForEach(".article__content-block", func(i int, innerEl *colly.HTMLElement) {
							if i == 1 {
								contents = append(contents, strings.TrimSpace(innerEl.Text))
							}
						})
					})
				case "https://www.forbes.ru/newrss.xml":
					collector.OnHTML("h1", func(el *colly.HTMLElement) {
						h1 = strings.TrimSpace(el.Text)
					})

					collector.OnHTML(".article__body", func(el *colly.HTMLElement) {
						el.ForEach(".article__item", func(i int, innerEl *colly.HTMLElement) {
							contents = append(contents, strings.TrimSpace(innerEl.Text))
						})
					})
				case "https://www.fontanka.ru/fontanka.rss":
					collector.OnHTML("h1", func(el *colly.HTMLElement) {
						h1 = strings.TrimSpace(el.Text)
					})

					// calendar.fontanka should be parsed in special way
					if strings.HasPrefix(item.Link, "https://calendar") {
						collector.OnHTML(".content-block", func(el *colly.HTMLElement) {
							el.ForEach("p", func(i int, innerEl *colly.HTMLElement) {
								contents = append(contents, strings.TrimSpace(innerEl.Text))
							})
						})
					}

					collector.OnHTML("article[data-io-article-url]", func(el *colly.HTMLElement) {
						el.ForEach("article[data-io-article-url]>div", func(i int, innerEl *colly.HTMLElement) {
							if i != 0 {
								contents = append(contents, strings.TrimSpace(innerEl.Text))
							}
						})
					})
				}

				err = collector.Visit(item.Link)
				if err != nil {
					b.lg.Errorf("Error visiting URL: %v", err)
				}

				_, err = b.db.StorePost(ctx, db.StorePost{
					Title:       item.Title,
					Url:         item.Link,
					PublishedAt: published,
					Description: item.Description,
					H1:          h1,
					Content:     strings.Join(contents, "\n"),
					SourceURL:   src,
				})
				if err != nil {
					b.lg.Errorf("Store Post error: %v", err)
				}
			}
			fmt.Printf("%+v\n", item)
		}
	}
}

func (b *Bot) sendPostsQuick() {
	ctx := context.Background()

	// Query users with >0 subscriptions
	uu, err := b.db.Cli.User.Query().Where(user.HasSources()).All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying users: %v", err)
	}
	for _, u := range uu {
		ctx := context.Background()

		settings, err := u.QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("failed querying settings: %v", err)
		}

		var frequency time.Duration
		switch settings.SendingFrequency {
		case "instant":
			frequency = 0
		case "1h":
			frequency = 1 * time.Hour
		case "4h":
			frequency = 4 * time.Hour
		}

		if time.Now().After(settings.LastSending.Add(frequency)) {
			posts, err := u.QuerySources().QueryPosts().Where(post.PublishedAtGT(settings.LastSending)).All(ctx)
			if err != nil {
				b.lg.Errorf("failed querying posts: %v", err)
			}

			for _, p := range posts {
				//TODO: send proper message, not just title
				msg := tgbotapi.NewMessage(int64(u.TgID), p.Title)
				msg.ParseMode = tgbotapi.ModeHTML
				b.BotAPI.Send(msg)

				_, err = settings.Update().SetLastSending(time.Now()).Save(ctx)
				if err != nil {
					b.lg.Errorf("failed updating user settings: %v", err)
				}
			}
		}
	}
}