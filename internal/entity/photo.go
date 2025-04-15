package entity

import "time"

type Photo struct {
	ID        int64
	Url       string
	CreatedAt time.Time
}
