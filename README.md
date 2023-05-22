# whatsapp-api-webhook-server-golang

- [Документация на русском языке](docs/README_RU.md).

whatsapp-api-webhook-server-golang is a library for integration with WhatsApp messenger using the API
service [green-api.com](https://green-api.com/en/). You should get a registration token and an account ID in
your [personal cabinet](https://console.green-api.com/) to use the library. There is a free developer account tariff.

## API

The documentation for the REST API can be found at the [link](https://green-api.com/en/docs/). The library is a wrapper
for the REST API, so the documentation at the link above also applies.

## Authorization

To send a message or perform other Green API methods, the WhatsApp account in the phone app must be authorized. To
authorize the account, go to your [cabinet](https://console.green-api.com/) and scan the QR code using the WhatsApp app.

## Example of preparing the environment for Ubuntu Server

### Go Installation

Go must be installed on the server. [Go installation instructions](https://go.dev/doc/install).

### Updating the system

Update the system:

```shell
sudo apt update
sudo apt upgrade -y
```

### Firewall

Set up the firewall:

Allow the SSH connection:

```shell
sudo ufw allow ssh
```

Base rules:

```shell
sudo ufw default deny incoming
sudo ufw default allow outgoing
```

Allow HTTP and HTTPS connections:

```shell
sudo ufw allow http
sudo ufw allow https
```

Enable the firewall:

```shell
sudo ufw enable
```

## How to run the web server

### Installation

Do not forget to create a module:

```shell
go mod init example
```

Installation:

```shell
go get github.com/green-api/whatsapp-api-webhook-server-golang
```

### Import

```
import (
	"github.com/green-api/whatsapp-api-webhook-server-golang/pkg"
)
```

### Examples

#### How to initialize an object

The WebhookToken attribute is optional.

```
webhook := pkg.Webhook{
    Address:      ":80",
    Pattern:      "/",
}
```

#### How to run the web server

The StartServer function takes a handler function. The handler function must have 1
parameter (`body map[string]interface{}`). When a new notification is received, your handler function will be executed.

Link to example: [main.go](examples/main.go).

```
_ := webhook.StartServer(func(body map[string]interface{}) {
    fmt.Println(body)
})
```

### Running the application

```shell
go run main.go
```

## Service methods documentation

[Service methods documentation](https://green-api.com/en/docs/api/)

## License

Licensed under [
Attribution-NoDerivatives 4.0 International (CC BY-ND 4.0)
](https://creativecommons.org/licenses/by-nd/4.0/) terms.
Please see file [LICENSE](LICENSE).
