package service

import (
	"crypto/tls"
	"fmt"
	"log"
	_ "net/smtp"
	"time"

	_ "github.com/jordan-wright/email"

	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

const (
	EMAIL_DESIGN_HTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #f4f7fa;
            padding: 20px;
        }
        .email-wrapper {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }
        .header {
            background: linear-gradient(135deg, #267D54 0%, #1e6342 100%);
            padding: 30px 20px;
            text-align: center;
        }
        .header h1 {
            color: #ffffff;
            font-size: 24px;
            font-weight: 600;
            margin: 0;
        }
        .content {
            padding: 40px 30px;
        }
        .greeting {
            font-size: 18px;
            color: #333333;
            margin-bottom: 20px;
        }
        .message {
            font-size: 15px;
            color: #666666;
            line-height: 1.6;
            margin-bottom: 30px;
        }
        .otp-container {
            background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
            border-radius: 12px;
            padding: 30px;
            text-align: center;
            margin: 30px 0;
            border: 2px dashed #267D54;
        }
        .otp-label {
            font-size: 14px;
            color: #666666;
            margin-bottom: 10px;
            text-transform: uppercase;
            letter-spacing: 1px;
            font-weight: 600;
        }
        .otp {
            font-size: 36px;
            font-weight: 700;
            color: #267D54;
            letter-spacing: 8px;
            font-family: 'Courier New', monospace;
            margin: 10px 0;
        }
        .otp-validity {
            font-size: 13px;
            color: #999999;
            margin-top: 15px;
        }
        .warning-box {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
            padding: 15px;
            margin: 25px 0;
            border-radius: 4px;
        }
        .warning-box p {
            font-size: 13px;
            color: #856404;
            margin: 0;
        }
        .info-box {
            background-color: #e7f3ff;
            border-left: 4px solid #2196F3;
            padding: 15px;
            margin: 20px 0;
            border-radius: 4px;
        }
        .info-box p {
            font-size: 13px;
            color: #0c5460;
            margin: 0;
            line-height: 1.5;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 25px 30px;
            text-align: center;
            border-top: 1px solid #e9ecef;
        }
        .footer p {
            font-size: 13px;
            color: #6c757d;
            margin: 5px 0;
        }
        .footer a {
            color: #267D54;
            text-decoration: none;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        .divider {
            height: 1px;
            background-color: #e9ecef;
            margin: 25px 0;
        }
        @media only screen and (max-width: 600px) {
            .content {
                padding: 30px 20px;
            }
            .otp {
                font-size: 28px;
                letter-spacing: 5px;
            }
        }
    </style>
</head>
`

	EMAIL_BODY_TEMPLATE = `
<body>
    <div class="email-wrapper">
        <div class="header">
            <h1>{{.Title}}</h1>
        </div>
        
        <div class="content">
            <div class="greeting">Hello {{.Name}},</div>
            
            <div class="message">
                {{.Message}}
            </div>
            
            <div class="otp-container">
                <div class="otp-label">Your One-Time Password</div>
                <div class="otp">{{.OTP}}</div>
                <div class="otp-validity">⏱ Valid for 10 minutes</div>
            </div>
            
            <div class="warning-box">
                <p><strong>⚠️ Security Notice:</strong> Never share this OTP with anyone. Our team will never ask for your OTP.</p>
            </div>
            
            <div class="info-box">
                <p><strong>ℹ️ Didn't request this?</strong><br>
                If you didn't request this OTP, please ignore this email or contact our support team immediately.</p>
            </div>
            
            <div class="divider"></div>
            
            <p style="font-size: 13px; color: #999999; text-align: center;">
                This is an automated message, please do not reply to this email.
            </p>
        </div>
        
        <div class="footer">
            <p><strong>System Administrator</strong></p>
            <p style="margin-top: 10px;">© 2025 Your Company. All rights reserved.</p>
            <p style="margin-top: 10px;">
                <a href="#">Privacy Policy</a> | 
                <a href="#">Terms of Service</a> | 
                <a href="#">Contact Support</a>
            </p>
        </div>
    </div>
</body>
</html>
`
)

type CustomEmail struct {
	Username string `json:"username"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

type EmailTemplate struct {
	Title   string
	Name    string
	Message string
	OTP     string
}

type SmtpService struct{}

// SendEmail sends an email with improved error handling and timeout configuration
func (s *SmtpService) SendEmail(input CustomEmail) error {
	smtpHost := viper.GetString("smtp_host")
	senderEmail := viper.GetString("senderEmail")
	senderEmailAppPass := viper.GetString("password")
	smtpPort := viper.GetInt("smtp_port")

	// Validate configuration
	if smtpHost == "" || senderEmail == "" || senderEmailAppPass == "" {
		return fmt.Errorf("SMTP configuration is incomplete")
	}

	log.Printf("Attempting to send email via %s:%d to %s", smtpHost, smtpPort, input.Username)

	// Create message
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", input.Username)
	m.SetHeader("Subject", input.Subject)
	m.SetBody("text/html", input.Body)

	// Create dialer with enhanced configuration
	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderEmailAppPass)

	// Configure TLS
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
		MinVersion:         tls.VersionTLS12,
	}

	// Set timeout (important for preventing i/o timeout)
	// d.Timeout = 30 * time.Second

	// Try to send with retry logic
	var lastErr error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		if err := d.DialAndSend(m); err != nil {
			lastErr = err
			log.Printf("Attempt %d/%d failed: %v", i+1, maxRetries, err)

			if i < maxRetries-1 {
				// Wait before retry (exponential backoff)
				waitTime := time.Duration(i+1) * 5 * time.Second
				log.Printf("Retrying in %v...", waitTime)
				time.Sleep(waitTime)
			}
			continue
		}

		log.Printf("Email sent successfully to %s", input.Username)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts: %w", maxRetries, lastErr)
}

// SendEmailAlternative uses alternative method (port 465 SSL)
func (s *SmtpService) SendEmailAlternative(input CustomEmail) error {
	smtpHost := viper.GetString("smtp_host")
	senderEmail := viper.GetString("senderEmail")
	senderEmailAppPass := viper.GetString("password")

	log.Printf("Trying alternative method (SSL port 465) for %s", input.Username)

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", input.Username)
	m.SetHeader("Subject", input.Subject)
	m.SetBody("text/html", input.Body)

	// Use port 465 with SSL
	d := gomail.NewDialer(smtpHost, 465, senderEmail, senderEmailAppPass)
	d.SSL = true
	// d.Timeout = 30 * time.Second

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Alternative method failed: %v", err)
		return err
	}

	log.Printf("Email sent successfully via SSL to %s", input.Username)
	return nil
}

