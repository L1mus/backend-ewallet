package mail

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"

	domainAuth "github.com/L1mus/backend-ewallet/internal/domain/auth"
	"github.com/wneessen/go-mail"
)

var (
	ErrMissingSMTPCredentials = errors.New("SMTP credentials are not configured")
	ErrBuildEmailBody         = errors.New("failed to build verification email body")
	ErrSendMail               = errors.New("failed to send verification email")
)

type SMTPClient interface {
	DialAndSendWithContext(ctx context.Context, msgs ...*mail.Msg) error
}

type SMTPClientFactory func(host string, opts ...mail.Option) (SMTPClient, error)

func defaultClientFactory(host string, opts ...mail.Option) (SMTPClient, error) {
	return mail.NewClient(host, opts...)
}

type Mailer struct {
	newClient   SMTPClientFactory
	smtpHost    string
	smtpUser    string
	smtpPass    string
	fromAddress string
}

func NewMailer() (*Mailer, error) {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpUser == "" || smtpPass == "" {
		return nil, ErrMissingSMTPCredentials
	}

	return &Mailer{
		newClient:   defaultClientFactory,
		smtpHost:    os.Getenv("SMTP_HOST"),
		smtpUser:    smtpUser,
		smtpPass:    smtpPass,
		fromAddress: os.Getenv("SMTP_FROM"),
	}, nil
}

func NewMailerWithClientFactory(factory SMTPClientFactory, smtpHost, smtpUser, smtpPass, fromAddress string) *Mailer {
	return &Mailer{
		newClient:   factory,
		smtpHost:    smtpHost,
		smtpUser:    smtpUser,
		smtpPass:    smtpPass,
		fromAddress: fromAddress,
	}
}

func (m *Mailer) SendMail(ctx context.Context, mailReceiver, code string) error {
	message := mail.NewMsg()
	if err := message.From(m.fromAddress); err != nil {
		return fmt.Errorf("%w: set from address: %v", ErrSendMail, err)
	}
	if err := message.To(mailReceiver); err != nil {
		return fmt.Errorf("%w: set to address: %v", ErrSendMail, err)
	}

	message.Subject("Verify Your Email Address")

	htmlBodyEmail, err := BuildVerificationEmailHTML(code)
	if err != nil {
		return err
	}
	message.SetBodyString(mail.TypeTextHTML, htmlBodyEmail)

	client, err := m.newClient(
		m.smtpHost,
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(m.smtpUser),
		mail.WithPassword(m.smtpPass),
	)
	if err != nil {
		return fmt.Errorf("%w: create client: %v", ErrSendMail, err)
	}

	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("%w: %v", ErrSendMail, err)
	}

	log.Printf("verification email sent to %s", mailReceiver)
	return nil
}

type VerificationEmailData struct {
	AppName   string
	Code      string
	ExpiresIn string
}

const verificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="x-apple-disable-message-reformatting" />
    <title>Email Verification</title>
  </head>
  <body
    style="
      margin: 0;
      padding: 0;
      background-color: #f4f8fb;
      font-family: Arial, Helvetica, sans-serif;
      color: #1f2937;
    "
  >
    <table
      role="presentation"
      width="100%"
      cellspacing="0"
      cellpadding="0"
      border="0"
      style="background-color: #f4f8fb; padding: 32px 16px"
    >
      <tr>
        <td align="center">
          <table
            role="presentation"
            width="100%"
            cellspacing="0"
            cellpadding="0"
            border="0"
            style="
              max-width: 600px;
              background-color: #ffffff;
              border-radius: 16px;
              overflow: hidden;
              box-shadow: 0 8px 24px rgba(15, 23, 42, 0.08);
            "
          >
            <tr>
              <td
                align="center"
                style="background-color: #00a6ff; padding: 32px 24px"
              >
                <div
                  style="
                    color: #ffffff;
                    font-size: 26px;
                    line-height: 32px;
                    font-weight: 700;
                  "
                >
                  {{.AppName}}
                </div>
                <div
                  style="
                    margin-top: 10px;
                    color: #e6f7ff;
                    font-size: 18px;
                    line-height: 20px;
                  "
                >
                  Email verification
                </div>
              </td>
            </tr>
            <tr>
              <td style="padding: 36px 32px 28px">
                <h1
                  style="
                    margin: 0 0 16px;
                    color: #111827;
                    font-size: 24px;
                    line-height: 32px;
                  "
                >
                  Verify your email address
                </h1>
                <p
                  style="
                    margin: 0 0 20px;
                    color: #4b5563;
                    font-size: 16px;
                    line-height: 25px;
                  "
                >
                  Use the verification code below to complete your account
                  verification.
                </p>
                <table
                  role="presentation"
                  width="100%"
                  cellspacing="0"
                  cellpadding="0"
                  border="0"
                >
                  <tr>
                    <td
                      align="center"
                      style="
                        background-color: #eaf8ff;
                        border: 1px solid #b5e8ff;
                        border-radius: 12px;
                        padding: 22px 12px;
                      "
                    >
                      <div
                        style="
                          color: #007fbe;
                          font-size: 32px;
                          line-height: 40px;
                          letter-spacing: 10px;
                          font-weight: 700;
                          font-family:
                            &quot;Courier New&quot;, Courier, monospace;
                        "
                      >
                        {{.Code}}
                      </div>
                    </td>
                  </tr>
                </table>
                <p
                  style="
                    margin: 24px 0 0;
                    color: #6b7280;
                    font-size: 14px;
                    line-height: 22px;
                  "
                >
                  This verification code will expire in
                  <strong style="color: #374151">{{.ExpiresIn}}</strong>.
                </p>
                <p
                  style="
                    margin: 12px 0 0;
                    color: #6b7280;
                    font-size: 14px;
                    line-height: 22px;
                  "
                >
                  For your security, do not share this code with anyone.
                  {{.AppName}} will never ask you for your verification code.
                </p>
              </td>
            </tr>
            <tr>
              <td
                style="
                  padding: 22px 32px;
                  background-color: #f8fafc;
                  border-top: 1px solid #e5e7eb;
                "
              >
                <p
                  style="
                    margin: 0;
                    color: #9ca3af;
                    font-size: 12px;
                    line-height: 18px;
                    text-align: center;
                  "
                >
                  If you did not request this verification, you can safely
                  ignore this email.
                </p>
              </td>
            </tr>
          </table>
        </td>
      </tr>
    </table>
  </body>
</html>`

var (
	verificationCodePattern = regexp.MustCompile(`^\d{6}$`)
	verificationEmailTpl    = template.Must(template.New("verification-email").Parse(verificationEmailTemplate))
)

func BuildVerificationEmailHTML(code string) (string, error) {
	if !verificationCodePattern.MatchString(code) {
		return "", domainAuth.ErrInvalidVerificationCode
	}

	data := VerificationEmailData{
		AppName:   "E-Wallet",
		Code:      code,
		ExpiresIn: "10 minutes",
	}

	var body bytes.Buffer
	if err := verificationEmailTpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("%w: %v", ErrBuildEmailBody, err)
	}

	return body.String(), nil
}
