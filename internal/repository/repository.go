package repository

import (
	"fmt"

	"gorm.io/gorm"
)

const (
	StatusSuccessBought     = "success_bought"
	StatusSuccessReserved   = "success_reserved"
	StatusSeatSold          = "seat_sold"
	StatusReservedByAnother = "reserved_by_another"
	StatusSeatNotFound      = "not_found"
	StatusUnknownStatus     = "unknown_status"
	StatusDBError           = "db_error"
)

type Hall1 struct {
	Seat   int    `json:"seat"` //gorm:"primaryKey"`
	Status string `json:"status"`
	Movie  string `json:"movie"`
	Time   string `json:"time"`
	User   string `json:"user"`
}

type SeatRepository interface {
	GetAll() ([]Hall1, error)
	BuySeatRepo(seat *Hall1, seat1 Hall1) (string, error)
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

func (r *seatRepository) BuySeatRepo(seat *Hall1, seat1 Hall1) (string, error) {
	if err := r.db.Where("seat = ?", seat.Seat).First(seat).Error; err != nil {
		return StatusSeatNotFound, fmt.Errorf("seat not found: %w", err)
	}

	switch seat.Status {
	case "sold":
		return StatusSeatSold, nil
	case "reserved":
		if seat.User != seat1.User {
			return StatusReservedByAnother, nil
		}
		if err := r.db.Model(seat).Updates(seat1).Error; err != nil {
			return StatusDBError, err
		}
		return StatusSuccessBought, nil
	case "available":
		if err := r.db.Model(seat).Updates(seat1).Error; err != nil {
			return StatusDBError, err
		}
		return StatusSuccessReserved, nil
	default:
		return StatusUnknownStatus, fmt.Errorf("unknown seat status: %s", seat.Status)
	}
}

func (r *seatRepository) CreateMovie(seat *Hall1) error {
	// if err := r.db.Exec("TRUNCATE TABLE hall1 RESTART IDENTITY CASCADE").Error; err != nil {
	// 	return err
	// }
	return r.db.Create(seat).Error
}
