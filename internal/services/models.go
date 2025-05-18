package services

import "gorm.io/gorm"

type ShowtimeRequest struct {
	MovieID uint   `json:"movie_id"`
	Start   string `json:"start"`
}

type BookSeatRequest struct {
	UserID     uint   `json:"user_id"`
	ShowtimeID uint   `json:"showtime_id"`
	SeatNumber string `json:"seat_number"`
}

type User struct {
	gorm.Model
	Name string
}

type Movie struct {
	gorm.Model
	Title     string
	Showtimes []Showtime
}

type Showtime struct {
	gorm.Model
	MovieID uint
	Start   string
	Seats   []Seat
}

type Seat struct {
	gorm.Model
	SeatNumber string
	IsBooked   bool
	ShowtimeID uint
	UserID     *uint // Nullable, who booked
}
