# gFly Notification - Mail

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

Mail channel for [`github.com/gflydev/notification`](https://github.com/gflydev/notification).
It plugs the gFly mailer into the notification dispatcher.

### Install

```bash
go get -u github.com/gflydev/notification
go get -u github.com/gflydev/notification/mail
```

### Usage

Register the mail channel once at startup:

```go
import (
    notificationMail "github.com/gflydev/notification/mail"
)

func main() {
    notificationMail.AutoRegister()
}
```

Make a notification mailable by implementing `ToEmail() mail.Data`:

```go
import (
    notifyMail "github.com/gflydev/notification/mail"
    "github.com/gflydev/notification"
)

type ResetPassword struct{}

func (n ResetPassword) ToEmail() notifyMail.Data {
    return notifyMail.Data{
        To:      "vinh@jivecode.com",
        Cc:      "",             // optional
        Bcc:     "",             // optional
        ReplyTo: "",             // optional
        Subject: "Mail title",
        Body:    "Mail body",    // HTML
    }
}

resetPassword := ResetPassword{}
if err := notification.Send(resetPassword); err != nil {
    log.Error(err)
}
```

The `Cc`, `Bcc` and `ReplyTo` fields are optional and only added to the outgoing
envelope when non-empty.
