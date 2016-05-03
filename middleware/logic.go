package middleware

import (
	"crypto/md5"
	"encoding/hex"
	_ "fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/honeyjonny/undesirereso/database"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func SetDbContext(dbcontext *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("dbcontext", dbcontext)
		context.Next()
	}
}

func CreateSessionForUser(user database.User) (session database.Session, timestamp string) {
	timestamp = timeToString(user.CreatedAt)

	session = database.Session{
		Cookie: getMD5Hash(user.UserName + timestamp),
		UserID: user.ID,
	}

	return session, timestamp
}

func timeToString(times time.Time) string {
	return times.Format(time.RFC3339)
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GetUserFromGinContext(context *gin.Context) (user database.User, err bool) {
	usr, exists := context.Get("user")

	if exists {
		user = usr.(database.User)
		return user, exists
	}

	return database.User{}, false
}

func TruncatePostContent(originalContent string) (retContent string) {
	if len(originalContent) > 144 {
		retContent = originalContent[0:143]
	} else {
		retContent = originalContent
	}

	return retContent
}

func AuthBySession(dbcontext *gorm.DB) gin.HandlerFunc {
	return func(context *gin.Context) {
		cookie, err := context.Cookie("_session")
		var notFound bool

		if err == nil {
			var user database.User

			notFound =
				dbcontext.
					Table("users").
					Joins("inner join sessions on sessions.user_id = users.id").
					Where(&database.Session{Cookie: cookie}).
					Find(&user).
					RecordNotFound()

			if !notFound {
				//fmt.Printf("find user %s\n", user.UserName)

				context.Set("user", user)
			}

		}

		if err != nil || notFound {
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
		}

		context.Next()
	}
}
