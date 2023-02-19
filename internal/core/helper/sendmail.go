package helper

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"jamo/backend/internal/core/domain/entity"
	"net/smtp"
	"os"
)

//go:embed email-templates/*
var templatesDir embed.FS

// SendMail sends mail
func SendMail(tempName string, data entity.ContactMessage) error {

	smtpHost := Config.SMTPHost
	smtpPort := Config.SMTPPort
	password := os.Getenv("smtp_password")
	from := Config.SMTPUsername

	headers := map[string]string{
		"From":                from,
		"To":                  data.To,
		"Subject":             "Mail from Amirdeen",
		"MIME-Version":        "1.0",
		"Content-Type":        "text/html; charset=utf-8;",
		"Content-Disposition": "inline",
	}

	headerMessage := ""

	for header, value := range headers {
		headerMessage += fmt.Sprintf("%s: %s\r\n", header, value)
	}

	msg, err := getMailBody(tempName, data)
	if err != nil {
		LogEvent("get mail body:", err.Error())
		return err
	}

	body := headerMessage + "\r\n" + msg

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	if err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{data.To}, []byte(body)); err != nil {
		return err
	}

	return nil
}

func getMailBody(templateName string, data entity.ContactMessage) (string, error) {
	tpl, err := template.ParseFS(templatesDir, "email-templates"+"/"+templateName)
	if err != nil {
		fmt.Println("error parsefs:", err)
		return "", err
	}

	buf := &bytes.Buffer{}

	if err := tpl.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
