package mailer

import (
	"sync"
	"fmt"
	"net/smtp"

    "japa/internal/config"

	"go.uber.org/zap"
)

type SMTPMailer struct {
	ServerConfig    config.ServerConfig // For server settings
	SiteConfig      config.SiteConfig   // Site details for email template fill
	EmailConfig     config.SMTPConfig   // Email settings
	Logger          *zap.Logger         // For structured logging
}

// Initialize SMTPMailer
func NewSMTPMailer(
	serverConfig config.ServerConfig,
	siteConfig   config.SiteConfig, 
	emailConfig  config.SMTPConfig, 
	logger       *zap.Logger,
) *SMTPMailer {
	return &SMTPMailer{
		ServerConfig: serverConfig,
		SiteConfig:   siteConfig,
		EmailConfig:  emailConfig,
		Logger:       logger,
	}
}

// Sends an email using an HTML template and dynamic data
func (s *SMTPMailer) Send(to any, emailData *EmailData) error {
	body, err := ParseEmailTemplate(s.ServerConfig.TemplatesDir, emailData, s.Logger)
	if err != nil {
		return err
	}
	
	if err := s.send(to, emailData.Subject, body); err != nil {
		s.Logger.Error("error sending mail via smtp", zap.Error(err))
		return err
	}

	return nil
}

// Handles the low-level email delivery.
// to @param is expecting either string or []string
func (s *SMTPMailer) send(to any, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.EmailConfig.EMAIL_HOST, s.EmailConfig.EMAIL_PORT)
	auth := smtp.PlainAuth("", s.EmailConfig.EMAIL_USERNAME, s.EmailConfig.EMAIL_PASSWORD, s.EmailConfig.EMAIL_HOST)

	var recipients []string

	switch v := to.(type) {
	case string:
		recipients = []string{v}
	case []string:
		recipients = v
	default:
		return fmt.Errorf("unsupported type for recipient: %T", v)
	}

	if len(recipients) <= 20 {
		// If receipients is less than 20 then loop normally

		for _, recipient := range recipients {
			msg := []byte(fmt.Sprintf(
				"To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
				recipient, subject, body,
			))

			err := smtp.SendMail(addr, auth, s.SiteConfig.SiteEmail, []string{recipient}, msg)
			if err != nil {
				s.Logger.Error("send mail failed", zap.String("to", recipient), zap.Error(err))
				return err
			}
		}
		return nil

	} else {
		// For larger recipient lists (20 or more), we use goroutines to send emails concurrently.
		// This helps improve performance, since sending emails via SMTP can be slow for each recipient.

		const maxConcurrentWorkers = 10 // Limit the number of concurrent emails being sent to avoid overloading the server or hitting SMTP rate limits.

		var waitGroup sync.WaitGroup                       // Used to wait for all email-sending goroutines to finish.
		semaphore := make(chan struct{}, maxConcurrentWorkers) // Semaphore/workerLimit to control how many goroutines can run at the same time.

		for _, recipient := range recipients {
			waitGroup.Add(1)             // Track this goroutine
			semaphore <- struct{}{}      // Block if maxConcurrentWorkers is reached.

			go func(email string) {
				defer waitGroup.Done()     // Mark this goroutine as done
				defer func() { <-semaphore }() // Free up a slot in the semaphore

				// Recover from any unexpected panic so the app doesn't crash
				defer func() {
					if panicErr := recover(); panicErr != nil {
						s.Logger.Error("Panic recovered in goroutine", zap.Any("error", panicErr))
					}
				}()

				// Attempt to send the email
				err := s.send(email, subject, body)
				if err != nil {
					s.Logger.Error("Failed to send email", zap.String("to", email), zap.Error(err))
				}
			}(recipient)
		}

		waitGroup.Wait() // Wait for all emails to be sent before continuing
	}

	return nil
}

