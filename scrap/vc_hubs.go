package scrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/source"
)

var (
	vcHubsLink   = "https://vc.ru/subs"
	vcSourceLink = "https://vc.ru/rss/all"
)

func VCHubs(cli *ent.Client) error {
	collector := colly.NewCollector()

	fmt.Println("starting VC hubs scraping")
	hubs, err := getVCHubs(collector)
	if err != nil {
		return err
	}

	err = storeVCHubs(cli, hubs)

	fmt.Println("finishing VC hubs scraping")
	return nil
}

func getVCHubs(collector *colly.Collector) ([]string, error) {
	var hubs []string

	collector.OnHTML(".subsites_catalog__subscribed", func(el *colly.HTMLElement) {
		counter := 0

		el.ForEach(".subsites_catalog_item>.subsite_card_simple_short>.subsite_card_simple_short__title",
			func(i int, element *colly.HTMLElement) {
				counter++

				href := element.Attr("href")
				if !(strings.HasPrefix(href, "https://vc.ru/u/") || strings.HasPrefix(href, "https://vc.ru/s/")) {
					hubs = append(hubs, strings.TrimPrefix(href, "https://vc.ru/"))
				}
			})
		fmt.Println("Hubs counter: ", counter)
	})

	err := collector.Visit(vcHubsLink)
	if err != nil {
		return nil, err
	}

	return hubs, nil
}

func storeVCHubs(cli *ent.Client, hubs []string) error {
	fmt.Println("Starting storing VC hubs")
	fmt.Printf("We have %d hubs to be checked\n", len(hubs))

	VCSource, err := cli.Source.Query().Where(source.URL(vcSourceLink)).Only(context.Background())
	if err != nil {
		return err
	}

	counter := 0

	for _, hub := range hubs {
		collector := colly.NewCollector()

		link := "https://vc.ru/rss/" + hub

		collector.OnXML("/rss/channel", func(e *colly.XMLElement) {
			title := e.ChildText("/title")
			lang := e.ChildText("/language")
			switch lang {
			case "ru-RU", "RU":
				lang = "RU"
			case "en-EN", "EN":
				lang = "EN"
			default:
				lang = "RU"
			}

			_, err = cli.Source.Create().SetURL(link).SetTitle(title).SetLanguage(source.Language(lang)).SetParent(VCSource).Save(context.Background())
			if err != nil {
				fmt.Println("Error storing source: ", err.Error())
			}
			counter++
		})

		collector.Visit(link)

		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("Stored hubs counter: ", counter)

	return nil
}
