package announce

import (
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func GetAnnounceDetail() func(url string) domain.Detail {
	return func(url string) domain.Detail {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		d := domain.Detail{}
		doc.Find(".classified-dpe__ges").Each(func(i int, s *goquery.Selection) {
			d.DpeScore = strings.TrimSpace(s.Find(".container-dpe p.pointer").Text())
			d.GesScore = strings.TrimSpace(s.Find(".container-ges p.pointer").Text())
		})
		doc.Find(".features-list li").Each(func(i int, s *goquery.Selection) {
			featureText := strings.TrimSpace(s.Find("span.feature").Text())
			d.Characteristic = append(d.Characteristic, featureText)
		})

		doc.Find("section.classified-main-infos").Each(func(i int, s *goquery.Selection) {
			d.Title = strings.TrimSpace(s.Find("h1").Text())
			parts := strings.Split(d.Title, " ")

			d.Type = parts[1]
			d.Rooms = parts[2]
			d.Space = parts[4]
			d.City = strings.TrimSpace(s.Find("h1 span").Text())
			d.Price = strings.TrimSpace(s.Find(".classified-price-per-m2 strong").Text())
			d.Description = strings.TrimSpace(s.Find(".classified-description p.truncated-description span").Text())
		})
		return d
	}
}
