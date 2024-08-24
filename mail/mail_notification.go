package mail

// ========================================================================================
//                             Mail Notification Structure
// ========================================================================================

type Data struct {
	To      string
	Cc      string
	Subject string
	Body    string
}

type IMailNotification interface {
	ToEmail() Data
}
