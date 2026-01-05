package handlers

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendRegistrationEmail(to, username, password string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	addr := fmt.Sprintf("%s:%s", host, port)

	auth := smtp.PlainAuth("", user, pass, host)

	subject := "Dobrodošli – pristupni podaci"
	body := fmt.Sprintf(`
Dobrodošli!

Vaš korisnički račun je uspješno kreiran.

Korisničko ime: %s
Inicijalna lozinka: %s

Preporučujemo da inicijalnu lozinku promijenite nakon prve prijave.

Lijep pozdrav,
Account Service mShop
`, username, password)

	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body,
	)

	return smtp.SendMail(addr, auth, user, []string{to}, msg)
}
