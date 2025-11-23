package fcm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"log/slog"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"google.golang.org/api/option"
)

func NewFirebaseApp(cfg *config.FirebaseConfig) *messaging.Client {
	ctx := context.Background()

	decodedKey, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
	if err != nil {
		slog.Error("Error decoding Firebase private key", "error", err)
		log.Fatalf("Error decoding private key: %v", err)
	}
	credentials := map[string]string{
		"type":                        "service_account",
		"project_id":                  cfg.ProjectID,
		"private_key":                 strings.ReplaceAll(string(decodedKey), "\\n", "\n"),
		"client_email":                cfg.ClientEmail,
		"client_id":                   "102206243799981696473",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        "https://www.googleapis.com/robot/v1/metadata/x509/" + cfg.ClientEmail,
		"universe_domain":             "googleapis.com",
	}

	credBytes, err := json.Marshal(credentials)
	if err != nil {
		slog.Error("Error converting Firebase credentials to JSON", "error", err)
		log.Fatalf("Error converting credentials to JSON: %v", err)
	}

	options := option.WithCredentialsJSON(credBytes)

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID}, options)
	if err != nil {
		slog.Error("Error initializing Firebase app", "error", err)
		log.Fatalf("Error initializing Firebase app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		slog.Error("Error obtaining Firebase messaging client", "error", err)
		log.Fatalf("Error obtaining Firebase messaging client: %v", err)
	}

	slog.Info("Firebase messaging client initialized successfully")
	return client
}
