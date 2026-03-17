package oauth_helper_session

import "time"

func ComputeTTL(defaultTTL time.Duration, expiry time.Time) time.Duration {
	ttl := defaultTTL

	if !expiry.IsZero() {
		remaining := time.Until(expiry)

		if ttl <= 0 || (remaining > 0 && remaining < ttl) {
			ttl = remaining
		}
	}

	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	return ttl
}
