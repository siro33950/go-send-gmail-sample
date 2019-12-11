package main

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"strings"
)

func main() {

	service, err := createService()
	if err != nil {
		log.Printf("[ERROR] Failed to process create service: %s", err)
		return
	}
	// Fromには偽装したアドレスもしくはそのaliasを指定する
	temp := []byte("From: hoge@gmal.com\r\n" +
		"To: fuga@gmail.com\r\n" +
		"Subject: testSubject\r\n" +
		"\r\ntestBody")
	var message gmail.Message
	message.Raw = base64.StdEncoding.EncodeToString(temp)
	message.Raw = strings.Replace(message.Raw, "/", "_", -1)
	message.Raw = strings.Replace(message.Raw, "+", "-", -1)
	message.Raw = strings.Replace(message.Raw, "=", "", -1)
	
	// 偽装したアドレスでAPIを叩く
	_, err = service.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("[ERROR] Failed to process send message: %s", err)
		return
	}
}

func createService() (service *gmail.Service, err error) {
	// サービスアカウント作成時にダウンロードしたJSONを読み込む
	json, err := ioutil.ReadFile("./credential.json")
	if err != nil {
		log.Printf("[ERROR] Failed to process read file: %s", err)
		return nil, err
	}

	// スコープはGSuiteで指定した物をそのまま記述する
	config, err := google.JWTConfigFromJSON(json, gmail.MailGoogleComScope)
	if err != nil {
		log.Printf("[ERROR] Failed to process get jwt config: %s", err)
		return nil, err
	}
	// 認証を偽装するアドレス
	// 指定したアドレスで認証したことになる
	config.Subject = "hoge@gmal.com"

	ctx := context.Background()
	tokenSource := config.TokenSource(ctx)
	return gmail.NewService(ctx, option.WithTokenSource(tokenSource))
}
