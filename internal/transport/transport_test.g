package transport_test

import (
	"bytes"
	"encoding/json"
	"go_cinema_app/internal/repository"
	"go_cinema_app/internal/transport"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	seat []repository.Hall1
}

func (m *mockService) GetSeats() ([]repository.Hall1, error) {
	return m.seat, nil
}
func (m *mockService) ReserveSeatSvc(u *repository.Hall1) (string, error) {
	m.seat = append(m.seat, *u)
	return "", nil
}
func (m *mockService) BuySeatSvc(u *repository.Hall1) (string, error) {
	m.seat = append(m.seat, *u)
	return "", nil
}

func (m *mockService) CreateMovie(u *repository.Hall1) error {
	m.seat = append(m.seat, *u)
	return nil
}

func (m *mockService) RemoveMovie() error {
	return nil
}

func (m *mockService) FindPriceSvc(u repository.Hall1) (string, error) {
	m.seat = append(m.seat, u)
	return "", nil
}

func TestSeatHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockSvc := &mockService{}
	handler := transport.NewSeatHandler(mockSvc)
	handler.RegisterRoutes(router)

	t.Run("GET /seat", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/seat", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var seat []repository.Hall1
		json.NewDecoder(w.Body).Decode(&seat)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1, len(seat))
	})

	t.Run("POST /seat", func(t *testing.T) {
		body := `{"Seat": 10, "Status": "reserved", "Movie": "Matrix", "Time": "10:00", "User": "b"}`
		req := httptest.NewRequest(http.MethodPost, "/seat", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

}
