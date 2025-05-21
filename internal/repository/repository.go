package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Hall1 struct {
	Seat   int    `json:"seat" gorm:"primaryKey"`
	Status string `json:"status"`
	Movie  string `json:"movie"`
	Time   string `json:"time"`
	User   string `json:"user"`
}

type SeatRepository interface {
	GetAll() ([]Hall1, error)
	BuySeatRepo(seat *Hall1, seat1 Hall1) error
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

func (r *seatRepository) BuySeatRepo(seat *Hall1, seat1 Hall1) error {
	if err := r.db.Where("seat = ?", seat.Seat).First(seat).Error; err != nil {
		return fmt.Errorf("seat not found: %w", err)
	}

	switch seat.Status {
	case "sold":
		return errors.New("seat is sold")
	case "reserved":
		if seat.User != seat1.User {
			return errors.New("seat is reserved by another user")
		}
		return r.db.Model(seat).Updates(seat1).Error
	case "available":
		return r.db.Model(seat).Updates(seat1).Error
	default:
		return fmt.Errorf("unknown seat status: %s", seat.Status)
	}

	// if err := r.db.Find(seat, seat.Seat).Error; err != nil {
	// 	return fmt.Errorf("seat not found: %w", err)
	// }

	// switch {
	// case seat.Status == "sold":
	// 	return errors.New("seat is sold")
	// case seat.Status == "reserved" && seat.User != seat1.User:
	// 	return errors.New("seat is reserved")
	// case seat.Status == "reserved" && seat.User == seat1.User:
	// 	return r.db.Model(&seat1).Where("seat = ?", seat1.Seat).Updates(seat1).Error
	// case seat.Status == "available":
	// 	return r.db.Model(&seat1).Where("seat = ?", seat1.Seat).Updates(seat1).Error
	// default:
	// 	return fmt.Errorf("unknown status: %s", seat.Status)
	// }
}

func (r *seatRepository) CreateMovie(seat *Hall1) error {
	return r.db.Create(seat).Error
}
