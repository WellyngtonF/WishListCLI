package item

import (
	"time"
)

type Item struct {
	Name            string
	Category        string
	Producer        string
	MaxPrice        float64
	ScrapingSources []string
	URL             string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	MinPrice        float64
}
