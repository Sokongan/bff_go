package service_oauth

import (
	"context"
	"errors"
	"fmt"
	oauth_domain "sso-bff/internal/domain/oauth"
	"sso-bff/internal/lib"
	"sso-bff/modules/oauth"
	oauth_helper_session "sso-bff/modules/oauth/helper/session"
	oauth_helper_token "sso-bff/modules/oauth/helper/token"

	"time"

	"golang.org/x/oauth2"
)

type SessionService struct {
	oauthClient oauth.OAuthClientPort
	sessions    oauth.SessionStorePort
	idTokens    oauth.IDTokenVerifierPort
	sessionTTL  time.Duration
}

func NewSessionService(
	oauthClient oauth.OAuthClientPort,
	sessions oauth.SessionStorePort,
	idTokens oauth.IDTokenVerifierPort,
	ttl time.Duration,
) *SessionService {

	return &SessionService{
		oauthClient: oauthClient,
		sessions:    sessions,
		idTokens:    idTokens,
		sessionTTL:  ttl,
	}
}
func (s *SessionService) CreateSession(
	ctx context.Context,
	token *oauth2.Token,
	clientToken string,
) (string, time.Duration, error) {

	if s == nil || s.sessions == nil || s.idTokens == nil {
		return "", 0, oauth.ErrServiceMisconfigured
	}
	if token == nil {
		return "", 0, errors.New("token is nil")
	}

	sessionID, err := lib.PKCEStateToken()
	if err != nil {
		return "", 0, fmt.Errorf("generate session id: %w", err)
	}

	idToken := oauth_helper_token.ExtractIDToken(token)
	if idToken == "" {
		return "", 0, oauth.ErrInvalidIDToken
	}

	claims, err := s.idTokens.Verify(ctx, idToken)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", oauth.ErrInvalidIDToken, err)
	}
	if claims.Subject == "" {
		return "", 0, oauth.ErrInvalidIDToken
	}

	ttl := oauth_helper_session.ComputeTTL(
		s.sessionTTL,
		claims.ExpiresAt,
	)

	session := oauth_helper_session.BuildSession(
		token,
		idToken,
		&claims,
		clientToken,
	)

	if err := s.sessions.SaveSession(
		ctx,
		sessionID,
		session,
		ttl,
	); err != nil {
		return "", 0, fmt.Errorf("%w: %v", oauth.ErrSessionStore, err)
	}

	return sessionID, ttl, nil
}

func (s *SessionService) GetSession(
	ctx context.Context,
	sessionID string,
) (oauth.Session, error) {

	if s == nil || s.sessions == nil {
		return oauth.Session{}, oauth.ErrServiceMisconfigured
	}
	return s.sessions.GetSession(ctx, sessionID)
}

func (s *SessionService) SubjectBySessionID(
	ctx context.Context,
	sessionID string,
) (string, error) {
	session, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	return session.Subject, nil
}

func (s *SessionService) DeleteSession(
	ctx context.Context,
	sessionID string,
) error {
	if s == nil || s.sessions == nil {
		return oauth.ErrServiceMisconfigured
	}
	return s.sessions.DeleteSession(ctx, sessionID)
}

func (s *SessionService) RefreshSession(
	ctx context.Context,
	sessionID string,
) (time.Duration, error) {
	if s == nil || s.sessions == nil ||
		s.oauthClient == nil || s.idTokens == nil {
		return 0, oauth.ErrServiceMisconfigured
	}

	session, err := s.sessions.GetSession(ctx, sessionID)
	if err != nil {
		return 0, err
	}
	if session.RefreshToken == "" {
		return 0, errors.New("refresh token missing")
	}

	token, err := s.oauthClient.Refresh(ctx, session.RefreshToken)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", oauth.ErrTokenExchange, err)
	}

	idToken := ""
	if raw := token.Extra("id_token"); raw != nil {
		if val, ok := raw.(string); ok {
			idToken = val
		}
	}

	claims := oauth_domain.IDTokenClaims{
		Subject:   session.Subject,
		ExpiresAt: token.Expiry,
	}

	if idToken != "" {

		verified, err := s.idTokens.Verify(ctx, idToken)
		if err != nil {
			return 0, fmt.Errorf(
				"%w: %v",
				oauth.ErrInvalidIDToken,
				err,
			)
		}

		if verified.Subject == "" {
			return 0, oauth.ErrInvalidIDToken
		}

		claims = verified
		session.IDToken = idToken
		session.Subject = verified.Subject
	}

	session.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		session.RefreshToken = token.RefreshToken
	}
	session.TokenType = token.TokenType
	session.Expiry = claims.ExpiresAt

	ttl := s.sessionTTL
	if !claims.ExpiresAt.IsZero() {
		remaining := time.Until(claims.ExpiresAt)
		if ttl <= 0 || (remaining > 0 && remaining < ttl) {
			ttl = remaining
		}
	}
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	if err := s.sessions.SaveSession(
		ctx,
		sessionID,
		session,
		ttl,
	); err != nil {
		return 0, fmt.Errorf("%w: %v", oauth.ErrSessionStore, err)
	}

	return ttl, nil
}
