package services

import (
	"math/rand"
	"microservice_template/internal/repository"
)

type SeatService interface {
	GetSeats() ([]repository.Hall1, error)
	BuySeat(seat *repository.Hall1) error
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

func (s *seatService) BuySeat(seat *repository.Hall1) error {
	return s.repo.BuySeat(seat)
}

func (s *seatService) CreateMovie(seat *repository.Hall1) error {
	movies := []string{"Matrix"}
	times := []string{"10:00"}
	statuses := []string{"available"}

	for i := 1; i <= 10; i++ {
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
