package main

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	ctx := context.Background()

	config, err := getJwtConfig()
	if err != nil {
		log.Printf("[ERROR] Failed to process get jwt config: %s", err)
		return
	}

	service, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		log.Printf("[ERROR] Failed to process create service: %s", err)
		return
	}

	message := createMessage()
	_, err = service.Users.Messages.Send("hoge@gmal.com", &message).Do()
	if err != nil {
		log.Printf("[ERROR] Failed to process send message: %s", err)
		return
	}
}

func getJwtConfig() (config *jwt.Config, err error) {
	json, err := ioutil.ReadFile("./credential.json")
	if err != nil {
		log.Printf("[ERROR] Failed to process read file: %s", err)
		return nil, err
	}

	config, err = google.JWTConfigFromJSON(json, gmail.MailGoogleComScope)
	if err != nil {
		log.Printf("[ERROR] Failed to process get jwt config: %s", err)
		return nil, err
	}
	config.Subject = "hoge@gmal.com"

	return config, nil
}

func createMessage() (message gmail.Message) {
	temp := []byte("From: 'hoge@gmal.com'\r\n" +
		"To: fuga@gmail.com\r\n" +
		"Subject: testSubject\r\n" +
		"\r\ntestBody")
	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)

	return message
}
