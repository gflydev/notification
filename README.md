# gFly Notification

    Copyright © 2023, gFly
    https://www.gFly.dev
    All rights reserved.

Channel-agnostic notification dispatcher for the gFly framework. A single
`notification.Send(...)` call fans a payload out to every registered channel
(Mail, SMS, Slack, Database, …) that the payload supports, concurrently.

### Install

```bash
go get -u github.com/gflydev/notification
# Mail channel (separate module)
go get -u github.com/gflydev/notification/mail
```

### Configuration

Sending is gated by the `NOTIFICATION_ENABLE` environment variable. When it is
unset or not truthy, `Send` logs and returns `nil` without dispatching — handy
for local/dev and test environments.

```bash
export NOTIFICATION_ENABLE=true
```

### Quick usage `main.go`

Register the channels you want once at startup:

```go
import (
    notificationMail "github.com/gflydev/notification/mail"
)

func main() {
    notificationMail.AutoRegister()
}
```

### Defining a notification

A notification is any struct that implements a channel interface. Implement
`mail.IMailNotification` (its `ToEmail() mail.Data`) to make it mailable:

```go
import (
    "github.com/gflydev/core/log"
    "github.com/gflydev/notification"
    notifyMail "github.com/gflydev/notification/mail"
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

func send() {
    if err := notification.Send(ResetPassword{}); err != nil {
        log.Error(err)
    }
}
```

`Send` returns:

- `nil` on success, or when notifications are disabled.
- `errors.InvalidParameter` when the notification is `nil`.
- `errors.NotImplemented` when no registered channel supports the payload.

Each channel runs in its own goroutine; a panic in one channel is recovered and
logged so it cannot crash the process or block the other channels.
