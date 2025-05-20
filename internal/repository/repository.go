package repository

import (
	"gorm.io/gorm"
)

type Hall1 struct {
	Seat   int    `gorm:"primaryKey"`
	Status string `json:"status"`
	Movie  string `json:"movie"`
	Time   string `json:"time"`
	User   string `json:"user"`
}

type SeatRepository interface {
	GetAll() ([]Hall1, error)
	BuySeat(seat *Hall1) error
	CreateMovie(seat *Hall1) error
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	db.AutoMigrate(&Hall1{})
	return &seatRepository{db: db}
}

func (r *seatRepository) GetAll() ([]Hall1, error) {
	var seat []Hall1
	err := r.db.Find(&seat).Error
	return seat, err
}

func (r *seatRepository) BuySeat(seat *Hall1) error {

	return r.db.Model(&seat).Where("seat = ?", seat.Seat).Updates(seat).Error

}

func (r *seatRepository) CreateMovie(seat *Hall1) error {
	return r.db.Create(seat).Error
}
