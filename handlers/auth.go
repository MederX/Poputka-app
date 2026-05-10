package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"kyrgyzstan-rideshare/config"
	"kyrgyzstan-rideshare/firebase"
	"kyrgyzstan-rideshare/utils"
)

type telegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
}

func AuthTelegram(c *gin.Context) {
	var body struct {
		InitData string `json:"init_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, "init_data required")
		return
	}

	if !utils.VerifyTelegramAuth(body.InitData, config.C.BotToken) {
		utils.Err(c, http.StatusUnauthorized, "invalid telegram data")
		return
	}

	vals, _ := url.ParseQuery(body.InitData)
	var tgUser telegramUser
	if err := json.Unmarshal([]byte(vals.Get("user")), &tgUser); err != nil {
		utils.Err(c, http.StatusBadRequest, "invalid user data")
		return
	}

	ctx := context.Background()
	userID := fmt.Sprintf("telegram_%d", tgUser.ID)
	ref := firebase.Client.Collection("users").Doc(userID)

	// Check if user exists to preserve role and created_at
	snap, _ := ref.Get(ctx)
	now := time.Now()

	userData := map[string]interface{}{
		"id":            userID,
		"telegram_id":   tgUser.ID,
		"first_name":    tgUser.FirstName,
		"last_name":     tgUser.LastName,
		"username":      tgUser.Username,
		"profile_photo": tgUser.PhotoURL,
		"updated_at":    now,
	}
	if !snap.Exists() {
		userData["role"] = "passenger"
		userData["created_at"] = now
	}

	if _, err := ref.Set(ctx, userData, firestoreMergeAll()); err != nil {
		utils.Err(c, http.StatusInternalServerError, "failed to save user")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     userID,
		"telegram_id": tgUser.ID,
		"exp":         time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenStr, err := token.SignedString([]byte(config.C.JWTSecret))
	if err != nil {
		utils.Err(c, http.StatusInternalServerError, "token error")
		return
	}

	// Fetch full user to return
	snap2, _ := ref.Get(ctx)
	var user map[string]interface{}
	snap2.DataTo(&user)

	c.JSON(http.StatusOK, gin.H{"token": tokenStr, "user": user})
}
