package middlewares

import (
	"net/http"
	"pasteGo/backend/api/rest/v1/handlers"
	auth "pasteGo/backend/api/rest/v1/handlers"
	"pasteGo/backend/api/rest/v1/types"
	"time"

	"github.com/gin-gonic/gin"
)

// ! dev headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie(types.CookieAccessToken)
		if err != nil {
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrGetCookies,
				Explanation: types.ErrGetCookiesExp,
			})
			c.Abort()
			return
		}

		claims, err := auth.ParseClaims(accessToken)
		if err != nil {
			handlers.DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrJWTProcessing,
				Explanation: types.ErrJWTProcessingExp,
			})
			c.Abort()
			return
		}
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			handlers.DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrJWTExpired,
				Explanation: types.ErrJWTExpiredExp,
			})
			c.Abort()
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}

func JwtRefreshMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie(types.CookieRefreshToken)
		if err != nil {
			handlers.DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrGetCookies,
				Explanation: types.ErrGetCookiesExp,
			})
			c.Abort()
			return
		}

		claims, err := auth.ParseClaims(refreshToken)
		if err != nil {
			handlers.DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrJWTProcessing,
				Explanation: types.ErrJWTProcessingExp,
			})
			c.Abort()
			return
		}
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			handlers.DumpCookies(c)
			c.IndentedJSON(http.StatusUnauthorized, types.APIResponse{
				Code:        types.ErrJWTExpired,
				Explanation: types.ErrJWTExpiredExp,
			})
			c.Abort()
			return
		}

		c.Set("userClaims", claims)
		c.Next()
	}
}
