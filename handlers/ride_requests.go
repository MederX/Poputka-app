package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kyrgyzstan-rideshare/config"
	"kyrgyzstan-rideshare/firebase"
	"kyrgyzstan-rideshare/models"
	"kyrgyzstan-rideshare/utils"
)

func CreateRideRequest(c *gin.Context) {
	userID := c.GetString("user_id")
	rideID := c.Param("id")

	var body struct {
		SeatsNeeded int    `json:"seats_needed" binding:"required,min=1"`
		Notes       string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()

	// Get ride
	rideSnap, err := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride models.Ride
	rideSnap.DataTo(&ride)

	if ride.SeatsAvailable < body.SeatsNeeded {
		utils.Err(c, http.StatusBadRequest, "not enough seats")
		return
	}

	// Get passenger info
	pasSnap, _ := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	var passenger models.User
	pasSnap.DataTo(&passenger)

	ref := firebase.Client.Collection("ride_requests").NewDoc()
	req := models.RideRequest{
		ID:             ref.ID,
		RideID:         rideID,
		PassengerID:    userID,
		PassengerName:  passenger.FirstName + " " + passenger.LastName,
		PassengerPhoto: passenger.ProfilePhoto,
		SeatsNeeded:    body.SeatsNeeded,
		Notes:          body.Notes,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}

	if _, err := ref.Set(ctx, req); err != nil {
		utils.Err(c, http.StatusInternalServerError, "create failed")
		return
	}

	// Notify driver
	go notifyUser(ctx, ride.DriverID,
		fmt.Sprintf("🚗 Новый запрос на место!\n<b>%s → %s</b>\nПассажир: %s",
			ride.FromCity, ride.ToCity, req.PassengerName))

	c.JSON(http.StatusCreated, gin.H{"data": req})
}

func ListRideRequests(c *gin.Context) {
	userID := c.GetString("user_id")
	rideID := c.Param("id")
	ctx := context.Background()

	// Verify driver owns this ride
	rideSnap, err := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride models.Ride
	rideSnap.DataTo(&ride)
	if ride.DriverID != userID {
		utils.Err(c, http.StatusForbidden, "not your ride")
		return
	}

	docs, err := firebase.Client.Collection("ride_requests").
		Where("ride_id", "==", rideID).Documents(ctx).GetAll()
	if err != nil {
		utils.Err(c, http.StatusInternalServerError, "query failed")
		return
	}

	requests := make([]map[string]interface{}, 0, len(docs))
	for _, d := range docs {
		var r map[string]interface{}
		d.DataTo(&r)
		requests = append(requests, r)
	}
	utils.OK(c, requests)
}

func UpdateRideRequest(c *gin.Context) {
	userID := c.GetString("user_id")
	rideID := c.Param("id")
	reqID := c.Param("reqId")

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, err.Error())
		return
	}
	if body.Status != "accepted" && body.Status != "rejected" {
		utils.Err(c, http.StatusBadRequest, "status must be accepted or rejected")
		return
	}

	ctx := context.Background()

	// Verify driver owns the ride
	rideSnap, err := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride models.Ride
	rideSnap.DataTo(&ride)
	if ride.DriverID != userID {
		utils.Err(c, http.StatusForbidden, "not your ride")
		return
	}

	reqSnap, err := firebase.Client.Collection("ride_requests").Doc(reqID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "request not found")
		return
	}
	var rideReq models.RideRequest
	reqSnap.DataTo(&rideReq)

	firebase.Client.Collection("ride_requests").Doc(reqID).Set(ctx,
		map[string]interface{}{"status": body.Status}, firestoreMergeAll())

	// Update seats if accepted
	if body.Status == "accepted" {
		newSeats := ride.SeatsAvailable - rideReq.SeatsNeeded
		if newSeats < 0 {
			newSeats = 0
		}
		firebase.Client.Collection("rides").Doc(rideID).Set(ctx,
			map[string]interface{}{"seats_available": newSeats}, firestoreMergeAll())
	}

	// Notify passenger
	var msg string
	if body.Status == "accepted" {
		msg = fmt.Sprintf("✅ Водитель принял ваш запрос!\n<b>%s → %s</b>, %s\nСвяжитесь: @%s",
			ride.FromCity, ride.ToCity, ride.DepartureDate, getUsernameFromID(ctx, ride.DriverID))
	} else {
		msg = fmt.Sprintf("❌ Водитель отклонил ваш запрос.\n<b>%s → %s</b>, %s",
			ride.FromCity, ride.ToCity, ride.DepartureDate)
	}

	go notifyUser(ctx, rideReq.PassengerID, msg)

	c.JSON(http.StatusOK, gin.H{"message": body.Status})
}

func getUsernameFromID(ctx context.Context, userID string) string {
	snap, err := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return ""
	}
	var u models.User
	snap.DataTo(&u)
	return u.Username
}

func notifyUser(ctx context.Context, userID string, msg string) {
	snap, err := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return
	}
	var u models.User
	snap.DataTo(&u)
	utils.SendNotification(u.TelegramID, msg, config.C.BotToken)
}
