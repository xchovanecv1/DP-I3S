package email

// ContentAdminAuth defines content for login token email template.
type ContentAdminAuth struct {
	Email    string
	Name     string
	UserName string
	RoomName string
	ItemName string
}

// AdminAuthNotif creates and sends a login token email with provided template content.
func (m *Mailer) AdminAuthNotif(name, address string, content ContentAdminAuth) error {
	msg := &message{
		from:     m.from,
		to:       NewEmail(name, address),
		subject:  "Upozornenie o aktivite zariadenia",
		template: "adminAuth",
		content:  content,
	}

	if err := msg.parse(); err != nil {
		return err
	}

	return m.Send(msg)
}
