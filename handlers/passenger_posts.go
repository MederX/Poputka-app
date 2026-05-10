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

func ListPassengerPosts(c *gin.Context) {
	ctx := context.Background()
	q := firebase.Client.Collection("passenger_posts").Where("status", "==", "open")

	if from := c.Query("from"); from != "" {
		q = q.Where("from_city", "==", from)
	}
	if to := c.Query("to"); to != "" {
		q = q.Where("to_city", "==", to)
	}

	docs, err := q.OrderBy("desired_date", firestore.Asc).Documents(ctx).GetAll()
	if err != nil {
		utils.Err(c, http.StatusInternalServerError, "query failed")
		return
	}

	posts := make([]map[string]interface{}, 0, len(docs))
	for _, d := range docs {
		var p map[string]interface{}
		d.DataTo(&p)
		posts = append(posts, p)
	}
	utils.OK(c, posts)
}

func CreatePassengerPost(c *gin.Context) {
	userID := c.GetString("user_id")
	var body struct {
		FromCity    string `json:"from_city" binding:"required"`
		ToCity      string `json:"to_city" binding:"required"`
		DesiredDate string `json:"desired_date" binding:"required"`
		SeatsNeeded int    `json:"seats_needed" binding:"required,min=1"`
		Notes       string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.Err(c, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	snap, _ := firebase.Client.Collection("users").Doc(userID).Get(ctx)
	var user models.User
	snap.DataTo(&user)

	ref := firebase.Client.Collection("passenger_posts").NewDoc()
	post := models.PassengerPost{
		ID:            ref.ID,
		PassengerID:   userID,
		PassengerName: user.FirstName + " " + user.LastName,
		FromCity:      body.FromCity,
		ToCity:        body.ToCity,
		DesiredDate:   body.DesiredDate,
		SeatsNeeded:   body.SeatsNeeded,
		Notes:         body.Notes,
		Status:        "open",
		CreatedAt:     time.Now(),
	}

	if _, err := ref.Set(ctx, post); err != nil {
		utils.Err(c, http.StatusInternalServerError, "create failed")
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": post})
}

func DeletePassengerPost(c *gin.Context) {
	userID := c.GetString("user_id")
	postID := c.Param("id")
	ctx := context.Background()

	snap, err := firebase.Client.Collection("passenger_posts").Doc(postID).Get(ctx)
	if err != nil {
		utils.Err(c, http.StatusNotFound, "post not found")
		return
	}
	var post models.PassengerPost
	snap.DataTo(&post)
	if post.PassengerID != userID {
		utils.Err(c, http.StatusForbidden, "not your post")
		return
	}

	firebase.Client.Collection("passenger_posts").Doc(postID).Delete(ctx)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
