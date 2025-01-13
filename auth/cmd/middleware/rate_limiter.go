package middlewares

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/juxue97/auth/internal/cache"
	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/common"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID and IP
		ip := r.RemoteAddr // Middleware.RealIP ensures this contains the client's IP
		limit := config.Configs.RateLimit.Limit
		window := config.Configs.RateLimit.Window
		// Apply rate limiting
		allowed, err := cache.CacheStorage.RateLimiter.Count(r.Context(), ip, limit, window)
		if err != nil {
			common.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			// Get the TTL of the rate limit key
			ttl, err := cache.CacheStorage.RateLimiter.GetRemainTime(r.Context(), ip)
			if err != nil {
				common.InternalServerError(w, r, err)
				return
			}
			// Format the retry-after duration
			retryAfter := int(ttl.Seconds())
			fmt.Println(retryAfter)
			if retryAfter < 0 {
				retryAfter = int(window.Seconds()) // Fallback to default window if TTL isn't available
			}

			// convert the retryAfter to string
			retryAfterStr := strconv.Itoa(retryAfter)
			common.TooManyRequestsError(w, r, retryAfterStr)
			return
		}

		// Proceed to the next middleware/handler
		next.ServeHTTP(w, r)
	})
}
