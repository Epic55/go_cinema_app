package repository

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

const (
	StatusSuccessReserved       = "success_reserved"
	StatusSuccessBought         = "success_bought"
	StatusSuccessCanBuy         = "success_can_buy"
	StatusSeatSold              = "seat_sold"
	StatusReservedByAnother     = "reserved_by_another"
	StatusReservedByTheSameUser = "U have already reserved this seat"
	StatusBoughtByTheSameUser   = "U have already bought this seat"
	StatusSeatNotFound          = "not_found"
	StatusUnknownStatus         = "unknown_status"
	StatusDBError               = "db_error"
	StatusFailureBuy            = "U need to reserver the seat first"
)

type Hall1 struct {
	Seat   int    `json:"seat"` //gorm:"primaryKey"`
	Status string `json:"status"`
	Movie  string `json:"movie"`
	Time   string `json:"time"`
	User   string `json:"user,omitempty"`   // User who reserved or bought the seat
	Price  int    `json:"price,omitempty"`  // Price for the seat reservation or purchase
	UserID int    `json:"UserID,omitempty"` // User ID for the seat reservation or purchase
	Pin    int    `json:"Pin,omitempty"`    // Pin for the seat reservation or purchase
}

type SeatRepository interface {
	GetAll() ([]Hall1, error)
	ReserveSeatRepo(seat *Hall1, seat1 Hall1) (string, error)
	CheckBeforeBuyRepo(seat Hall1, seat1 Hall1) (string, error)
	FindPriceRepo(seat Hall1) (string, error)
	BuySeatRepo(seat *Hall1, seat1 Hall1) (string, error)
	CreateMovie(seat *Hall1) error
	RemoveMovie() error
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
	err := r.db.Order("seat asc").Select("seat", "status", "movie", "time").Find(&seat).Error
	return seat, err
}

func (r *seatRepository) ReserveSeatRepo(seat *Hall1, seat1 Hall1) (string, error) {
	if err := r.db.Where("seat = ?", seat.Seat).First(seat).Error; err != nil {
		return StatusSeatNotFound, fmt.Errorf("seat not found: %w", err)
	}

	switch seat.Status {
	case "sold":
		if seat.User == seat1.User {
			return StatusBoughtByTheSameUser, nil
		}

		return StatusSeatSold, nil
	case "reserved":
		if seat.User != seat1.User {
			return StatusReservedByAnother, nil
		}
		//UNCOMMENT THIS TO SIMULATE USERS TRYING TO BOOK THE SAME SEAT
		//fmt.Println("1 - ", StatusReservedByTheSameUser)
		return StatusReservedByTheSameUser, nil
	case "available":
		if err := r.db.Model(seat).Where("seat = ?", seat.Seat).Updates(seat1).Error; err != nil {
			return StatusDBError, err
		}

		return StatusSuccessReserved, nil
	default:
		return StatusUnknownStatus, fmt.Errorf("unknown seat status: %s", seat.Status)
	}
}

func (r *seatRepository) CheckBeforeBuyRepo(seat Hall1, seat1 Hall1) (string, error) {
	fmt.Println("------ CheckBeforeBuyRepo")
	if err := r.db.First(&seat, seat.Seat).Error; err != nil {
		//if err := r.db.Where("seat = ?", seat.Seat).First(seat).Error; err != nil {
		return StatusSeatNotFound, fmt.Errorf("seat not found: %w", err)
	}

	fmt.Println("------ CheckBeforeBuyRepo seat - ", seat)
	fmt.Println("------ CheckBeforeBuyRepo seat1 - ", seat1)
	switch seat.Status {
	case "sold":
		if seat.User == seat1.User {
			fmt.Println("------ CheckBeforeBuyRepo sold  seat.User == seat1.User ")
			return StatusBoughtByTheSameUser, nil
		}
		fmt.Println("------ CheckBeforeBuyRepo sold - ")
		return StatusSeatSold, nil
	case "reserved":
		if seat.User != seat1.User {
			return StatusReservedByAnother, nil
		}

		return StatusSuccessCanBuy, nil
	case "available":
		return StatusFailureBuy, nil
	default:
		return StatusUnknownStatus, fmt.Errorf("unknown seat status: %s", seat.Status)
	}
}

func (r *seatRepository) BuySeatRepo(seat *Hall1, seat1 Hall1) (string, error) {
	if err := r.db.Where("seat = ?", seat.Seat).First(seat).Error; err != nil {
		return StatusSeatNotFound, fmt.Errorf("seat not found: %w", err)
	}

	// switch seat.Status {
	// case "sold":
	// 	if seat.User == seat1.User {
	// 		return StatusBoughtByTheSameUser, nil
	// 	}

	// 	return StatusSeatSold, nil
	// case "reserved":
	// 	if seat.User != seat1.User {
	// 		return StatusReservedByAnother, nil
	// 	}

	if err := r.db.Model(seat).Where("seat = ?", seat.Seat).Omit("user_id", "pin").Updates(seat1).Error; err != nil {
		return StatusDBError, err
	}

	return StatusSuccessBought, nil
	// case "available":
	// 	return StatusFailureBuy, nil
	// default:
	// 	return StatusUnknownStatus, fmt.Errorf("unknown seat status: %s", seat.Status)
	// }
}

func (r *seatRepository) CreateMovie(seat *Hall1) error {
	return r.db.Create(seat).Error
}

func (r *seatRepository) RemoveMovie() error {
	if err := r.db.Exec("TRUNCATE TABLE hall1 RESTART IDENTITY CASCADE").Error; err != nil {
		return err
	}
	return nil
}

func (r *seatRepository) FindPriceRepo(seat Hall1) (string, error) {
	if err := r.db.First(&seat, seat.Seat).Error; err != nil {
		return StatusSeatNotFound, fmt.Errorf("seat not found: %w", err)
	}

	return strconv.Itoa(seat.Price), nil
}
