package mail

// ========================================================================================
//                             Mail Notification Structure
// ========================================================================================

// Data is the payload returned by a notification's ToEmail method describing
// the message to deliver over the mail channel.
type Data struct {
	To      string // Required. Primary recipient address.
	Cc      string // Optional. Carbon-copy address.
	Bcc     string // Optional. Blind carbon-copy address.
	ReplyTo string // Optional. Reply-To address.
	Subject string // Required. Mail subject line.
	Body    string // Required. HTML mail body.
}

// IMailNotification is implemented by notifications that can be delivered by email.
type IMailNotification interface {
	ToEmail() Data
}
