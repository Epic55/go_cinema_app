package services

import (
	"go_cinema_app/internal/repository"
	"math/rand"
)

type SeatService interface {
	GetSeats() ([]repository.Hall1, error)
	ReserveSeatSvc(seat *repository.Hall1) (string, error)
	CheckBeforeBuySvc(seat repository.Hall1) (string, error)
	FindPriceSvc(seat repository.Hall1) (string, error)
	BuySeatSvc(seat *repository.Hall1) (string, error)
	CreateMovie(seat *repository.Hall1) error
	RemoveMovie() error
}

type seatService struct {
	repo repository.SeatRepository
}

func NewSeatService(r repository.SeatRepository) SeatService {
	return &seatService{repo: r}
}

func (s *seatService) GetSeats() ([]repository.Hall1, error) {
	return s.repo.GetAll()
}

func (s *seatService) ReserveSeatSvc(seat *repository.Hall1) (string, error) {
	////SIMULATE USERS TRYING TO BOOK THE SAME SEAT
	// var wg sync.WaitGroup
	// wg.Add(3)
	// for i := 1; i <= 2; i++ {
	// 	go func() {
	// 		defer wg.Done()
	// 		s.repo.ReserveSeatRepo(seat, *seat)
	// 	}()
	// }
	// wg.Wait()
	// return repository.StatusUnknownStatus, nil
	return s.repo.ReserveSeatRepo(seat, *seat)
}

func (s *seatService) CheckBeforeBuySvc(seat repository.Hall1) (string, error) {
	return s.repo.CheckBeforeBuyRepo(seat, seat)
}

func (s *seatService) BuySeatSvc(seat *repository.Hall1) (string, error) {
	return s.repo.BuySeatRepo(seat, *seat)
}

func (s *seatService) CreateMovie(seat *repository.Hall1) error {
	movies := []string{seat.Movie}
	time := []string{seat.Time}
	statuses := []string{"available"}
	prices := []int{seat.Price}

	for i := 1; i <= seat.Seat; i++ {
		h := repository.Hall1{
			Seat:   i,
			Status: statuses[rand.Intn(len(statuses))],
			Movie:  movies[rand.Intn(len(movies))],
			Time:   time[rand.Intn(len(time))],
			Price:  prices[rand.Intn(len(prices))],
		}
		s.repo.CreateMovie(&h)
	}
	return nil
}

func (s *seatService) RemoveMovie() error {
	return s.repo.RemoveMovie()
}

func (s *seatService) FindPriceSvc(seat repository.Hall1) (string, error) {
	return s.repo.FindPriceRepo(seat)
}
