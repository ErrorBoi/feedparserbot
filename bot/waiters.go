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
	"github.com/ErrorBoi/feedparserbot/ent"
	"github.com/ErrorBoi/feedparserbot/ent/globalsettings"
	"github.com/ErrorBoi/feedparserbot/ent/post"
	"github.com/ErrorBoi/feedparserbot/ent/schema"
	"github.com/ErrorBoi/feedparserbot/ent/user"
	"github.com/ErrorBoi/feedparserbot/ent/usersettings"
	"github.com/ErrorBoi/feedparserbot/utils"
)

var sources = []string{
	"https://vc.ru/rss/all",
	"https://rb.ru/feeds/all/",
	"https://www.forbes.ru/newrss.xml",
	"https://www.fontanka.ru/fontanka.rss",
}

func (b *Bot) parseSources() {
	now := time.Now().Add(-5 * time.Minute)

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
				b.lg.Infof("POST: %+v", item)
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
				case "https://rb.ru/feeds/all/":
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

				p, err := b.db.StorePost(ctx, db.StorePost{
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

				ctx := context.Background()
				pSource, err := p.QuerySource().Only(ctx)
				if err != nil {
					b.lg.Errorf("failed querying source from post: %v", err)
				}

				uu, err := pSource.QueryUsers().
					Where(
						user.HasSettingsWith(
							usersettings.SendingFrequencyNEQ(usersettings.SendingFrequencyInstant),
						),
					).All(ctx)
				for _, u := range uu {
					us, err := u.QuerySettings().Only(ctx)
					if err != nil {
						b.lg.Errorf("failed querying user settings: %v", err)
					}

					hasUrgentWords := false
					for _, urgent := range us.UrgentWords {
						if strings.Contains(item.Title, urgent) {
							hasUrgentWords = true
							break
						}
					}

					hasBannedWords := false
					for _, banned := range us.BannedWords {
						if strings.Contains(item.Title, banned) {
							hasBannedWords = true
							break
						}
					}

					if hasUrgentWords && !hasBannedWords {
						pi := utils.PostInfo{
							SourceTitle: pSource.Title,
							PostTitle:   p.Title,
							URL:         p.URL,
							Description: p.Description,
						}
						msg := tgbotapi.NewMessage(int64(u.TgID), utils.FormatPost(pi))
						msg.ParseMode = tgbotapi.ModeHTML
						b.BotAPI.Send(msg)

						_, err = us.Update().SetLastSending(time.Now()).Save(ctx)
						if err != nil {
							b.lg.Errorf("failed updating user settings: %v", err)
						}
					}
				}

				globalSettings, err := b.db.Cli.Globalsettings.Query().Where(globalsettings.ID(1)).Only(ctx)
				if err != nil {
					b.lg.Errorf("failed querying global settings: %v", err)
				}

				hasClickbaitWords := false
				for _, clickbait := range globalSettings.ClickbaitWords {
					if strings.Contains(item.Title, clickbait) {
						hasClickbaitWords = true
						break
					}
				}

				if hasClickbaitWords {
					editors, err := b.db.Cli.User.Query().Where(user.RoleIn(user.RoleAdmin, user.RoleEditor)).All(ctx)
					if err != nil {
						b.lg.Errorf("failed querying admins and editors: %v", err)
					}

					for _, editor := range editors {
						us, err := editor.QuerySettings().Only(ctx)
						if err != nil {
							b.lg.Errorf("failed querying user settings: %v", err)
						}

						language := string(us.Language)

						var (
							title string
							description string
						)
						switch language {
						case "RU":
							title = p.TitleTranslations.RU
							description = p.DescriptionTranslations.RU
						case "EN":
							title = p.TitleTranslations.EN
							description = p.DescriptionTranslations.EN
						}

						text := fmt.Sprintf(ClickbaitFormatMessage[language], p.ID, p.URL, title, description)


						msg := tgbotapi.NewMessage(int64(editor.TgID), text)
						msg.ParseMode = tgbotapi.ModeHTML
						b.BotAPI.Send(msg)
					}
				}
			}
		}
	}
}

func (b *Bot) sendPostsQuick() {
	ctx := context.Background()

	// Query users with >0 subscriptions
	uu, err := b.db.Cli.User.Query().
		Where(
			user.HasSources(),
			user.HasSettingsWith(
				usersettings.SendingFrequencyIn(
					usersettings.SendingFrequencyInstant,
					usersettings.SendingFrequency1h,
					usersettings.SendingFrequency4h,
				)),
		).
		All(ctx)
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
			b.sendPost(u, settings)
		}
	}
}

