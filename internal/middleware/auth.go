package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yz13-dev/checkhouse-remote-checker/internal/auth"
)

func getLocalRegion() string {
	region := strings.TrimSpace(os.Getenv("REGION"))
	if region == "" {
		return "default"
	}
	return region
}

func TokenAuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.GetHeader("Authorization")
		if bearer == "" {
			log.Println("Authorization header is missing")
			c.AbortWithStatus(401)
			return
		}

		token := strings.TrimPrefix(bearer, "Bearer ")

		if token == "" {
			log.Println("Token is missing")
			c.AbortWithStatus(401)
			return
		}

		claims, err := jwtService.Verify(token)
		if err != nil {
			log.Println("Token is invalid")
			c.AbortWithStatus(401)
			return
		}
		claimRegion := claims.Region
		localRegion := getLocalRegion()
		// log.Println(claimRegion, localRegion)

		if claimRegion != localRegion {
			log.Println("Region mismatch")
			c.AbortWithStatus(401)
			return
		}

		log.Println("Region matched")

	}
}
