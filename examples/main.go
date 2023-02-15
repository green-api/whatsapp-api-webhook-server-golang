package main

import (
	"fmt"
	"log"

	"github.com/green-api/whatsapp-api-webhook-server-golang/pkg"
)

func main() {
	webhook := pkg.Webhook{
		Address: ":5000",
		Pattern: "/",
	}

	err := webhook.StartServer(func(body map[string]interface{}) {
		fmt.Println(body)
	})
	if err != nil {
		log.Fatal(err)
	}
}
