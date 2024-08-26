package mail

import (
	"github.com/gflydev/core/log"
	"github.com/gflydev/mail"
	"github.com/gflydev/notification"
)

// ========================================================================================
//                                 Mail Notification Handler
// ========================================================================================

func AutoRegister() {
	notification.Register(notification.Type((*IMailNotification)(nil)), newMailHandler)
}

func newMailHandler(notification any) notification.INotifiable {
	return &mailNotification{
		Notification: notification.(IMailNotification),
	}
}

type mailNotification struct {
	Notification IMailNotification
}

func (h *mailNotification) Notify() {
	data := h.Notification.ToEmail()

	envelop := mail.Envelop{
		To:      []string{data.To},
		Subject: data.Subject,
		HTML:    data.Body,
	}

	if data.Cc != "" {
		envelop.Cc = []string{data.Cc}
	}

	log.Tracef("Send via Mail data %v", data)

	mail.Send(envelop)
}
