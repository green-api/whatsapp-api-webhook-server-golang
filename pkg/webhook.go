package pkg

import (
	"encoding/base64"
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
		defer request.Body.Close()

		if w.WebhookToken != "" {
			if !w.validateToken(request.Header.Get("Authorization")) {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		body, err := io.ReadAll(request.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if !json.Valid(body) {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		var data map[string]interface{}

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("Error unmarshaling JSON: %v", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		handler(data)

		writer.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(w.Address, nil)
}

func (w Webhook) validateToken(authHeader string) bool {
	if authHeader == "" {
		return false
	}

	if strings.HasPrefix(authHeader, "Basic ") {
		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return false
		}
		return string(decoded) == w.WebhookToken
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token == w.WebhookToken
}
