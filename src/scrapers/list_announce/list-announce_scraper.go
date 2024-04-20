package list_announce

import (
	"fmt"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func GetListAnnounce() func(url string) (domain.ListAnnounceCard, error) {
	return func(url string) (domain.ListAnnounceCard, error) {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}
		// Find the review items
		var announces domain.ListAnnounceCard
		doc.Find(".classified-card__content").Each(func(i int, s *goquery.Selection) {
			price := strings.TrimSpace(s.Find(".content__price").Text())
			link, exists := s.Find(".content__link").Attr("href")
			if !exists {
				log.Println("Lien non trouvé")
				return
			}
			description := strings.TrimSpace(s.Find(".content__description").Text())
			squareMeters := strings.TrimSpace(s.Find(".title__options span:contains('m²')").Text())
			rooms := strings.TrimSpace(s.Find(".title__options span:contains('pièces')").Text())
			estateType := strings.TrimSpace(s.Find(".title__estate-type").Text())

			announce := domain.AnnounceCard{
				ProductType:  estateType,
				Price:        price,
				SquareFetter: squareMeters,
				BeadRoom:     rooms,
				URL:          link,
				Description:  description,
			}
			announces = append(announces, announce)
		})
		return announces, nil
	}
}
