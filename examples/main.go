package main

import (
	"fmt"
	"log"

	"github.com/green-api/whatsapp-api-webhook-server-golang/pkg"
)

func main() {
	webhook := pkg.Webhook{
		Address:      ":80",
		Pattern:      "/",
		WebhookToken: "FV8OtZ8BmXKqM6Fot74D",
	}

	err := webhook.StartServer(func(body map[string]interface{}) {
		fmt.Println(body)
	})
	if err != nil {
		log.Fatal(err)
	}
}
