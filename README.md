# whatsapp-api-webhook-server-golang

whatsapp-api-webhook-server-golang - библиотека для интеграции с мессенджером WhatsApp через API
сервиса [green-api.com](https://green-api.com/). Чтобы воспользоваться библиотекой, нужно получить регистрационный токен
и ID аккаунта в [личном кабинете](https://console.green-api.com/). Есть бесплатный тариф аккаунта разработчика.

## API

Документация к REST API находится по [ссылке](https://green-api.com/docs/api/). Библиотека является оберткой к REST API,
поэтому документация по ссылке выше применима и к самой библиотеке.

## Авторизация

Чтобы отправить сообщение или выполнить другие методы Green API, аккаунт WhatsApp в приложении телефона должен быть в
авторизованном состоянии. Для авторизации аккаунта перейдите в [личный кабинет](https://console.green-api.com/) и
сканируйте QR-код с использованием приложения WhatsApp.

## Пример подготовки среды для Ubuntu Server

### Установка Go

На сервере должен быть установлен Go. [Инструкция по установке Go](https://go.dev/doc/install).

### Обновление системы

Обновим систему:

```shell
sudo apt update
sudo apt upgrade -y
```

### Брандмауэр

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

## Как запустить веб-сервер

### Установка

Не забудьте создать модуль:

```shell
go mod init example
```

Установка:

```shell
go get github.com/green-api/whatsapp-api-webhook-server-golang
```

### Импорт

```
import (
	"github.com/green-api/whatsapp-api-webhook-server-golang/pkg"
)
```

### Примеры

#### Как инициализировать объект

Атрибут WebhookToken является опциональным.

```
webhook := pkg.Webhook{
    Address:      ":80",
    Pattern:      "/",
}
```

#### Как запустить веб-сервер

Функция StartServer принимает функцию-обработчик. Функция-обработчик должна содержать 1
параметр (`body map[string]interface{}`). При получении нового уведомления ваша функция-обработчик будет выполнена.

Ссылка на пример: [main.go](examples/main.go).

```
_ := webhook.StartServer(func(body map[string]interface{}) {
    fmt.Println(body)
})
```

### Запуск приложения

```shell
go run main.go
```

## Документация по методам сервиса

[Документация по методам сервиса](https://green-api.com/docs/api/)

## Лицензия

Лицензировано на условиях [
Attribution-NoDerivatives 4.0 International (CC BY-ND 4.0)
](https://creativecommons.org/licenses/by-nd/4.0/).
Смотрите файл [LICENSE](LICENSE).
