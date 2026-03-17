package oauth_helper_token

import "golang.org/x/oauth2"

func ExtractIDToken(token *oauth2.Token) string {
	if token == nil {
		return ""
	}

	raw := token.Extra("id_token")
	if raw == nil {
		return ""
	}

	if val, ok := raw.(string); ok {
		return val
	}

	return ""
}
