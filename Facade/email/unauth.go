package email

// ContentUnauthNotif defines content for login token email template.
type ContentUnauthNotif struct {
	Email    string
	Name     string
	ItemName string
	RoomName string
}

// UnatuhNorification creates and sends a login token email with provided template content.
func (m *Mailer) UnatuhNorification(name, address string, content ContentUnauthNotif) error {
	msg := &message{
		from:     m.from,
		to:       NewEmail(name, address),
		subject:  "Upozornenie o aktivite zariadenia",
		template: "unauth",
		content:  content,
	}

	if err := msg.parse(); err != nil {
		return err
	}

	return m.Send(msg)
}
