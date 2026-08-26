package mail

import (
	"context"
	"errors"
	"strings"
	"testing"

	domainAuth "github.com/L1mus/backend-ewallet/internal/domain/auth"
	"github.com/joho/godotenv"
	"github.com/wneessen/go-mail"
)

func TestBuildVerificationEmailHTML_Success(t *testing.T) {
	code := "422123"

	html, err := BuildVerificationEmailHTML(code)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, code) {
		t.Errorf("expected generated HTML to contain code %q, but it didn't", code)
	}
}

func TestBuildVerificationEmailHTML_InvalidCode(t *testing.T) {
	code := "4221A3"

	_, err := BuildVerificationEmailHTML(code)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if !errors.Is(err, domainAuth.ErrInvalidVerificationCode) {
		t.Errorf("expected error to be ErrInvalidVerificationCode, got: %v", err)
	}
}

type fakeSMTPClient struct {
	dialAndSendErr error
}

func (f *fakeSMTPClient) DialAndSendWithContext(_ context.Context, _ ...*mail.Msg) error {
	return f.dialAndSendErr
}

func newTestMailer(factory SMTPClientFactory) *Mailer {
	return NewMailerWithClientFactory(factory, "smtp.test.local", "test-user", "test-pass", "no-reply@example.com")
}

func TestMailer_SendMail_Success(t *testing.T) {
	factory := func(host string, opts ...mail.Option) (SMTPClient, error) {
		return &fakeSMTPClient{}, nil
	}
	m := newTestMailer(factory)

	err := m.SendMail(context.Background(), "testmail@gmail.com", "422113")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMailer_SendMail_NoFromAddress(t *testing.T) {
	factory := func(host string, opts ...mail.Option) (SMTPClient, error) {
		t.Fatal("client factory should not be called when From address is invalid")
		return nil, nil
	}
	m := NewMailerWithClientFactory(factory, "smtp.test.local", "test-user", "test-pass", "" /* fromAddress kosong */)

	err := m.SendMail(context.Background(), "testmail@gmail.com", "422113")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrSendMail) {
		t.Errorf("expected error to be ErrSendMail, got: %v", err)
	}
	if !strings.Contains(err.Error(), "set from address") {
		t.Errorf("expected error message to mention 'set from address', got: %v", err)
	}
}

func TestMailer_SendMail_NoMailReceiver(t *testing.T) {
	factory := func(host string, opts ...mail.Option) (SMTPClient, error) {
		t.Fatal("client factory should not be called when receiver address is invalid")
		return nil, nil
	}
	m := newTestMailer(factory)

	err := m.SendMail(context.Background(), "", "422113")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrSendMail) {
		t.Errorf("expected error to be ErrSendMail, got: %v", err)
	}
}

func TestMailer_SendMail_FailedToCreateClient(t *testing.T) {
	factory := func(host string, opts ...mail.Option) (SMTPClient, error) {
		return nil, errors.New("could not resolve smtp host")
	}
	m := newTestMailer(factory)

	err := m.SendMail(context.Background(), "testmail@gmail.com", "422113")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrSendMail) {
		t.Errorf("expected error to be ErrSendMail, got: %v", err)
	}
}

func TestMailer_SendMail_DialAndSendFails(t *testing.T) {
	factory := func(host string, opts ...mail.Option) (SMTPClient, error) {
		return &fakeSMTPClient{dialAndSendErr: errors.New("connection refused")}, nil
	}
	m := newTestMailer(factory)

	err := m.SendMail(context.Background(), "testmail@gmail.com", "422113")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, ErrSendMail) {
		t.Errorf("expected error to be ErrSendMail, got: %v", err)
	}
}

func TestNewMailer_MissingCredentials(t *testing.T) {
	t.Setenv("SMTP_USER", "")
	t.Setenv("SMTP_PASS", "")
	t.Setenv("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	t.Setenv("SMTP_FROM", "no-reply@example.com")

	_, err := NewMailer()
	if !errors.Is(err, ErrMissingSMTPCredentials) {
		t.Errorf("expected error to be ErrMissingSMTPCredentials, got: %v", err)
	}
}

func TestNewMailer_Success(t *testing.T) {
	t.Setenv("SMTP_USER", "test-user")
	t.Setenv("SMTP_PASS", "test-pass")
	t.Setenv("SMTP_HOST", "sandbox.smtp.mailtrap.io")
	t.Setenv("SMTP_FROM", "no-reply@example.com")

	m, err := NewMailer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected a non-nil Mailer")
	}
}

func TestMailer_SendMail_Integration_Success(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Skip("skipping integration test: .env not found")
	}

	m, err := NewMailer()
	if err != nil {
		t.Skip("skipping integration test: SMTP credentials not configured")
	}

	if err := m.SendMail(context.Background(), "testmail@gmail.com", "422113"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
