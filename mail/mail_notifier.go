package mail

import (
	"github.com/gflydev/core/log"
	"github.com/gflydev/mail"
	"github.com/gflydev/notification"
)

// ========================================================================================
//                            Register Mail Notification Handler
// ========================================================================================

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

func (h *mailNotification) Notify() {
	data := h.Data.ToEmail()

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
