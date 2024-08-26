# gFly Notification - Mail

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

### Usage

Install
```bash
go get -u github.com/gflydev/notification@v1.0.0
go get -u github.com/gflydev/notification/mail@v1.0.1
```

Quick usage `main.go`
```go
import (
    notificationMail "github.com/gflydev/notification/mail"
    "github.com/gflydev/notification"
)

notificationMail.AutoRegister()
```

Notification file
```go
type ResetPassword struct {
}

func (n ResetPassword) ToEmail() notifyMail.Data {
    return notifyMail.Data{
        To:      "vinh@jivecode.com",
        Subject: "Mail title",
        Body:    "Mail body",
    }
}

resetPassword := ResetPassword{}
if err := notification.Send(resetPassword); err != nil {
    log.Error(err)
}
```