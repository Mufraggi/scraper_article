package list_announce

import (
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func GetListAnnounce(url string) func(url string) []domain.AnnounceCard {
	return func(url string) []domain.AnnounceCard {
		res, err := http.Get("https://immobilier.lefigaro.fr/annonces/immobilier-vente-appartement-nice+06000.html")
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
		// Find the review items
		var announces []domain.AnnounceCard
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
		return announces
	}
}
