---
name: greenapi-webhook-server-golang
version: 1.0.0
description: Receive GREEN-API WhatsApp webhooks in Go with the official library whatsapp-api-webhook-server-golang (HTTP endpoint that accepts incoming notification POSTs). Use when the task is to receive/handle GREEN-API webhooks in Go instead of polling.
license: MIT
metadata:
  sdk: github.com/green-api/whatsapp-api-webhook-server-golang
  docs: https://green-api.com/en/docs/api/receiving/technology-webhook-endpoint/
---

# GREEN-API webhook server for Go

## When to apply

Use when incoming GREEN-API notifications should be **pushed** to your server
(webhook technology) rather than pulled by polling. For polling and for calling
API methods (sending messages etc.) use the client SDK
`github.com/green-api/whatsapp-api-client-golang-v2` — this library only
receives.

## How GREEN-API delivers webhooks

Docs: https://green-api.com/en/docs/api/receiving/technology-webhook-endpoint/

- GREEN-API sends each notification as an HTTP **POST** with a JSON body to the
  `webhookUrl` configured in instance settings (console or `SetSettings`).
- Your endpoint must be publicly reachable and answer **200 OK**; failed
  deliveries are retried at 1-minute intervals, guaranteed within 24 hours.
- If `webhookUrlToken` is set, GREEN-API sends it in the `Authorization`
  header (`Bearer <token>` by default; `Basic <token>` also supported).
- Notification body format and `typeWebhook` values:
  https://green-api.com/en/docs/api/receiving/notifications-format/

## Install

```shell
go mod init myproject
go get github.com/green-api/whatsapp-api-webhook-server-golang
```

## The entire API surface

One type, one method (package `pkg`):

```go
import "github.com/green-api/whatsapp-api-webhook-server-golang/pkg"

webhook := pkg.Webhook{
	Address:      ":5000",        // address for http.ListenAndServe
	Pattern:      "/",            // URL path to handle
	WebhookToken: "your-token",   // optional; must equal webhookUrlToken from instance settings
}

err := webhook.StartServer(func(body map[string]interface{}) {
	// body is the parsed notification JSON
	fmt.Println(body["typeWebhook"])
})
if err != nil {
	log.Fatal(err) // StartServer blocks; returns only on server error
}
```

Behavior implemented by the library (do not re-implement):

- Validates `Authorization` header against `WebhookToken` when set — accepts
  `Bearer <token>` and `Basic <base64(token)>`; responds 401 on mismatch.
- Responds 400 on invalid JSON, 200 after your handler returns.
- The handler receives the notification as `map[string]interface{}` — there
  are no typed notification structs in this library. Read fields via type
  assertions (`body["typeWebhook"].(string)`) or re-marshal into your own
  structs.

## Handler example: react to incoming text messages

```go
err := webhook.StartServer(func(body map[string]interface{}) {
	typeWebhook, _ := body["typeWebhook"].(string)
	if typeWebhook != "incomingMessageReceived" {
		return
	}
	senderData, _ := body["senderData"].(map[string]interface{})
	chatId, _ := senderData["chatId"].(string)

	messageData, _ := body["messageData"].(map[string]interface{})
	typeMessage, _ := messageData["typeMessage"].(string)

	var text string
	switch typeMessage {
	case "textMessage":
		if d, ok := messageData["textMessageData"].(map[string]interface{}); ok {
			text, _ = d["textMessage"].(string)
		}
	case "extendedTextMessage":
		if d, ok := messageData["extendedTextMessageData"].(map[string]interface{}); ok {
			text, _ = d["text"].(string)
		}
	}
	log.Printf("message from %s: %s", chatId, text)
})
```

## Pitfalls

- **Enable webhooks in instance settings**: set `webhookUrl` to your public
  endpoint and turn on the notification types you need (`incomingWebhook`,
  `outgoingWebhook`, `stateWebhook`, …) — via console.green-api.com or
  `Account().SetSettings(...)` in the client SDK. With empty `webhookUrl`
  nothing will ever arrive.
- **Webhooks and polling are alternatives**: while `webhookUrl` is set,
  `ReceiveNotification` polling won't get the notifications.
- `StartServer` blocks the goroutine and uses `http.HandleFunc` on the default
  mux — run it once per process.
- Keep the handler fast: GREEN-API just needs the 200; offload slow work to a
  goroutine/queue.
- The server is plain HTTP. For HTTPS put it behind a reverse proxy
  (nginx/caddy) or a tunnel; `webhookUrl` should be the public URL.
