package services_test

import (
	"go_cinema_app/internal/repository"
	"go_cinema_app/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	seat []repository.Hall1
}

func (m *mockRepo) GetAll() ([]repository.Hall1, error) {
	return m.seat, nil
}

func (m *mockRepo) ReserveSeatRepo(seat *repository.Hall1, seat1 repository.Hall1) (string, error) {
	m.seat = append(m.seat, *seat)
	return "success_reserved", nil
}

func (m *mockRepo) BuySeatRepo(seat *repository.Hall1, seat1 repository.Hall1) (string, error) {
	m.seat = append(m.seat, *seat)
	return "", nil
}

func (m *mockRepo) CreateMovie(seat *repository.Hall1) error {
	return nil
}

func (m *mockRepo) RemoveMovie() error {
	return nil
}

func (m *mockRepo) FindPriceRepo(seat repository.Hall1) (string, error) {
	return "", nil
}

func TestSeatService(t *testing.T) {
	mock := &mockRepo{}
	svc := services.NewSeatService(mock)

	// Test Create
	u := &repository.Hall1{Seat: 1, Status: "reserved", Movie: "Matrix", Time: "10:00", User: "b"}
	status, err := svc.ReserveSeatSvc(u)
	assert.NoError(t, err)

	// Test Get
	seat, err := svc.GetSeats()
	assert.NoError(t, err)
	assert.Len(t, seat, 1)
	assert.Equal(t, "success_reserved", status)
	assert.Equal(t, 1, seat[0].Seat)
}
