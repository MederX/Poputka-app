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

func ListCargo(c *gin.Context) {
	ctx := context.Background()
	q := firebase.Client.Collection("cargo_requests").Where("status", "==", "open")

	if from := c.Query("from"); from != "" {
		q = q.Where("from_city", "==", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("to_city", "==", to)
	}

	docs, err := q.OrderBy("created_at", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		utils.Err(c, http.StatusInternalServerError, "query failed")
		return
	}

	items := make([]map[string]interface{}, 0, len(docs))
	for _, d := range docs {
		var item map[string]interface{}
		d.DataTo(&item)
		items = append(items, item)
	}
	utils.OK(c, items)
}

func CreateCargo(c *gin.Context) {
	userID := c.GetString("user_id")
	var body struct {
		FromCity    string  `json:"from_city" binding:"required"`
		ToCity      string  `json:"to_city" binding:"required"`
		Description string  `json:"description" binding:"required"`
		WeightKg    float64 `json:"weight_kg" binding:"required,min=0"`
		Notes       string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	snap, _ := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	var user models.User
	snap.DataTo(&user)

	ref := firebase.Client.Collection("cargo_requests").NewDoc()
	cargo := models.CargoRequest{
		ID:          ref.ID,
		UserID:      userID,
		UserName:    user.FirstName + " " + user.LastName,
		FromCity:    body.FromCity,
		ToCity:      body.ToCity,
		Description: body.Description,
		WeightKg:    body.WeightKg,
		Notes:       body.Notes,
		Status:      "open",
		CreatedAt:   time.Now(),
	}

	if _, err := ref.Set(ctx, cargo); err != nil {
		utils.Err(c, http.StatusInternalServerError, "create failed")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cargo})
}

func DeleteCargo(c *gin.Context) {
	userID := c.GetString("user_id")
	cargoID := c.Param("id")
	ctx := context.Background()

	snap, err := firebase.Client.Collection("cargo_requests").Doc(cargoID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "cargo not found")
		return
	}
	var cargo models.CargoRequest
	snap.DataTo(&cargo)
	if cargo.UserID != userID {
		utils.Err(c, http.StatusForbidden, "not your request")
		return
	}

	firebase.Client.Collection("cargo_requests").Doc(cargoID).Delete(ctx)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
