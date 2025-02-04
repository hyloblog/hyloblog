package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/hyloblog/hyloblog/internal/app/handler/response"
	"github.com/hyloblog/hyloblog/internal/assert"
	"github.com/hyloblog/hyloblog/internal/authz"
	"github.com/hyloblog/hyloblog/internal/config"
	"github.com/hyloblog/hyloblog/internal/session"
	"github.com/hyloblog/hyloblog/internal/util"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	sesh, ok := r.Context().Value(session.CtxSessionKey).(*session.Session)
	assert.Assert(ok)

	if errors.Is(err, authz.SubscriptionError) {
		sesh.Println("authz error:", err)
		unauthorised(w, r)
		return
	}

	sesh.Println("internal server error:", err)
	internalServerError(w, r)
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	sesh, ok := r.Context().Value(session.CtxSessionKey).(*session.Session)
	assert.Assert(ok)
	w.WriteHeader(http.StatusInternalServerError)
	if err := response.NewTemplate(
		[]string{"500.html"},
		util.PageInfo{
			Data: struct {
				Title        string
				UserInfo     *session.UserInfo
				DiscordURL   string
				OpenIssueURL string
			}{
				Title:        "Hyloblog – Internal server error",
				UserInfo:     session.ConvertSessionToUserInfoError(sesh),
				DiscordURL:   config.Config.Hyloblog.DiscordURL,
				OpenIssueURL: config.Config.Hyloblog.OpenIssueURL,
			},
		},
	).Respond(w, r); err != nil {
		sesh.Println(
			"pathological error:",
			err,
		)
	}
}

func unauthorised(w http.ResponseWriter, r *http.Request) {
	sesh, ok := r.Context().Value(session.CtxSessionKey).(*session.Session)
	assert.Assert(ok)
	w.WriteHeader(http.StatusInternalServerError)
	if err := response.NewTemplate(
		[]string{"401.html"},
		util.PageInfo{
			Data: struct {
				Title    string
				UserInfo *session.UserInfo
			}{
				Title:    "Hyloblog – Unauthorised",
				UserInfo: session.ConvertSessionToUserInfoError(sesh),
			},
		},
	).Respond(w, r); err != nil {
		sesh.Println(
			"pathological error:",
			err,
		)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	sesh, ok := r.Context().Value(session.CtxSessionKey).(*session.Session)
	assert.Assert(ok)
	sesh.Println("404", r.URL)
	w.WriteHeader(http.StatusNotFound)
	if err := response.NewTemplate(
		[]string{"404.html"},
		util.PageInfo{
			Data: struct {
				Title    string
				UserInfo *session.UserInfo
			}{
				Title:    "Hyloblog – Page not found",
				UserInfo: session.ConvertSessionToUserInfoError(sesh),
			},
		},
	).Respond(w, r); err != nil {
		sesh.Println(
			"pathological error:",
			err,
		)
	}
}

func NotFoundSubdomain(w http.ResponseWriter, r *http.Request) {
	sesh, ok := r.Context().Value(session.CtxSessionKey).(*session.Session)
	assert.Assert(ok)
	sesh.Println("404 (subdomain)", r.Host, r.URL)
	w.WriteHeader(http.StatusNotFound)
	sesh.Println("userinfo", session.ConvertSessionToUserInfoError(sesh))
	if err := response.NewTemplate(
		[]string{"404_subdomain.html"},
		util.PageInfo{
			Data: struct {
				Title              string
				UserInfo           *session.UserInfo
				Hyloblog          string
				RequestedSubdomain string
				StartURL           string
			}{
				Title:              "Hyloblog – Site not found",
				UserInfo:           session.ConvertSessionToUserInfoError(sesh),
				Hyloblog:          config.Config.Hyloblog.Hyloblog,
				RequestedSubdomain: r.Host,
				StartURL: fmt.Sprintf(
					"%s://%s",
					config.Config.Hyloblog.Protocol,
					config.Config.Hyloblog.RootDomain,
				),
			},
		},
	).Respond(w, r); err != nil {
		sesh.Println(
			"pathological error:",
			err,
		)
	}
}

func NotFoundDomain(w http.ResponseWriter, r *http.Request) {
	sesh, ok := r.Context().Value(session.CtxSessionKey).(*session.Session)
	assert.Assert(ok)
	sesh.Println("404 (domain)", r.Host, r.URL)
	w.WriteHeader(http.StatusNotFound)
	if err := response.NewTemplate(
		[]string{"404_domain.html"},
		util.PageInfo{
			Data: struct {
				Title           string
				UserInfo        *session.UserInfo
				Hyloblog       string
				RequestedDomain string
				DiscordURL      string
				DomainGuideURL  string
			}{
				Title:           "Hyloblog – Site not found",
				UserInfo:        session.ConvertSessionToUserInfoError(sesh),
				Hyloblog:       config.Config.Hyloblog.Hyloblog,
				RequestedDomain: r.Host,
				DomainGuideURL:  config.Config.Hyloblog.CustomDomainGuideURL,
				DiscordURL:      config.Config.Hyloblog.DiscordURL,
			},
		},
	).Respond(w, r); err != nil {
		sesh.Println(
			"pathological error:",
			err,
		)
	}
}
