package sql_repo

import (
	"context"
	"github.com/Mufraggi/scraper_article/src/domain"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
	"strings"
	"time"
)

type AdminRepository struct {
	p *pgxpool.Pool
}

type IAdminRepository interface {
	InsertAnnounce(detail *domain.Detail) (*uuid.UUID, error)
}

// InitAdminRepository initializes a new IAdminRepository with the provided PostgreSQL connection pool.
func InitAdminRepository(p *pgxpool.Pool) IAdminRepository {
	return &AdminRepository{
		p: p,
	}
}
func (a AdminRepository) InsertAnnounce(detail *domain.Detail) (*uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	announce := domain.AnnounceSql{
		AnnouncementId: detail.Id,
		DpeScore:       detail.DpeScore,
		GesScore:       detail.GesScore,
		Characteristic: detail.Characteristic,
		Title:          detail.Title,
		Space:          toInt(detail.Space),
		Rooms:          toInt(detail.Rooms),
		City:           detail.City,
		Price:          priceInt(detail.Price),
		Description:    detail.Description,
	}
	q := `INSERT INTO announcement (announce_mongo_id, dpeScore, gesScore, characteristic, title, space, rooms, city, price, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;`
	var id uuid.UUID
	err := a.p.QueryRow(ctx, q, argsProperty(announce)...).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func priceInt(p string) float64 {
	priceStr := strings.Replace(p, " €", "", -1)
	price, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		return 0
	}
	return float64(price)
}
func argsProperty(p domain.AnnounceSql) []interface{} {
	return []interface{}{
		p.AnnouncementId.String(),
		p.DpeScore,
		p.GesScore,
		p.Characteristic,
		p.Title,
		p.Space,
		p.Rooms,
		p.City,
		p.Price,
		p.Description,
	}
}

func toInt(p string) int {
	priceStr := strings.Replace(p, " €", "", -1)
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		return 0
	}
	return price
}
