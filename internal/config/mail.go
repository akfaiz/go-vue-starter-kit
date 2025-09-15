package config

import "github.com/akfaiz/go-vue-starter-kit/pkg/env"

type Mail struct {
	SMTP MailSMTP
	From MailFrom
}

type MailSMTP struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MailFrom struct {
	Address string
	Name    string
}

func loadMailConfig() Mail {
	return Mail{
		SMTP: MailSMTP{
			Host:     env.GetString("MAIL_HOST", "127.0.0.1"),
			Port:     env.GetInt("MAIL_PORT", 2525),
			Username: env.GetString("MAIL_USERNAME"),
			Password: env.GetString("MAIL_PASSWORD"),
		},
		From: MailFrom{
			Address: env.GetString("MAIL_FROM_ADDRESS"),
			Name:    env.GetString("MAIL_FROM_NAME"),
		},
	}
}
