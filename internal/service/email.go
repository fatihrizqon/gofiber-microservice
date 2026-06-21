package service

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/fatihrizqon/gofiber-microservice/internal/worker"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// EmailConfig is now removed.

type IEmailService interface {
	SendVerificationEmail(to, name, verificationLink string) error
	SendResetPasswordEmail(to, name, resetLink string) error
	SendNotificationEmail(to, subject, message string) error
}

type EmailService struct {
	log         *logrus.Logger
	asynqClient *asynq.Client
}

func NewEmailService(log *logrus.Logger, asynqClient *asynq.Client) IEmailService {
	return &EmailService{
		log:         log,
		asynqClient: asynqClient,
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
	task, err := worker.NewEmailDeliveryTask(to, subject, body)
	if err != nil {
		s.log.Errorf("Failed to create email delivery task: %v", err)
		return
	}

	info, err := s.asynqClient.Enqueue(task)
	if err != nil {
		s.log.Errorf("Failed to enqueue email delivery task: %v", err)
		return
	}
	s.log.Infof("Enqueued email delivery task: id=%s queue=%s", info.ID, info.Queue)
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
