package services_test

import (
	"microservice_template/internal/repository"
	"microservice_template/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	seat []repository.Seat
}

func (m *mockRepo) GetAll() ([]repository.Seat, error) {
	return m.seat, nil
}

func (m *mockRepo) BuySeat(seat *repository.Seat) error {
	m.seat = append(m.seat, *seat)
	return nil
}

func TestSeatService(t *testing.T) {
	mock := &mockRepo{}
	svc := services.NewSeatService(mock)

	// Test Create
	u := &repository.Seat{Name: "Alice", Email: "alice@example.com"}
	err := svc.AddSeat(u)
	assert.NoError(t, err)

	// Test Get
	seat, err := svc.GetSeats()
	assert.NoError(t, err)
	assert.Len(t, seat, 1)
	assert.Equal(t, "Alice", seat[0].Name)
}
