package email

// ContentAuthUser defines content for login token email template.
type ContentAuthUser struct {
	Email    string
	Name     string
	ItemName string
	RoomName string
}

// AuthUserNotif creates and sends a login token email with provided template content.
func (m *Mailer) AuthUserNotif(name, address string, content ContentAuthUser) error {
	msg := &message{
		from:     m.from,
		to:       NewEmail(name, address),
		subject:  "Potvrdenie aktivity zariadenia",
		template: "authUser",
		content:  content,
	}

	if err := msg.parse(); err != nil {
		return err
	}

	return m.Send(msg)
}
