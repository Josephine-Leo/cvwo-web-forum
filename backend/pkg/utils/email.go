package utils

import (
	"bytes"
	"path/filepath"
	"text/template"

	"fmt"
	"os"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	URL     string
	Subject string
}

// Template
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	fmt.Println("Am parsing templates...")

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

// email template parser
func SendEmail(to string, data *EmailData, templateName string) error {
	from := "cvwoJose@gmail.com"
	smtpUser := "cvwojose@gmail.com"
	smtpPass := "bhmd vwot nkos hbvu " //edit
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	tmpl, err := ParseTemplateDir("pkg/templates")
	if err != nil {
		return err
	}

	t := tmpl.Lookup(templateName)
	if t == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	return d.DialAndSend(m)
}
