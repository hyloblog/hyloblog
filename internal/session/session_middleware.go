package session

import (
	"context"
	"net/http"
	"time"

	"github.com/hyloblog/hyloblog/internal/model"
)

const (
	CtxSessionKey = "session"
)

type SessionService struct {
	store *model.Store
}

func NewSessionService(s *model.Store) *SessionService {
	return &SessionService{
		store: s,
	}
}

func (s *SessionService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
		 * XXX: The errors in this function must all be fatal, i.e.
		 * internal server error with no content, until we integrate
		 * tightly with our handler package.
		 */

		logger := newAnonymousLoggerFromRequest(r)

		cookie, err := r.Cookie(CookieName)
		if err != nil {
			logger.Printf("Error getting cookie: %v", err)
			/* create unauth session */
			session, err := CreateUnauthSession(
				s.store, w, unauthSessionDuration, logger,
			)
			if err != nil {
				logger.Printf("Error creating unauth session: %v", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), CtxSessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		/* cookie exists retrieve session */
		session, err := GetSession(s.store, w, cookie.Value, logger)
		if err != nil {
			logger.Printf("Error getting session: %v", err)
			/* expire cookie if error */
			http.SetCookie(w, &http.Cookie{
				Name:     CookieName,
				Value:    "",
				Expires:  time.Now().Add(-1 * time.Hour),
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		if session.expiresAt.Before(time.Now()) {
			logger.Println("session expired")
			/* expire cookie if session expired */
			http.SetCookie(w, &http.Cookie{
				Name:     CookieName,
				Value:    "",
				Expires:  time.Now().Add(-1 * time.Hour),
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		ctx := context.WithValue(r.Context(), CtxSessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
