package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/gocolly/colly"

	"github.com/ErrorBoi/feedparserbot/db"
)

var VC = "https://vc.ru/rss/all"

var sources = []string{
	"https://vc.ru/rss/all",
	"https://rb.ru/feeds/all",
	"https://www.forbes.ru/newrss.xml",
	"https://www.fontanka.ru/fontanka.rss",
}

func (b *Bot) parseSources() {
	now := time.Now().Add(-60 * time.Minute)

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(VC)
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

			var h1 string
			collector.OnHTML("h1", func(el *colly.HTMLElement) {
				h1 = strings.TrimSpace(el.Text)
				h1 = strings.TrimSuffix(h1, "Материал редакции")
			})

			var contents []string
			collector.OnHTML(".content--full", func(el *colly.HTMLElement) {
				el.ForEach(".layout--a", func(i int, innerEl *colly.HTMLElement) {
					contents = append(contents, strings.TrimSpace(innerEl.Text))
				})
			})

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
				SourceURL:   VC,
			})
			if err != nil {
				b.lg.Errorf("Store Post error: %v", err)
			}
		}
		fmt.Printf("%+v\n", item)
	}
}
