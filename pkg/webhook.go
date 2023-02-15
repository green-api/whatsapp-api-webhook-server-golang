package pkg

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Webhook struct {
	Address      string
	Pattern      string
	WebhookToken string
}

func (w Webhook) StartServer(handler func(map[string]interface{})) error {
	http.HandleFunc(w.Pattern, func(writer http.ResponseWriter, request *http.Request) {
		if w.WebhookToken != "" {
			token := request.Header.Get("Authorization")
			token = strings.ReplaceAll(token, "Bearer ", "")
			if token != w.WebhookToken {
				return
			}
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			log.Fatal(err)
		}

		var data map[string]interface{}

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatal(err)
		}

		handler(data)

		writer.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(w.Address, nil)
}
