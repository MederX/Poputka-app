package firebase

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var Client *firestore.Client

func Init() {
	ctx := context.Background()
	projectID := os.Getenv("FIREBASE_PROJECT_ID")

	var app *firebase.App
	var err error

	// Support credentials from file or env var (for Heroku)
	if creds := os.Getenv("GOOGLE_CREDENTIALS"); creds != "" {
		conf := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, conf, option.WithCredentialsJSON([]byte(creds)))
	} else if credsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); credsFile != "" {
		conf := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, conf, option.WithCredentialsFile(credsFile))
	} else {
		conf := &firebase.Config{ProjectID: projectID}
		app, err = firebase.NewApp(ctx, conf)
	}
	if err != nil {
		log.Fatalf("firebase init: %v", err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("firestore init: %v", err)
	}
}
