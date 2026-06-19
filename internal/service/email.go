package service

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderName  string
	SenderEmail string
}

type IEmailService interface {
	SendVerificationEmail(to, name, verificationLink string) error
	SendResetPasswordEmail(to, name, resetLink string) error
	SendNotificationEmail(to, subject, message string) error
}

type EmailService struct {
	config EmailConfig
	log    *logrus.Logger
}

func NewEmailService(config EmailConfig, log *logrus.Logger) IEmailService {
	return &EmailService{
		config: config,
		log:    log,
	}
}

func (s *EmailService) parseTemplate(templateName string, data interface{}) (string, error) {
	tmplPath := filepath.Join("templates", "email", templateName)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func (s *EmailService) sendEmail(to, subject, body string) {
	go func() {
		m := gomail.NewMessage()
		sender := fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.SenderEmail)
		m.SetHeader("From", sender)
		m.SetHeader("To", to)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)

		d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

		if err := d.DialAndSend(m); err != nil {
			s.log.Errorf("Failed to send email to %s: %v", to, err)
		} else {
			s.log.Infof("Email successfully sent to %s", to)
		}
	}()
}

func (s *EmailService) SendVerificationEmail(to, name, verificationLink string) error {
	data := map[string]string{
		"Name":             name,
		"VerificationLink": verificationLink,
	}

	body, err := s.parseTemplate("verification.html", data)
	if err != nil {
		return err
	}

	s.sendEmail(to, "Verifikasi Alamat Email Anda", body)
	return nil
}

func (s *EmailService) SendResetPasswordEmail(to, name, resetLink string) error {
	data := map[string]string{
		"Name":      name,
		"ResetLink": resetLink,
	}

	body, err := s.parseTemplate("reset_password.html", data)
	if err != nil {
		return err
	}

	s.sendEmail(to, "Reset Password Akun Anda", body)
	return nil
}

func (s *EmailService) SendNotificationEmail(to, subject, message string) error {
	data := map[string]string{
		"Message": message,
	}

	body, err := s.parseTemplate("notification.html", data)
	if err != nil {
		return err
	}

	s.sendEmail(to, subject, body)
	return nil
}
