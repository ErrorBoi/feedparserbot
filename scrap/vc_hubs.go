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
	vcSubsLink   = "https://vc.ru/subs"
	vcSourceLink = "https://vc.ru/rss/all"
)

func VCHubs(cli *ent.Client) error {
	collector := colly.NewCollector()

	fmt.Println("starting VC hubs scraping")
	hubs, err := getVCHubs(collector)
	if err != nil {
		return err
	}
	fmt.Println("HUBS ACHTUNG:")
	fmt.Println(hubs)

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

				fmt.Printf("Link: %s. Title: %s\n", element.Attr("href"), element.ChildText("span"))

				//if strings.TrimSpace(element.Text) != "" {
				//	//ctx, cancel := context.WithTimeout(context.Background(), utils.WriteTimeout)
				//	//defer cancel()
				//
				//
				//
				//	fmt.Printf("HTML Element's text is %s\n", element.Text)
				//}
			})
		fmt.Println("Hubs counter: ", counter)
		//colly.NewHTMLElementFromSelectionNode()
	})

	err := collector.Visit(vcSubsLink)
	if err != nil {
		return nil, err
	}

	return hubs, nil
}

func storeVCHubs(cli *ent.Client, hubs []string) error {
	fmt.Println("Starting storing VC hubs")
	fmt.Printf("We have %d hubs to be checked\n", len(hubs))

	//ctx, cancel := context.WithTimeout(context.Background(), utils.WriteTimeout)
	//defer cancel()

	VCSource, err := cli.Source.Query().Where(source.URL(vcSourceLink)).Only(context.Background())
	if err != nil {
		return err
	}

	counter := 0

	for _, hub := range hubs {
		collector := colly.NewCollector()

		link := "https://vc.ru/rss/" + hub
		fmt.Println("Visiting hub ", link)

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

			//ctx, cancel := context.WithTimeout(context.Background(), utils.WriteTimeout)
			//defer cancel()

			_, err = cli.Source.Create().SetURL(link).SetTitle(title).SetLanguage(source.Language(lang)).SetParent(VCSource).Save(context.Background())
			if err != nil {
				fmt.Println("ERROR CREATING SOURCE: ", err.Error())
			} else {
				fmt.Println("Stored hub ", link)
			}
			counter++
		})
		if err != nil {
			fmt.Println(err)
		}

		err = collector.Visit(link)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(1*time.Second)
	}

	fmt.Println("Stored hubs counter: ", counter)

	return nil
}
