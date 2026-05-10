package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"kyrgyzstan-rideshare/firebase"
	"kyrgyzstan-rideshare/utils"
)

func GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	ctx := context.Background()

	snap, err := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "user not found")
		return
	}
	var user map[string]interface{}
	snap.DataTo(&user)
	utils.OK(c, user)
}

func UpdateMe(c *gin.Context) {
	userID := c.GetString("user_id")
	var body struct {
		PhoneNumber string `json:"phone_number"`
		Role        string `json:"role"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, "invalid body")
		return
	}

	updates := map[string]interface{}{"updated_at": time.Now()}
	if body.PhoneNumber != "" {
		updates["phone_number"] = body.PhoneNumber
	}
	if body.Role == "driver" || body.Role == "passenger" {
		updates["role"] = body.Role
	}

	ctx := context.Background()
	if _, err := firebase.Client.Collection("users").Doc(userID).Set(ctx, updates, firestoreMergeAll()); err != nil {
		utils.Err(c, http.StatusInternalServerError, "update failed")
		return
	}

	snap, _ := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	var user map[string]interface{}
	snap.DataTo(&user)
	utils.OK(c, user)
}
