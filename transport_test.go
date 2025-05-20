package transport_test

import (
	"bytes"
	"encoding/json"
	"microservice_template/internal/repository"
	"microservice_template/internal/transport"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockService struct {
	seat []repository.Seat
}

func (m *mockService) GetSeats() ([]repository.Seat, error) {
	return m.seat, nil
}
func (m *mockService) AddSeat(u *repository.Seat) error {
	m.seat = append(m.seat, *u)
	return nil
}

func TestSeatHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockSvc := &mockService{}
	handler := transport.NewSeatHandler(mockSvc)
	handler.RegisterRoutes(router)

	t.Run("POST /seat", func(t *testing.T) {
		body := `{"name":"Bob","email":"bob@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/seat", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GET /seat", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/seat", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var seat []repository.Seat
		json.NewDecoder(w.Body).Decode(&seat)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, 1, len(seat))
	})
}
