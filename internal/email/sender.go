package email

import (
	"github.com/hylodoc/hyloblog/internal/email/emailaddr"
	"github.com/hylodoc/hyloblog/internal/email/emailqueue"
	"github.com/hylodoc/hyloblog/internal/email/postbody"
	"github.com/hylodoc/hyloblog/internal/model"
)

type Sender interface {
	SendRegisterLink(token string) error
	SendLoginLink(token string) error

	SendNewSubscriberEmail(sitename, unsublink string) error

	SendNewPostUpdate(
		posttitle, postlink, unsublink string,
		body postbody.PostBody,
	) error
}

type sender struct {
	to, from  emailaddr.EmailAddr
	emailmode model.EmailMode
	_store    *model.Store
}

func NewSender(
	to, from emailaddr.EmailAddr, mode model.EmailMode, store *model.Store,
) Sender {
	return &sender{to, from, mode, store}
}

func (s *sender) send(subject, body string, stream model.PostmarkStream) error {
	return s.sendwithheaders(subject, body, nil, stream)
}

func (s *sender) sendwithheaders(
	subject, body string, headers map[string]string,
	stream model.PostmarkStream,
) error {
	return emailqueue.NewEmail(
		s.from.Addr(),
		s.to.Addr(),
		subject,
		body,
		s.emailmode,
		stream,
		headers,
	).Queue(s._store)
}
