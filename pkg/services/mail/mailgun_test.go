package mail

import (
	"github.com/babafemi99/WR/internal/config"
	"log"
	"testing"
)

func Test_mailgun_SendEmail(t *testing.T) {
	cfg := config.New()
	newMailgun := NewMailgun(cfg)
	data := struct {
		Type     string
		Email    string
		Password string
	}{
		Type:     "Admin",
		Email:    "bayo@bb.com",
		Password: "6geheie9eyhneo",
	}
	patterns := []string{"welcome_user.tmpl"}

	err := newMailgun.SendEmail("ooluwa27@gmail.com", nil, data, patterns...)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Errorf("Error sending email: %v", err)
	}

	log.Println("shutting down")
	err = newMailgun.ShutDown()
	if err != nil {
		return
	}
}
