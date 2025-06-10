package transport

import (
	"context"
	"fmt"
	pb "go_cinema_app/buyingGRPC"
	"go_cinema_app/internal/repository"
	"go_cinema_app/internal/services"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SeatHandler struct {
	service services.SeatService
}

func NewSeatHandler(s services.SeatService) *SeatHandler {
	return &SeatHandler{service: s}
}

func (h *SeatHandler) RegisterRoutes(r *gin.Engine) {
	seat := r.Group("/seat")
	{
		seat.GET("", h.GetSeats)
		seat.POST("", h.ReserveSeatTransport)
	}

	buy := r.Group("/buyseat")
	{
		buy.POST("", h.BuySeatTransport)
	}

	createmovie := r.Group("/createmovie")
	{
		createmovie.POST("", h.CreateMovie)
	}

	removemovie := r.Group("/removemovie")
	{
		removemovie.POST("", h.RemoveMovie)
	}
}

func (h *SeatHandler) GetSeats(c *gin.Context) {
	seat, err := h.service.GetSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve seat"})
		return
	}
	c.JSON(http.StatusOK, seat)
}

func (h *SeatHandler) ReserveSeatTransport(c *gin.Context) {
	var seat repository.Hall1
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if seat.User == "" || seat.Movie == "" || seat.Time == "" || seat.Seat == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fulfill all data"})
		return
	}

	status, err := h.service.ReserveSeatSvc(&seat)
	if err != nil {
		// For DB or unknown errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch status {
	case repository.StatusSuccessReserved:
		c.JSON(http.StatusCreated, gin.H{"message": "Seat is reserved successfully", "seat": seat})
	case repository.StatusSeatSold:
		c.JSON(http.StatusConflict, gin.H{"error": "Seat is already sold"})
	case repository.StatusReservedByAnother:
		c.JSON(http.StatusForbidden, gin.H{"error": "Seat is reserved by another user"})
	case repository.StatusSeatNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
	case repository.StatusUnknownStatus:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unknown seat status"})
	case repository.StatusReservedByTheSameUser:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "U have already reserved this seat"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
	}
}

func (h *SeatHandler) BuySeatTransport(c *gin.Context) {
	var seat repository.Hall1
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if seat.User == "" || seat.Movie == "" || seat.Time == "" || seat.Seat == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fulfill all data"})
		return
	}

	status, err := h.service.CheckBeforeBuySvc(seat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch status {
	case repository.StatusReservedByAnother:
		c.JSON(http.StatusForbidden, gin.H{"message": "Seat is reserved by another user"})
		return
	case repository.StatusSeatSold:
		c.JSON(http.StatusConflict, gin.H{"message": "Seat is already sold"})
		return
	case repository.StatusBoughtByTheSameUser:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "U have already bought this seat"})
		return
	case repository.StatusFailureBuy:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "U need to reserve the seat first"})
		return
	case repository.StatusSeatNotFound:
		c.JSON(http.StatusNotFound, gin.H{"message": "Seat not found"})
		return
	case repository.StatusUnknownStatus:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown seat status"})
		return
	case repository.StatusSuccessCanBuy:
		fmt.Println("User can buy seat")
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}

	//GRPC CONNECTION TO BANK APP
	conn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewBuyingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Find the price for the seat
	price, err := h.service.FindPriceSvc(seat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	priceConverted, err := strconv.Atoi(price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pin := int32(seat.Pin)

	resp, err := client.Buying(ctx, &pb.BuyingRequest{UserId: int32(seat.UserID), Pin: pin, Price: int64(priceConverted)})
	if err != nil {
		log.Printf("Error 1: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("%s", resp.GetMessage())

	status, err = h.service.BuySeatSvc(&seat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch status {
	case repository.StatusSuccessBought:
		c.JSON(http.StatusCreated, gin.H{"message": "Seat is bought successfully", "seat": seat})
	case repository.StatusUnknownStatus:
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unknown seat status"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
	}
}

func (h *SeatHandler) CreateMovie(c *gin.Context) {
	var seat repository.Hall1
	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.service.CreateMovie(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create movie"})
		return
	}

	c.JSON(http.StatusCreated, seat)
}

func (h *SeatHandler) RemoveMovie(c *gin.Context) {
	if err := h.service.RemoveMovie(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove movie"})
		return
	}

	c.JSON(http.StatusCreated, "movie removed successfully")
}
