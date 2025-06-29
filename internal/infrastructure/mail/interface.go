package mailer


// For responsiveness, all mail providers 
// and functions must satisfy this interface
type Mailer interface {
	Send(to any, emailData *EmailData) error
}
