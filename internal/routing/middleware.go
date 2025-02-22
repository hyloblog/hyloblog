package routing

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hyloblog/hyloblog/internal/analytics"
	"github.com/hyloblog/hyloblog/internal/app/handler"
	"github.com/hyloblog/hyloblog/internal/app/handler/response"
	"github.com/hyloblog/hyloblog/internal/assert"
	"github.com/hyloblog/hyloblog/internal/config"
	"github.com/hyloblog/hyloblog/internal/model"
	"github.com/hyloblog/hyloblog/internal/routing/internal/usersite"
	"github.com/hyloblog/hyloblog/internal/session"
)

type RoutingService struct {
	store *model.Store
}

func NewRoutingService(s *model.Store) *RoutingService {
	return &RoutingService{store: s}
}

func (s *RoutingService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.tryRoute(w, r, next); err != nil {
			sesh, ok := r.Context().Value(
				session.CtxSessionKey,
			).(*session.Session)
			assert.Assert(ok)
			assert.Assert(sesh != nil)
			switch {
			case errors.Is(err, usersite.ErrPageNotFound):
				handler.NotFound(w, r)
				break
			case errors.Is(err, usersite.ErrUnknownSubdomain):
				handler.NotFoundSubdomain(w, r)
				break
			case errors.Is(err, usersite.ErrUnknownDomain):
				handler.NotFoundDomain(w, r)
				break
			default:
				sesh.Println("unknown host error:", err)
				handler.HandleError(w, r, err)
				break
			}
		}
	})
}

func (s *RoutingService) tryRoute(
	w http.ResponseWriter, r *http.Request, hyloblog http.Handler,
) error {
	site, err := usersite.GetSite(r.Host, s.store)
	if err != nil {
		if errors.Is(err, usersite.ErrIsService) {
			hyloblog.ServeHTTP(w, r)
			return nil
		}
		var redirectErr usersite.RedirectRule
		if errors.As(err, &redirectErr) {
			response.NewRedirect(
				redirectErr.To(),
				http.StatusSeeOther,
			).Respond(w, r)
			return nil
		}
		return fmt.Errorf("get site: %w", err)
	}
	/* site visit is only recorded after checking for email token because
	 * the redirect would cause two visits to be recorded. */
	if site.RecordEmailClick(r.URL, s.store) {
		http.Redirect(
			w, r,
			stripEmailToken(r.URL),
			http.StatusPermanentRedirect,
		)
		return nil
	}
	analytics.NewMixpanelClientWrapper(
		config.Config.Mixpanel.Token,
	).TrackUserSite(r.Host, r.URL.String(), r)
	if err := site.RecordVisit(r.URL.Path, s.store); err != nil {
		return fmt.Errorf("record visit: %w", err)
	}
	filepath, err := site.GetBinding(r.URL.Path, s.store)
	if err != nil {
		return fmt.Errorf("get filepath: %w", err)
	}
	http.ServeFile(w, r, filepath)
	return nil
}

func stripEmailToken(url *url.URL) string {
	u := *url
	u.RawQuery = ""
	return u.String()
}
