package mailer

// ResponsiveMailer tries multiple mail providers in order until one succeeds.
type ResponsiveMailer struct {
	Providers []Mailer // List of mailers to try, in priority order
}

// Send attempts to deliver the email using the configured providers.
func (rm *ResponsiveMailer) Send(to any, emailData *EmailData) error {
	var lastErr error // Stores the most recent error to return if all fail

	// Loop through each provider in order
	for _, provider := range rm.Providers {
		// Try sending the email with the current provider
		err := provider.Send(to, emailData)

		if err == nil {
			// If sending succeeded, return nil immediately
			return nil
		}

		// If sending failed, store this error and continue to the next provider
		lastErr = err
	}

	// If we got here, 
	// all providers failed.
	return lastErr // Return the last error
}
