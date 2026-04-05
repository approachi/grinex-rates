package domain

import "time"

type Rate struct {
	ID        int64
	Ask       int64
	Bid       int64
	CreatedAt time.Time
}

type OrderBookEntry struct {
	Price  int64
	Volume int64
}
