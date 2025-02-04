package email

import (
	"fmt"

	"github.com/hyloblog/hyloblog/internal/config"
	"github.com/hyloblog/hyloblog/internal/email/internal/emailtemplate"
	"github.com/hyloblog/hyloblog/internal/model"
)

const (
	magicRegisterLinkSubject = "Confirm your Hyloblog Account"
	magicLoginLinkSubject    = "Login to Hyloblog"
)

func (s *sender) SendRegisterLink(token string) error {
	text, err := emailtemplate.NewRegisterLink(
		fmt.Sprintf(
			"%s://%s/%s?token=%s",
			config.Config.Hyloblog.Protocol,
			config.Config.Hyloblog.RootDomain,
			"magic/registercallback",
			token,
		),
	).Render(s.emailmode)
	if err != nil {
		return fmt.Errorf("cannot render template: %w", err)
	}
	if err := s.send(
		magicRegisterLinkSubject, text, model.PostmarkStreamOutbound,
	); err != nil {
		return fmt.Errorf("send error: %w", err)
	}
	return nil
}

func (s *sender) SendLoginLink(token string) error {
	text, err := emailtemplate.NewLoginLink(
		fmt.Sprintf(
			"%s://%s/%s?token=%s",
			config.Config.Hyloblog.Protocol,
			config.Config.Hyloblog.RootDomain,
			"magic/logincallback",
			token,
		),
	).Render(s.emailmode)
	if err != nil {
		return fmt.Errorf("cannot render template: %w", err)
	}
	if err := s.send(
		magicLoginLinkSubject, text, model.PostmarkStreamOutbound,
	); err != nil {
		return fmt.Errorf("send error: %w", err)
	}
	return nil
}
