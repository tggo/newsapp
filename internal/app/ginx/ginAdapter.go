package ginx

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinZap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
func GinZap(logger *zap.Logger, timeFormat string, notLogged map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger = logger.WithOptions(zap.WithCaller(false))
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)

		if skipped, ok := notLogged[path]; ok && skipped {
			return
		}

		var additionalFields []zap.Field
		additionalFields = append(additionalFields, zap.Int("status", c.Writer.Status()))
		additionalFields = append(additionalFields, zap.String("method", c.Request.Method))
		additionalFields = append(additionalFields, zap.String("path", path))
		if query != "" {
			additionalFields = append(additionalFields, zap.String("query", query))
		}
		additionalFields = append(additionalFields, zap.String("ip", c.ClientIP()))
		additionalFields = append(additionalFields, zap.String("user-agent", c.Request.UserAgent()))
		// additionalFields = append(additionalFields, zap.String("time", end.Format(timeFormat)))
		additionalFields = append(additionalFields, zap.Duration("latency", latency))

		// TODO: add userID to log if present
		// if val, ok := c.Get(UserObjKey); ok && val != nil {
		// 	u, okUser := val.(*accountModel.User)
		// 	if okUser {
		// 		additionalFields = append(additionalFields, zap.String("user_id", u.ID.String()))
		// 		additionalFields = append(additionalFields, zap.Any("roles", u.Roles))
		// 	}
		// }

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e, additionalFields...)
			}
		} else {
			logger.Info(path, additionalFields...)
		}
	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs, but the stack info is too large.
func RecoveryWithZap(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, errDump := httputil.DumpRequest(c.Request, false)
				if brokenPipe && errDump == nil {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Time("time", time.Now()),
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
