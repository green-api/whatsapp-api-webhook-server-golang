# whatsapp-api-webhook-server-golang

whatsapp-api-webhook-server-golang - библиотека на Go, созданная для интеграции с WhatsApp через API
сервиса [GREEN API](https://green-api.com/). Чтобы начать использовать библиотеку, вам нужно получить ID и token
аккаунта в [личном кабинете](https://console.green-api.com/).

## API

Документация к REST API находится [здесь](https://green-api.com/docs/api/). Библиотека является оберткой к REST API,
поэтому документация по ссылке выше применима и к самой библиотеке.

## Подготовка среды

На сервере должен быть установлен Go. Установить Go можно так:

```shell
snap install go --classic
```

Проверьте, что вы установили Go:

```shell
go version
```

### Пример подготовки среды на Ubuntu Server

Обновим систему:

```shell
sudo apt update
sudo apt upgrade -y
```

Настроим брандмауэр:

Разрешим соединение по SSH:

```shell
sudo ufw allow ssh
```

Базовые правила:

```shell
sudo ufw default deny incoming
sudo ufw default allow outgoing
```

Разрешаем соединения по HTTP и HTTPS:

```shell
sudo ufw allow http
sudo ufw allow https
```

Активируем брандмауэр:

```shell
sudo ufw enable
```

## Как перенаправить входящие уведомления на сервер

Чтобы перенаправить входящие уведомления на сервер, нужно в личном кабинете установить адрес отправки уведомлений (URL).

![](https://raw.githubusercontent.com/green-api/whatsapp-api-webhook-server-python/master/media/ChangeWebhookServerURL.png)

## Установка

```shell
go get github.com/green-api/whatsapp-api-webhook-server-golang
```

### Установка и запуск примера

Установка:

```shell
wget https://raw.githubusercontent.com/green-api/whatsapp-api-webhook-server-golang/master/examples/main.go
```

Запуск:

```shell
go run main.go
```

## Пример

### Импорт

```
import (
	"github.com/green-api/whatsapp-api-webhook-server-golang/pkg"
)
```

### Как инициализировать объект

Атрибут WebhookToken является опциональным.

```
webhook := pkg.Webhook{
    Address:      ":80",
    Pattern:      "/",
}
```

### Запуск сервера

Функция StartServer принимает функцию-обработчик. Функция-обработчик должна содержать 1
параметр (`body map[string]interface{}`). При получении нового уведомления ваша функция-обработчик будет выполнена.

```
_ := webhook.StartServer(func(body map[string]interface{}) {
    fmt.Println(body)
})
```

## Лицензия

Лицензия MIT. [LICENSE](LICENSE)
