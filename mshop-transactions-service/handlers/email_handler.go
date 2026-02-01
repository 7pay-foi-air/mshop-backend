package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
)

type EmailAttachment struct {
	Filename string
	Content  []byte
	MimeType string
}

func SendEmail(
	to string,
	subject string,
	body string,
	attachment *EmailAttachment,
) error {

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		return fmt.Errorf("SMTP env vars not configured")
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	auth := smtp.PlainAuth("", user, pass, host)

	var msg bytes.Buffer
	boundary := "mshop-boundary-123456"

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	if attachment != nil {
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n\r\n")

		msg.WriteString("--" + boundary + "\r\n")
		msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
		msg.WriteString(body + "\r\n\r\n")

		msg.WriteString("--" + boundary + "\r\n")
		msg.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", attachment.MimeType, attachment.Filename))
		msg.WriteString("Content-Transfer-Encoding: base64\r\n")
		msg.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", attachment.Filename))

		encoded := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Content)))
		base64.StdEncoding.Encode(encoded, attachment.Content)

		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			msg.Write(encoded[i:end])
			msg.WriteString("\r\n")
		}

		msg.WriteString("\r\n--" + boundary + "--")
	} else {
		msg.WriteString("\r\n" + body)
	}

	return smtp.SendMail(addr, auth, user, []string{to}, msg.Bytes())
}
