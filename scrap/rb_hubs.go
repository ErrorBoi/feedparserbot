package scrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"

	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/source"
)

var (
	rbHubsLink   = "https://rb.ru/list/rss"
	rbSourceLink = "https://rb.ru/feeds/all"
)

func RBHubs(cli *ent.Client) error {
	collector := colly.NewCollector()

	fmt.Println("starting RB hubs scraping")
	hubs, err := getRBHubs(collector)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", hubs)

	err = storeRBHubs(cli, hubs)

	fmt.Println("finishing RB hubs scraping")
	return nil
}

func getRBHubs(collector *colly.Collector) ([]string, error) {
	var hubs []string

	collector.OnHTML(".article__content-block", func(el *colly.HTMLElement) {
		counter := 0

		el.ForEach("li",
			func(i int, element *colly.HTMLElement) {
				counter++
				hubs = append(hubs, strings.TrimSpace(element.Text))
			})
		fmt.Println("Hubs counter: ", counter)
	})

	err := collector.Visit(rbHubsLink)
	if err != nil {
		return nil, err
	}

	return hubs, nil
}

func storeRBHubs(cli *ent.Client, hubs []string) error {
	fmt.Println("Starting storing RB hubs")
	fmt.Printf("We have %d hubs to be checked\n", len(hubs))

	RBSource, err := cli.Source.Query().Where(source.URL(rbSourceLink)).Only(context.Background())
	if err != nil {
		return err
	}

	counter := 0

	for _, hub := range hubs {
		arr := strings.Split(hub, " — ")
		if len(arr) == 1 {
			arr = strings.Split(arr[0], " — ")
			if len(arr) == 1 {
				arr = strings.Split(arr[0], " ")
			}
		}
		title, url := arr[0], arr[1]

		_, err = cli.Source.Create().SetURL(url).SetTitle("Rusbase: "+title).SetLanguage(RBSource.Language).SetParent(RBSource).Save(context.Background())
		if err != nil {
			fmt.Println("Error storing source: ", err.Error())
		}
		counter++
	}

	fmt.Println("Stored hubs counter: ", counter)

	return nil
}
