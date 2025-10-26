module github.com/Hmbown/driftlock/productized

go 1.24.1

require (
	github.com/Hmbown/driftlock/api-server v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.5.1
	github.com/segmentio/kafka-go v0.4.48
	github.com/stripe/stripe-go/v75 v75.11.0
	golang.org/x/crypto v0.40.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.3
)

replace (
	github.com/Hmbown/driftlock => ..
	github.com/Hmbown/driftlock/api-server => ../api-server
)
