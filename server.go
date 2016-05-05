package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/honeyjonny/sociality/database"
	"github.com/honeyjonny/sociality/middleware"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	_ "time"
)

func main() {

	dbconfig := database.DbConfig{
		Dialect:          "postgres",
		ConnectionString: "user=sadm dbname=social password=ChangeThis sslmode=disable",
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	router.Static("/public/css/", "./public/css/")
	router.Static("/public/img/", "./public/img/")

	router.LoadHTMLGlob("./templates/*")

	dbcontext := dbconfig.CreateConnection()

	router.Use(middleware.SetDbContext(dbcontext))

	authorized := router.Group("/", middleware.AuthBySession(dbcontext))

	authorized.GET("/users", func(c *gin.Context) {

		if _, exists := middleware.GetUserFromGinContext(c); exists {

			var usrDtos []middleware.UserDTO

			dbctx := c.MustGet("dbcontext").(*gorm.DB)

			dbctx.Debug().
				Table("users").
				Select("id as uid, user_name as username, created_at as created").
				Scan(&usrDtos)

			c.HTML(http.StatusOK, "users.tmpl", gin.H{
				"title": "Users page",
				"users": usrDtos,
			})
		}
	})

	authorized.GET("/home", func(c *gin.Context) {

		if user, exists := middleware.GetUserFromGinContext(c); exists {

			var posts []middleware.PostDTO

			dbctx := c.MustGet("dbcontext").(*gorm.DB)

			dbctx.Debug().
				Table("posts").
				Joins("inner join users on users.id = posts.user_id").
				Where("users.id = ?", user.ID).
				Order("posts.created_at desc").
				Select("posts.created_at as created, posts.text as content").
				Scan(&posts)

			c.HTML(http.StatusOK, "home.tmpl", gin.H{
				"title":    "Home page",
				"username": user.UserName,
				"posts":    posts,
			})
		}
	})

	authorized.GET("/friends", func(c *gin.Context) {

		dbctx := c.MustGet("dbcontext").(*gorm.DB)

		id, isquery := c.GetQuery("add")

		user, exists := middleware.GetUserFromGinContext(c)

		if isquery && (len(id) > 0) && exists {

			uid, err := strconv.ParseInt(id, 0, 64)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid query parameter",
				})

				return
			}

			uuid := uint(uid)

			if uuid == user.ID {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "you cannot follow youself",
				})

				return
			}

			if dbctx.
				Table("users").
				Where("users.id = ?", uuid).
				First(&database.User{}).
				RecordNotFound() {

				c.JSON(http.StatusNotFound, gin.H{
					"error": "user is not exist",
				})

				return
			}

			newFollower := database.Follower{
				ObjectID:  user.ID,
				SubjectID: uuid,
			}

			notAlreadyFollow :=
				dbctx.
					Where(&newFollower).
					First(&database.Follower{}).
					RecordNotFound()

			if !notAlreadyFollow {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "you already follow this user",
				})

				return
			}

			dbctx.Create(&newFollower)

			c.Header("Location", "/friends")
			c.JSON(http.StatusSeeOther, gin.H{
				"success": "user added to friends",
			})
		}

		if exists && !isquery {

			var userDtos []middleware.UserDTO

			dbctx.Debug().
				Table("followers").
				Joins("inner join users on users.id = followers.subject_id").
				Where("followers.object_id = ?", user.ID).
				Order("users.created_at desc").
				Select("users.id as uid, users.user_name as username, users.created_at as created").
				Scan(&userDtos)

			c.HTML(http.StatusOK, "friends.tmpl", gin.H{
				"title":   "Friends",
				"friends": userDtos,
			})
		}
	})

	authorized.GET("/messages", func(c *gin.Context) {

		dbctx := c.MustGet("dbcontext").(*gorm.DB)

		id, isquery := c.GetQuery("pm")

		user, exists := middleware.GetUserFromGinContext(c)

		if isquery && (len(id) > 0) && exists {

			var messageDtos []middleware.MessageDTO

			dbctx.Debug().
				Table("messages").
				Where("object_id = ? and subject_id = ?", user.ID, id).
				Order("created_at asc").
				Select("created_at as created, text as message").
				Scan(&messageDtos)

			c.HTML(http.StatusOK, "messages.tmpl", gin.H{
				"title":    "PM",
				"uid":      id,
				"messages": messageDtos,
			})

		} else {
			c.Header("Location", "/friends")
			c.JSON(http.StatusSeeOther, gin.H{
				"error": "got arguments or life!",
			})
		}
	})

	authorized.POST("/messages", func(c *gin.Context) {

		dbctx := c.MustGet("dbcontext").(*gorm.DB)

		if user, exists := middleware.GetUserFromGinContext(c); exists {

			var newMessage middleware.MessageForm

			if c.Bind(&newMessage) == nil {

				value, _ := strconv.ParseUint(newMessage.Uid, 0, 64)

				uuid := uint(value)

				if dbctx.
					Table("users").
					Where("users.id = ?", uuid).
					First(&database.User{}).
					RecordNotFound() {

					c.JSON(http.StatusNotFound, gin.H{
						"error": "user is not exist",
					})

					return
				}

				msg := database.Message{
					ObjectID:  user.ID,
					SubjectID: uuid,
					Text:      newMessage.Message,
				}

				dbctx.Create(&msg)

				href := fmt.Sprintf("/messages?pm=%s", newMessage.Uid)
				c.Header("Location", href)
				c.JSON(http.StatusSeeOther, gin.H{
					"id":      msg.ID,
					"created": msg.CreatedAt,
				})

			} else {

				c.JSON(http.StatusBadRequest, gin.H{
					"error": "form invalid",
				})
			}
		}
	})

	authorized.POST("/posts", func(c *gin.Context) {
		if user, exists := middleware.GetUserFromGinContext(c); exists {

			var newPost middleware.PostForm

			if c.Bind(&newPost) == nil {

				dbctx := c.MustGet("dbcontext").(*gorm.DB)

				truncatedContent := middleware.TruncatePostContent(newPost.Content)

				dbPost := database.Post{
					UserID: user.ID,
					Text:   truncatedContent,
				}

				dbctx.Create(&dbPost)

				c.Header("Location", "/home")
				c.JSON(http.StatusSeeOther, gin.H{
					"created": dbPost,
				})

			} else {

				c.JSON(http.StatusBadRequest, gin.H{
					"error": "form invalid",
				})
			}
		}
	})

	authorized.GET("/logout", func(c *gin.Context) {
		if user, exists := middleware.GetUserFromGinContext(c); exists {

			dbctx := c.MustGet("dbcontext").(*gorm.DB)

			/*			dbctx.
						Table("sessions").
						Where(&database.Session{UserID: user.ID}).
						Delete(database.Session{})*/

			dbctx.
				Unscoped().
				Table("sessions").
				Where(&database.Session{UserID: user.ID}).
				Delete(database.Session{})

			c.Header("Location", "/")
			c.JSON(http.StatusSeeOther, gin.H{
				"logout": user.UserName,
			})
		}
	})

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{
			"title": "Register form",
			"body":  "Register, please",
		})

	})

	router.POST("/register", func(c *gin.Context) {
		var newUser middleware.LoginForm

		dbctx := c.MustGet("dbcontext").(*gorm.DB)

		if c.Bind(&newUser) == nil {

			var checkUsr database.User

			notFound :=
				dbctx.
					Where(&database.User{
						UserName: newUser.Username,
					}).
					First(&checkUsr).
					RecordNotFound()

			if !notFound {

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "user already exists",
				})

			} else {

				dbUser := database.User{
					UserName: newUser.Username,
					Password: newUser.Password,
				}

				dbctx.Create(&dbUser)

				c.Header("Location", "/login")
				c.JSON(http.StatusSeeOther, gin.H{
					"registered": dbUser.UserName,
				})

			}

		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "form invalid",
			})

		}

	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Login form",
			"body":  "Login, please",
		})

	})

	router.POST("/login", func(c *gin.Context) {
		var form middleware.LoginForm
		var dbUser database.User

		dbctx := c.MustGet("dbcontext").(*gorm.DB)

		if c.Bind(&form) == nil {

			notFound :=
				dbctx.
					Where(&database.User{
						UserName: form.Username,
						Password: form.Password,
					}).
					First(&dbUser).
					RecordNotFound()

			if notFound {

				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized",
				})

				return

			} else {

				session := middleware.CreateSessionForUser(dbUser)

				/*				dbctx.
								Table("sessions").
								Where(&database.Session{UserID: dbUser.ID}).
								Delete(database.Session{})*/

				dbctx.
					Unscoped().
					Table("sessions").
					Where(&database.Session{UserID: dbUser.ID}).
					Delete(database.Session{})

				dbctx.
					Create(&session)

				c.SetCookie("_session", session.Cookie, 0, "/", "", false, false)

				c.Header("Location", "/home")
				c.JSON(http.StatusFound, gin.H{
					"logined": dbUser.UserName,
				})

				return
			}

		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "form invalid",
			})

			return
		}
	})

	router.GET("/", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Typical Social Network",
			"body":  "Social network for you and your colleagues!",
		})
	})

	router.Run(":8080")
}