/*
## 🧠 1. Why was it called **“semaphore”?**

> It’s just a **nickname** or term developers use.

### ✅ In Go:

* A `semaphore` is **not a keyword or built-in type**.
* It’s just a **`chan struct{}`** (a channel of empty structs).
* But we **treat it** like a **counting gate** — allowing only `N` goroutines to run at the same time.

So we name it `semaphore` or `workerLimit` to make its purpose clearer.

---

## ⚙️ 2. What does `make(chan struct{}, 10)` really do?

```go
workerLimit := make(chan struct{}, 10)
```

This creates a **buffered channel** with:

* A capacity of **10**
* That holds **empty values** (`struct{}` is a type with **zero memory**)

Think of it like:

> “This bucket can hold up to 10 tokens.”

---

## 📥 3. What does `workerLimit <- struct{}{}` mean?

```go
workerLimit <- struct{}{}
```

This sends an **empty struct** into the channel — **taking up one slot**.

* If the channel already holds 10 items (full), this line **blocks (waits)**.
* If there’s space, it goes through and the goroutine continues.

### ✅ So this line is saying:

> "I want to run. If the limit is reached, I’ll wait until a slot is free."

---

## 📤 4. What does `<-workerLimit` mean?

```go
<-workerLimit
```

This **removes** one item from the channel — **freeing up one slot**.

Think of it as:

> "I'm done. Let someone else take this slot."

So when the goroutine ends:

* It signals: “I'm done, next worker can run.”

---

## 🧠 Summary Table

| Code                        | Meaning                                     |
| --------------------------- | ------------------------------------------- |
| `make(chan struct{}, 10)`   | A 10-slot queue for controlling concurrency |
| `workerLimit <- struct{}{}` | Reserve a slot. Wait if all are taken.      |
| `<-workerLimit`             | Release the slot for someone else to use    |

---

## 🧴 Real-life Analogy

Imagine a **restaurant with 10 seats**:

* When someone walks in, they **take a seat** (`<- struct{}{}`).
* If all seats are taken, the next person **waits outside** (blocks).
* When someone finishes eating, they **leave their seat** (`<-semaphore`), and the next person is let in.

That’s **exactly** what’s happening here.

Excellent questions — you’re now getting into the real **syntax and design decisions of Go**. Let's walk through both of your questions carefully:

---

## ✅ 1. Why use `struct{}` and can we use other types?

Yes, you can use **any type** in a channel — but `struct{}` is the most efficient for control signals because it **uses zero memory**.

### Example using another type (e.g. `int`):

```go
sem := make(chan int, 3) // Semaphore with 3 slots

// Acquire a slot
sem <- 1  // or any value, it doesn't matter

// Release a slot
<-sem
```

You could even use strings:

```go
sem := make(chan string, 2)

sem <- "token1"
sem <- "token2"

<-sem // take out "token1"
```

But using values here **wastes memory** and isn’t needed when the value itself doesn't matter — we only care about *blocking or not blocking*.

### So why `struct{}`?

```go
struct{}{} // literally zero-sized placeholder

sem := make(chan struct{}, 5)
sem <- struct{}{}  // fills one slot
<-sem              // frees one slot
```

Using `struct{}` is like saying:

> "I don't care what the data is, I just want to block/unblock."

---

## ✅ 2. What’s the purpose of `struct{}{}` and `()`, especially in `go func() {...}()`?

### 🔸 `struct{}{}`

* `struct{}` defines the **type**.
* `struct{}{}` is how you **instantiate** it.

Think of it like:

```go
type Empty struct{}     // a custom empty struct type
value := Empty{}        // create a value of that type
```

So:

```go
struct{}{}  // anonymous empty struct type + value
```

---

### 🔸 Why `()` in `go func(...) {...}()`

This part can be confusing, so let’s demystify:

```go
go func(email string) {
	// code here
}(recipient)
```

That’s an **anonymous function** (function literal) being **defined and immediately called**.

### Broken down:

#### Step 1: Declare function

```go
func(email string) {
	fmt.Println(email)
}
```

#### Step 2: Call the function immediately

```go
func(email string) {
	fmt.Println(email)
}("user@example.com")  // now it runs immediately
```

#### Step 3: Run it as a goroutine

```go
go func(email string) {
	fmt.Println(email)
}("user@example.com")  // in background
```

### ✅ Summary Table

| Syntax                   | Meaning                                        |
| ------------------------ | ---------------------------------------------- |
| `struct{}`               | Defines an empty struct type                   |
| `struct{}{}`             | Instantiates it (zero memory value)            |
| `func(...) { ... }()`    | Anonymous function, immediately executed       |
| `go func(...) { ... }()` | Run that anonymous function **as a goroutine** |

---
*/
