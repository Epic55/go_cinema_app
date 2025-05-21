package services

import (
	"math/rand"
	"microservice_template/internal/repository"
)

type SeatService interface {
	GetSeats() ([]repository.Hall1, error)
	BuySeatSvc(seat *repository.Hall1) (string, error)
	CreateMovie(seat *repository.Hall1) error
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

func (s *seatService) BuySeatSvc(seat *repository.Hall1) (string, error) {
	return s.repo.BuySeatRepo(seat, *seat)
}

func (s *seatService) CreateMovie(seat *repository.Hall1) error {
	movies := []string{seat.Movie}
	times := []string{seat.Time}
	statuses := []string{"available"}

	for i := 1; i <= seat.Seat; i++ {
		h := repository.Hall1{
			Seat:   i,
			Status: statuses[rand.Intn(len(statuses))],
			Movie:  movies[rand.Intn(len(movies))],
			Time:   times[rand.Intn(len(times))],
		}

		s.repo.CreateMovie(&h)
	}
	return s.repo.CreateMovie(seat)
}
