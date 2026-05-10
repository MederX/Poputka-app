package handlers

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"kyrgyzstan-rideshare/firebase"
	"kyrgyzstan-rideshare/models"
	"kyrgyzstan-rideshare/utils"
)

func ListRides(c *gin.Context) {
	ctx := context.Background()
	q := firebase.Client.Collection("rides").Where("status", "==", "active")

	if from := c.Query("from"); from != "" {
		q = q.Where("from_city", "==", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("to_city", "==", to)
	}
	if date := c.Query("date"); date != "" {
		q = q.Where("departure_date", "==", date)
	}

	docs, err := q.OrderBy("departure_date", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		utils.Err(c, http.StatusInternalServerError, "query failed")
		return
	}

	rides := make([]map[string]interface{}, 0, len(docs))
	for _, d := range docs {
		var r map[string]interface{}
		d.DataTo(&r)
		rides = append(rides, r)
	}
	utils.OK(c, rides)
}

func CreateRide(c *gin.Context) {
	userID := c.GetString("user_id")
	var body struct {
		FromCity      string `json:"from_city" binding:"required"`
		ToCity        string `json:"to_city" binding:"required"`
		DepartureDate string `json:"departure_date" binding:"required"`
		DepartureTime string `json:"departure_time" binding:"required"`
		SeatsTotal    int    `json:"seats_total" binding:"required,min=1"`
		PricePerSeat  int    `json:"price_per_seat" binding:"required,min=0"`
		VehicleType   string `json:"vehicle_type"`
		CargoAllowed  bool   `json:"cargo_allowed"`
		Notes         string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	// Get driver info
	snap, err := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusInternalServerError, "user not found")
		return
	}
	var user models.User
	snap.DataTo(&user)

	ref := firebase.Client.Collection("rides").NewDoc()
	ride := models.Ride{
		ID:             ref.ID,
		DriverID:       userID,
		DriverName:     user.FirstName + " " + user.LastName,
		DriverPhoto:    user.ProfilePhoto,
		FromCity:       body.FromCity,
		ToCity:         body.ToCity,
		DepartureDate:  body.DepartureDate,
		DepartureTime:  body.DepartureTime,
		SeatsTotal:     body.SeatsTotal,
		SeatsAvailable: body.SeatsTotal,
		PricePerSeat:   body.PricePerSeat,
		Currency:       "KGS",
		VehicleType:    body.VehicleType,
		CargoAllowed:   body.CargoAllowed,
		Notes:          body.Notes,
		Status:         "active",
		CreatedAt:      time.Now(),
	}

	if _, err := ref.Set(ctx, ride); err != nil {
		utils.Err(c, http.StatusInternalServerError, "create failed")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": ride})
}

func GetRide(c *gin.Context) {
	ctx := context.Background()
	snap, err := firebase.Client.Collection("rides").Doc(c.Param("id")).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride map[string]interface{}
	snap.DataTo(&ride)
	utils.OK(c, ride)
}

func UpdateRide(c *gin.Context) {
	userID := c.GetString("user_id")
	rideID := c.Param("id")
	ctx := context.Background()

	snap, err := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride models.Ride
	snap.DataTo(&ride)
	if ride.DriverID != userID {
		utils.Err(c, http.StatusForbidden, "not your ride")
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, "invalid body")
		return
	}
	// Prevent overwriting protected fields
	delete(body, "id")
	delete(body, "driver_id")
	delete(body, "created_at")

	if _, err := firebase.Client.Collection("rides").Doc(rideID).Set(ctx, body, firestoreMergeAll()); err != nil {
		utils.Err(c, http.StatusInternalServerError, "update failed")
		return
	}
	snap2, _ := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	var updated map[string]interface{}
	snap2.DataTo(&updated)
	utils.OK(c, updated)
}

func DeleteRide(c *gin.Context) {
	userID := c.GetString("user_id")
	rideID := c.Param("id")
	ctx := context.Background()

	snap, err := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride models.Ride
	snap.DataTo(&ride)
	if ride.DriverID != userID {
		utils.Err(c, http.StatusForbidden, "not your ride")
		return
	}

	firebase.Client.Collection("rides").Doc(rideID).Delete(ctx)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func CompleteRide(c *gin.Context) {
	userID := c.GetString("user_id")
	rideID := c.Param("id")
	ctx := context.Background()

	snap, err := firebase.Client.Collection("rides").Doc(rideID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "ride not found")
		return
	}
	var ride models.Ride
	snap.DataTo(&ride)
	if ride.DriverID != userID {
		utils.Err(c, http.StatusForbidden, "not your ride")
		return
	}

	firebase.Client.Collection("rides").Doc(rideID).Set(ctx,
		map[string]interface{}{"status": "completed"}, firestoreMergeAll())
	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}