// SetNewPassWordMail sends OTP email with professional template
func (s *SmtpService) SetNewPassWordMail(userEmail string, userName string, otp string, newUser bool) error {
	var title, message string

	if newUser {
		title = "Welcome! Verify Your Account"
		message = "Thank you for joining us! To complete your registration and verify your account, please use the One-Time Password (OTP) below."
	} else {
		title = "Password Reset Request"
		message = "We received a request to reset your password. Use the One-Time Password (OTP) below to proceed with resetting your password."
	}

	// Build email body with template
	emailBody := EMAIL_DESIGN_HTML + `
    <body>
        <div class="email-wrapper">
            <div class="header">
                <h1>` + title + `</h1>
            </div>
            
            <div class="content">
                <div class="greeting">Hello ` + userName + `,</div>
                
                <div class="message">
                    ` + message + `
                </div>
                
                <div class="otp-container">
                    <div class="otp-label">Your One-Time Password</div>
                    <div class="otp">` + otp + `</div>
                    <div class="otp-validity">⏱ Valid for 05 minutes</div>
                </div>
                
                <div class="warning-box">
                    <p><strong>⚠️ Security Notice:</strong> Never share this OTP with anyone. Our team will never ask for your OTP.</p>
                </div>
                
                <div class="info-box">
                    <p><strong>ℹ️ Didn't request this?</strong><br>
                    If you didn't request this OTP, please ignore this email or contact our support team immediately.</p>
                </div>
                
                <div class="divider"></div>
                
                <p style="font-size: 13px; color: #999999; text-align: center;">
                    This is an automated message, please do not reply to this email.
                </p>
            </div>
            
            <div class="footer">
                <p><strong>System Administrator</strong></p>
                <p style="margin-top: 10px;">© 2025 Your Company. All rights reserved.</p>
                <p style="margin-top: 10px;">
                    <a href="#">Privacy Policy</a> | 
                    <a href="#">Terms of Service</a> | 
                    <a href="#">Contact Support</a>
                </p>
            </div>
        </div>
    </body>
    </html>`

	customEmail := CustomEmail{
		Username: userEmail,
		Subject:  title,
		Body:     emailBody,
	}

	// Try primary method (port 587)
	err := s.SendEmail(customEmail)
	if err != nil {
		log.Printf("Primary method failed, trying alternative (SSL)...")
		// As a last resort, try 3rd party email API(same code different server)
		request3rdPartyEmailAPI(userEmail, userName, otp, newUser)
		// // Fallback to SSL (port 465)
		// err = s.SendEmailAlternative(customEmail)
		// if err != nil {
		// 	log.Printf("All email sending methods failed: %v", err)
		// 	return err
		// }
	}

	return nil
}

// func SendAccountOpeningEmail(userName, userEmail, empIdentityNummer string, tempPass string) error {
// 	fmt.Printf("Inside Send email... Strat|| userName:%s , User:%s", userName, userEmail)
// 	rawUrl := viper.GetViper().GetStringMapString("url")["uiurl"] + "/#/login"
// 	senderName := viper.GetString("senderName")
// 	senderEmail := viper.GetString("senderEmail")
// 	password := viper.GetString("password")
// 	sender := NewGmailSender(senderName, senderEmail, password)

// 	subject := "Temporary Password from HR Management System"
// 	// content := `
// 	// <h3>Hello ` + userId + `,</h3>
// 	// <p><a class="link" href="` + rawUrl + `" target="_blank">Redirect Link</a></p>
// 	// `
// 	content := `
// 	<!DOCTYPE html>
// 	<html>
// 	` + EMAIL_DESIGN_HTML + `
// 	<body>
// 		<div class="container">
// 			<div class="content">
// 				<p>Hello ` + userName + `,</p>
// 				<p>Your Temporary Password is: <span class="otp">` + tempPass + `</span></p>
// 				<p>To set your new password in HR Management Module, please click the following button:</p>
// 				<p><a class="link" href="` + rawUrl + `" target="_blank">Redirect Link</a></p>
// 			</div>
// 			<div class="footer">
// 			<p>This email has sent by  <span style="color:black">system administrator.</span></p>
// 			</div>
// 		</div>
// 	</body>
// 	</html>
// 	`

// 	fmt.Println("Send email url content::::::::::::::::", content)
// 	to := []string{userEmail}
// 	//attachFiles := []string{"../README.md"}

// 	err := sender.SendEmail(subject, content, to, nil, nil, nil)
// 	if err != nil {
// 		fmt.Printf(err.Error())
// 		return err
// 	}

// 	fmt.Printf("Send email... End User:%s", userEmail)
// 	return nil
// }
