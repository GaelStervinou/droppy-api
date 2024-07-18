package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "go-api/docs"
	"go-api/internal/http/controllers"
	"go-api/internal/http/middlewares"
	"go-api/internal/storage/postgres"
	"go-api/pkg/environment"
	"log"
	"os"
)

// @title Droppy API
// @version 1.0
// @description This is the API documentation for Droppy

// @contact.name   Droppy API Support
// @contact.email  stervinou.g36@gmail.com
// @host https://droppy.gael-stervinou.fr
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	if env, ok := os.LookupEnv("ENV"); ok {
		environment.SetEnv(env)
	}

	fmt.Println("ENV: " + environment.GetEnv())

	postgres.Init()
	postgres.AutoMigrate()
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AddAllowHeaders("Authorization")
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
	config.AllowOrigins = []string{"https://droppy-420013.web.app"}
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	r.Use(cors.New(config))
	r.Use(gin.Recovery())

	v1 := r.Group("/")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/refresh", controllers.RefreshToken)
			auth.POST("/", controllers.Login)
			auth.POST("", controllers.Login)
			auth.POST("/oauth_token", controllers.FirebaseLogin)
		}

		user := v1.Group("/users")
		{
			user.GET("/:id", middlewares.CurrentUserMiddleware(true), controllers.GetUserById)
			user.POST("/", controllers.Create)
			user.POST("", controllers.Create)
			user.GET("/search", controllers.SearchUsers)
			user.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchUserById)
			user.GET("/my-feed/ws", middlewares.CurrentUserMiddleware(true), controllers.GetCurrentUserFeedWS)
			user.GET("/my-feed", middlewares.CurrentUserMiddleware(true), controllers.GetCurrentUserFeed)
			user.GET("/:id/drops", middlewares.CurrentUserMiddleware(true), controllers.DropsByUserId)

			user.GET("/:id/following", middlewares.CurrentUserMiddleware(true), controllers.GetUserFollowing)
			user.GET("/:id/followers", middlewares.CurrentUserMiddleware(true), controllers.GetUserFollowers)
		}

		follow := v1.Group("/follows")
		{
			follow.POST("/", middlewares.CurrentUserMiddleware(true), controllers.FollowUser)
			follow.GET("/pending", middlewares.CurrentUserMiddleware(true), controllers.GetMyPendingRequestsWS)
			follow.POST("/accept/:id", middlewares.CurrentUserMiddleware(true), controllers.AcceptRequest)
			follow.POST("/reject/:id", middlewares.CurrentUserMiddleware(true), controllers.RejectRequest)
			follow.DELETE("/:id", middlewares.CurrentUserMiddleware(true), controllers.DeleteFollow)
		}

		drop := v1.Group("/drops")
		{
			drop.POST("/", middlewares.CurrentUserMiddleware(true), controllers.CreateDrop)
			drop.GET("/has-user-dropped", middlewares.CurrentUserMiddleware(true), controllers.HasUserDroppedTodayWS)
			drop.GET("/:id", middlewares.CurrentUserMiddleware(true), controllers.GetOneDrop)
			drop.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchDrop)
			drop.DELETE("/:id", middlewares.CurrentUserMiddleware(true), controllers.DeleteDrop)
			drop.POST("/:id/comments", middlewares.CurrentUserMiddleware(true), controllers.CommentDrop)
			drop.POST("/:id/like", middlewares.CurrentUserMiddleware(true), controllers.LikeDrop)
			drop.DELETE("/:id/like", middlewares.CurrentUserMiddleware(true), controllers.UnlikeDrop)
		}

		content := v1.Group("/contents")
		{
			content.GET("/search", middlewares.CurrentUserMiddleware(true), controllers.SearchContentForCurrentDrop)
		}

		report := v1.Group("/reports")
		{
			report.POST("/", middlewares.CurrentUserMiddleware(true), controllers.CreateReport)
		}

		group := v1.Group("/groups")
		{
			group.POST("", middlewares.CurrentUserMiddleware(true), controllers.CreateGroup)
			group.POST("/", middlewares.CurrentUserMiddleware(true), controllers.CreateGroup)
			group.GET("/:id", middlewares.CurrentUserMiddleware(true), controllers.GetOneGroup)
			group.GET("/:id/feed", middlewares.CurrentUserMiddleware(true), controllers.GetGroupFeed)
			group.PATCH("/:id", middlewares.CurrentUserMiddleware(true), controllers.PatchGroup)
			group.GET("/search", middlewares.CurrentUserMiddleware(true), controllers.SearchGroups)
			group.DELETE("/:id", middlewares.CurrentUserMiddleware(true), controllers.DeleteGroup)

			group.POST("/members/:id/join", middlewares.CurrentUserMiddleware(true), controllers.JoinGroup)
			group.POST("/members/:id/:userId", middlewares.CurrentUserMiddleware(true), controllers.AddUserToGroup)
			group.PATCH("/members/:groupId/:memberId", middlewares.CurrentUserMiddleware(true), controllers.PatchGroupMember)
			group.DELETE("/members/:groupId/:memberId", middlewares.CurrentUserMiddleware(true), controllers.DeleteGroupMember)
		}

		comment := v1.Group("/comments")
		{
			comment.POST("/:id/responses", middlewares.CurrentUserMiddleware(true), controllers.RespondToComment)
			comment.DELETE("/:id/responses/:responseId", middlewares.CurrentUserMiddleware(true), controllers.DeleteCommentResponse)
		}

		if environment.IsDev() {
			fixtures := v1.Group("/fixtures")
			{
				fixtures.GET("/all", controllers.PopulateAll)
				fixtures.GET("/users", controllers.PopulateUsers)
				fixtures.GET("/follows", controllers.PopulateFollows)
				fixtures.GET("/drops", controllers.PopulateDrops)
				fixtures.GET("/comments", controllers.PopulateComments)
				fixtures.GET("/groups", controllers.PopulateGroups)
			}
		}

		admin := v1.Group("/admin")
		{
			admin.GET("/users", middlewares.AdminRequired(), controllers.GetAllUsers)
			admin.GET("/users/search", middlewares.AdminRequired(), controllers.AdminSearchUser)
			admin.GET("/users/count", middlewares.AdminRequired(), controllers.GetAllUsersCount)
			admin.PUT("/users/:id", middlewares.AdminRequired(), controllers.UpdateUser)
			admin.GET("/groups", middlewares.AdminRequired(), controllers.GetAllGroups)
			admin.GET("/groups/count", middlewares.AdminRequired(), controllers.GetAllGroupsCount)
			admin.DELETE("/groups/:id", middlewares.AdminRequired(), controllers.AdminDeleteGroup)
			admin.GET("/drops", middlewares.AdminRequired(), controllers.GetAllDrops)
			admin.GET("/drops/count", middlewares.AdminRequired(), controllers.GetAllDropsCount)
			admin.DELETE("/drops/:id", middlewares.AdminRequired(), controllers.AdminDeleteDrop)
			admin.GET("/comments", middlewares.AdminRequired(), controllers.GetAllComments)
			admin.DELETE("/comments/:id", middlewares.AdminRequired(), controllers.DeleteComment)
			admin.GET("/reports", middlewares.AdminRequired(), controllers.GetAllReports)
			admin.PUT("/reports/:id", middlewares.AdminRequired(), controllers.AdminManageReport)
			admin.POST("/drops/schedule", middlewares.AdminRequired(), controllers.AdminScheduleDrop)
			admin.POST("/drops/send-now", middlewares.AdminRequired(), controllers.AdminSendDropNow)
		}
	}
	if environment.IsDev() {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	r.Static("/assets", "./assets")

	/*c := cron.New()
	_, err = c.AddFunc("0 0 * * *", drop_notif.GenerateRandomNotification)

	if err != nil {
		log2.Errorf("Error adding cron job: %v", err)
	}

	c.Start()
	fmt.Println("Scheduler started...")*/

	err = r.Run(":3000")

	if err != nil {
		log.Fatal(err)
	}
}