func (b *Bot) sendPostsAM() {
	ctx := context.Background()

	// Query users with >0 subscriptions
	uu, err := b.db.Cli.User.Query().
		Where(
			user.HasSources(),
			user.HasSettingsWith(usersettings.SendingFrequencyEQ(usersettings.SendingFrequencyAm)),
		).
		All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying users: %v", err)
	}

	for _, u := range uu {
		ctx := context.Background()

		settings, err := u.QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("failed querying settings: %v", err)
		}

		if time.Now().After(settings.LastSending) {
			b.sendPost(u, settings)
		}
	}
}

func (b *Bot) sendPostsPM() {
	ctx := context.Background()

	// Query users with >0 subscriptions
	uu, err := b.db.Cli.User.Query().
		Where(
			user.HasSources(),
			user.HasSettingsWith(usersettings.SendingFrequencyEQ(usersettings.SendingFrequencyPm)),
		).
		All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying users: %v", err)
	}

	for _, u := range uu {
		ctx := context.Background()

		settings, err := u.QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("failed querying settings: %v", err)
		}

		if time.Now().After(settings.LastSending) {
			b.sendPost(u, settings)
		}
	}
}

func (b *Bot) sendPostsDaily() {
	ctx := context.Background()

	// Query users with >0 subscriptions
	uu, err := b.db.Cli.User.Query().
		Where(
			user.HasSources(),
			user.HasSettingsWith(
				usersettings.SendingFrequencyIn(
					usersettings.SendingFrequencyMon,
					usersettings.SendingFrequencyTue,
					usersettings.SendingFrequencyWed,
					usersettings.SendingFrequencyThu,
					usersettings.SendingFrequencyFri,
					usersettings.SendingFrequencySat,
					usersettings.SendingFrequencySun,
				)),
		).
		All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying users: %v", err)
	}
	for _, u := range uu {
		ctx := context.Background()

		settings, err := u.QuerySettings().Only(ctx)
		if err != nil {
			b.lg.Errorf("failed querying settings: %v", err)
		}

		var weekday time.Weekday
		switch settings.SendingFrequency {
		case "sun":
			weekday = 0
		case "mon":
			weekday = 1
		case "tue":
			weekday = 2
		case "wed":
			weekday = 3
		case "thu":
			weekday = 4
		case "fri":
			weekday = 5
		case "sat":
			weekday = 6
		}

		now := time.Now()
		if now.After(settings.LastSending) && now.Weekday() == weekday {
			b.sendPost(u, settings)
		}
	}
}

func (b *Bot) sendPost(u *ent.User, us *ent.UserSettings) {
	ctx := context.Background()

	posts, err := u.QuerySources().QueryPosts().Where(post.UpdatedAtGT(us.LastSending)).All(ctx)
	if err != nil {
		b.lg.Errorf("failed querying posts: %v", err)
	}

	sentCounter := 0
	for _, p := range posts {
		hasBannedWords := false
		for _, banned := range us.BannedWords {
			if strings.Contains(p.Title, banned) {
				hasBannedWords = true
				break
			}
		}

		if !hasBannedWords {
			gs, err := b.db.Cli.Globalsettings.Query().Where(globalsettings.ID(1)).Only(ctx)
			if err != nil {
				b.lg.Errorf("failed querying global settings: %v", err)
			}


			hasClickbaitWords := false
			for _, clickbait := range gs.ClickbaitWords {
				if strings.Contains(p.Title, clickbait) {
					hasClickbaitWords = true
					break
				}
			}

			if !hasClickbaitWords || p.Subject != nil {
				src, err := p.QuerySource().Only(ctx)
				if err != nil {
					b.lg.Errorf("failed querying source from post: %v", err)
				}

				var titleTranslations schema.Translations
				if hasClickbaitWords {
					titleTranslations = p.SubjectTranslations
				} else {
					titleTranslations = p.TitleTranslations
				}

				var (
					postTitle string
					description string
				)
				switch us.Language {
				case usersettings.LanguageEN:
					postTitle = titleTranslations.EN
					description = p.DescriptionTranslations.EN
				case usersettings.LanguageRU:
					postTitle = titleTranslations.RU
					description = p.DescriptionTranslations.RU
				}



				pi := utils.PostInfo{
					SourceTitle: src.Title,
					PostTitle:   postTitle,
					URL:         p.URL,
					Description: description,
				}
				msg := tgbotapi.NewMessage(int64(u.TgID), utils.FormatPost(pi))
				msg.ParseMode = tgbotapi.ModeHTML
				b.BotAPI.Send(msg)

				sentCounter++
			}
		}
	}

	if sentCounter > 0 {
		_, err = us.Update().SetLastSending(time.Now()).Save(ctx)
		if err != nil {
			b.lg.Errorf("failed updating user settings: %v", err)
		}
	}
}