package scrap

import (
	"context"
	"fmt"

	"github.com/gocolly/colly"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/source"
	"github.com/ErrorBoi/feedparserbot/utils"
)

var sources = []string{
	"https://vc.ru/rss/all",
	"https://rb.ru/feeds/all",
	"https://www.forbes.ru/newrss.xml",
	"https://www.fontanka.ru/fontanka.rss",
}

func Sources(cli *ent.Client) error {
	collector := colly.NewCollector()

	fmt.Println("starting sources scraping")

	var err error
	for _, feed := range sources {
		collector.OnXML("/rss/channel", func(e *colly.XMLElement) {
			title := e.ChildText("/title")
			lang := e.ChildText("/language")
			switch lang {
			case "ru-RU", "RU":
				lang = "ru"
			case "en-EN", "EN":
				lang = "en"
			default:
				lang = "ru"
			}

			ctx, cancel := context.WithTimeout(context.Background(), utils.WriteTimeout)
			defer cancel()

			_, err = cli.Source.Create().SetURL(feed).SetTitle(title).SetLanguage(source.Language(lang)).Save(ctx)
		})
		if err != nil {
			return err
		}

		err = collector.Visit(feed)
		if err != nil {
			return err
		}
	}

	fmt.Println("finished sources scraping")
	return nil
}
