package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/babafemi99/WR/internal/config"
	"html/template"
	"log"
	"net/smtp"
)

type IMailer interface {
	SendEmail(toAddr string, content []byte, Dant any, patterns ...string) error
	ShutDown() error
}
type mailgun struct {
	conn     *smtp.Client
	fromAddr string
}

func (m mailgun) ShutDown() error {
	return m.conn.Close()
}

func (m mailgun) SendEmail(toAddr string, content []byte, Dant any, patterns ...string) error {

	for i := range patterns {
		patterns[i] = "templates/" + patterns[i]
	}

	err := m.conn.Mail(m.fromAddr)
	if err != nil {
		return err
	}

	err = m.conn.Rcpt(toAddr)
	if err != nil {
		return err
	}

	data, err := m.conn.Data()
	if err != nil {
		return err
	}

	_, err = data.Write(composeEmail(toAddr, m.fromAddr, patterns, Dant))
	if err != nil {
		return err
	}

	_ = data.Close()
	return nil
}

func NewMailgun(cfg *config.Config) IMailer {

	auth := smtp.PlainAuth("", cfg.EmailUsername, cfg.EmailPassword, cfg.EmailHost)
	log.Println(cfg.EmailHost)

	conn, err := smtp.Dial(fmt.Sprintf("%s:%s", cfg.EmailHost, cfg.EmailPort))
	if err != nil {
		log.Fatalf("[IMAILER: failed to dial smtp: %v]", err)
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	err = conn.StartTLS(tlsConfig)
	if err != nil {
		log.Fatalf("[IMAILER: tls-config: %v]", err)
	}

	err = conn.Auth(auth)
	if err != nil {
		log.Fatalf("[IMAILER: failed to authenticate: %v]", err)
	}

	m := mailgun{conn: conn, fromAddr: cfg.EmailFromAddr}
	return m
}

//func ComposeContent() []byte {
//	from := "info@weddingregistry.com"
//	to := "ooluwa27@gmail.com"
//	subject := "Test Email"
//	body := "This is a test email from Go."
//
//	// Generate plain string representing email data
//	emailString := fmt.Sprintf("From: %s\r\n", from)
//	emailString += fmt.Sprintf("To: %s\r\n", to)
//	emailString += fmt.Sprintf("Subject: %s\r\n", subject)
//	emailString += "\r\n" + body
//
//	return []byte(emailString)
//
//}

func composeEmail(recipient, sender string, patterns []string, data interface{}) []byte {
	// Create a new buffer to store the email message
	var buf bytes.Buffer

	// Write the "To" and "From" headers
	fmt.Fprintf(&buf, "To: %s\r\n", recipient)
	fmt.Fprintf(&buf, "From: %s\r\n", sender)

	// Load and execute templates for subject, plain text, and HTML
	subjectTemplate, plainBodyTemplate, htmlBodyTemplate := loadTemplates(patterns)

	subject := executeTemplate(subjectTemplate, data)
	fmt.Fprintf(&buf, "Subject: %s\r\n", subject)

	// Start the MIME structure for a multipart/alternative email
	fmt.Fprintf(&buf, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=boundary-string\r\n")
	fmt.Fprintf(&buf, "\r\n")
	fmt.Fprintf(&buf, "--boundary-string\r\n")

	plainBody := executeTemplate(plainBodyTemplate, data)
	fmt.Fprintf(&buf, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
	fmt.Fprintf(&buf, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(&buf, "\r\n")
	fmt.Fprintf(&buf, "%s\r\n", plainBody)

	if htmlBodyTemplate != nil {
		htmlBody := executeTemplate(htmlBodyTemplate, data)
		fmt.Fprintf(&buf, "--boundary-string\r\n")
		fmt.Fprintf(&buf, "Content-Type: text/html; charset=\"utf-8\"\r\n")
		fmt.Fprintf(&buf, "Content-Transfer-Encoding: quoted-printable\r\n")
		fmt.Fprintf(&buf, "\r\n")
		fmt.Fprintf(&buf, "%s\r\n", htmlBody)
	}

	// Close the MIME structure
	fmt.Fprintf(&buf, "--boundary-string--\r\n")
	// Convert the buffer to a byte slice and return it
	return buf.Bytes()
}

func loadTemplate(ts *template.Template, patterns []string, name string) *template.Template {
	ts, err := ts.ParseFS(EmbeddedFiles, patterns...)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return nil
	}

	return ts.Lookup(name)
}

func loadTemplates(patterns []string) (*template.Template, *template.Template, *template.Template) {
	ts := template.New("")

	subjectTemplate := loadTemplate(ts, patterns, "subject")
	plainBodyTemplate := loadTemplate(ts, patterns, "plainBody")
	htmlBodyTemplate := loadTemplate(ts, patterns, "htmlBody")

	return subjectTemplate, plainBodyTemplate, htmlBodyTemplate
}

func executeTemplate(tmpl *template.Template, data interface{}) string {
	var buffer bytes.Buffer
	if tmpl != nil {
		if err := tmpl.Execute(&buffer, data); err != nil {
			// Handle the error appropriately
			fmt.Println("Error executing template:", err)
			return ""
		}
	}
	return buffer.String()
}
