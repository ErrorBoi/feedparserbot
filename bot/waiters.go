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

func parseVCContent(coll *colly.Collector) (string, []string) {
	var (
		h1       string
		contents []string
	)

	coll.OnHTML("h1", func(el *colly.HTMLElement) {
		h1 = strings.TrimSpace(el.Text)
		h1 = strings.TrimSuffix(h1, "Материал редакции")
	})

	coll.OnHTML(".content--full", func(el *colly.HTMLElement) {
		el.ForEach(".layout--a", func(i int, innerEl *colly.HTMLElement) {
			contents = append(contents, strings.TrimSpace(innerEl.Text))
		})
	})
	return h1, contents
}

func parseRBContent(coll *colly.Collector) {

}
