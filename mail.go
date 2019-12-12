package main

import (
	"context"
	"encoding/base64"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
)

func main() {

	service, err := createService()
	if err != nil {
		log.Printf("[ERROR] Failed to process create service: %s", err)
		return
	}

	// Fromには偽装したアドレスもしくはそのaliasを指定する
	message := createMessage("送信者", "hoge_alias@gmail.com", "fuga@gmail.com", "テストメール", "テスト")

	// meを指定すれば偽装したアカウントでAPIを叩いてくれる
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

func createMessage(senderName string, senderAddress string, toAddress string, subject string, body string) gmail.Message {
	// body以外に日本語を指定したい場合はエンコードが必要
	fromStr := "From: =?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(senderName)) + "?= <" + senderAddress + ">\n"
	toStr := "To: " + toAddress + "\n"
	subStr := "Subject: =?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?=\n"
	contentType := "Content-Type: text/plain; charset=UTF-8\n"
	bodyStr := "\n" + body

	var message gmail.Message
	// 全体をURLエンコードしてRawに入れる
	message.Raw = base64.URLEncoding.EncodeToString([]byte(fromStr + toStr + subStr + contentType + bodyStr))
	return message
}
