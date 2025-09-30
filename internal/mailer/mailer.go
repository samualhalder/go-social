package mailer

import "embed"

const (
	FromName                 = "GO-SOCIAL"
	MaxRetries               = 3
	UserRegisterMailTemplate = "registermail.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}
