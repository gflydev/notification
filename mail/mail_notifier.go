package mail

import (
	"github.com/gflydev/core/log"
	"github.com/gflydev/mail"
	"github.com/gflydev/notification"
)

// ========================================================================================
//                            Register Mail Notification Handler
// ========================================================================================

// AutoRegister wires the mail channel into the notification dispatcher so that any
// notification implementing IMailNotification is delivered by email.
func AutoRegister() {
	notification.Register(newMailHandler, (*IMailNotification)(nil))
}

func newMailHandler(notification any) notification.INotifiable {
	return &mailNotification{
		Data: notification.(IMailNotification),
	}
}

// ========================================================================================
//                                 Mail Notification Handler
// ========================================================================================

type mailNotification struct {
	Data IMailNotification
}

// buildEnvelop maps a notification Data payload onto a mail.Envelop, populating the
// optional Cc, Bcc and ReplyTo fields only when they are provided.
func buildEnvelop(data Data) mail.Envelop {
	envelop := mail.Envelop{
		To:      []string{data.To},
		Subject: data.Subject,
		HTML:    data.Body,
	}

	if data.Cc != "" {
		envelop.Cc = []string{data.Cc}
	}

	if data.Bcc != "" {
		envelop.Bcc = []string{data.Bcc}
	}

	if data.ReplyTo != "" {
		envelop.ReplyTo = []string{data.ReplyTo}
	}

	return envelop
}

func (h *mailNotification) Notify() {
	data := h.Data.ToEmail()

	envelop := buildEnvelop(data)

	log.Tracef("Send via Mail data %v", data)

	mail.Send(envelop)
}
