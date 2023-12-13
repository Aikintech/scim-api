package jobs

import (
	"fmt"
	"os"

	"github.com/aikintech/scim-api/pkg/constants"
	"github.com/aikintech/scim-api/pkg/models"
	mailjet "github.com/mailjet/mailjet-apiv3-go"
)

type Mail struct {
	mailClient *mailjet.Client
}

type SendMailRequest struct {
	FromEmail  string
	FromName   string
	ToEmail    string
	ToName     string
	Subject    string
	TemplateID int64
	Variables  map[string]interface{}
}

func NewMail() *Mail {
	return &Mail{
		mailClient: mailjet.NewMailjetClient(os.Getenv("MAILJET_API_KEY"), os.Getenv("MAILJET_SECRET_KEY")),
	}
}

func (m *Mail) SendUserWelcomeMail(user models.User) {
	r := SendMailRequest{
		FromEmail:  constants.NO_REPLY_EMAIL,
		FromName:   "SCIM",
		ToEmail:    user.Email,
		ToName:     user.FirstName,
		Subject:    "Welcome to the SCIM APP community!",
		TemplateID: constants.MAILJET_WELCOME_MAIL_TEMPLATE_ID,
		Variables: map[string]interface{}{
			"name": user.FirstName,
		},
	}

	m.Send(r)
}

func (m *Mail) SendUserPasswordResetMail(user models.User, code string) {
	r := SendMailRequest{
		FromEmail:  constants.SUPPORT_EMAIL,
		FromName:   "SCIM Support",
		ToEmail:    user.Email,
		ToName:     user.FirstName,
		Subject:    "Reset your password for SCIM APP",
		TemplateID: constants.MAILJET_RESET_PASSWORD_MAIL_TEMPLATE_ID,
		Variables: map[string]interface{}{
			"name": user.FirstName,
			"code": code,
		},
	}

	m.Send(r)
}

func (m *Mail) SendUserVerificationMail(user models.User, code string) {
	r := SendMailRequest{
		FromEmail:  constants.NO_REPLY_EMAIL,
		FromName:   "SCIM",
		ToEmail:    user.Email,
		ToName:     user.FirstName,
		Subject:    "Please verify your email address",
		TemplateID: constants.MAILJET_VERIFY_MAIL_TEMPLATE_ID,
		Variables: map[string]interface{}{
			"name": user.FirstName,
			"code": code,
		},
	}

	m.Send(r)
}

func (m *Mail) Send(request SendMailRequest) {
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: request.FromEmail,
				Name:  request.FromName,
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: request.ToEmail,
					Name:  request.ToName,
				},
			},
			TemplateID:       request.TemplateID,
			TemplateLanguage: true,
			Subject:          request.Subject,
			Variables:        request.Variables,
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := m.mailClient.SendMailV31(&messages)

	if err != nil {
		fmt.Println("Error sending mail")
		fmt.Println(err.Error())
	}

	fmt.Printf("Data: %+v\n", res)
}
