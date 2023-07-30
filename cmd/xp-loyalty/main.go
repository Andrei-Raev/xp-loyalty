package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wintermonth2298/xp-loyalty/docs"
	"github.com/wintermonth2298/xp-loyalty/internal/handler"
	"github.com/wintermonth2298/xp-loyalty/internal/model"
	"github.com/wintermonth2298/xp-loyalty/internal/repository/mongo"
	"github.com/wintermonth2298/xp-loyalty/internal/service"
	"github.com/wintermonth2298/xp-loyalty/pkg/config"
	"github.com/wintermonth2298/xp-loyalty/pkg/mongo_client"
	"github.com/wintermonth2298/xp-loyalty/pkg/server"
)

// @title XP-loyality App API
// @version 1.0
// @description API Server for XP-loyality Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	configPath := path.Join(exPath, "config/config.json")

	var cfgPath string
	flag.StringVar(&cfgPath, "config", configPath, "path to config")
	flag.Parse()

	cfg, err := config.New(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := mongo_client.New(cfg.Mongo.URI, cfg.Mongo.Name)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())

	// images
	imageRepo := mongo.NewImagesRepository(db)
	imageService := service.NewImagesService(imageRepo)
	imageHandler := handler.NewImagesHandler(imageService)

	// awards
	awardsRepo := mongo.NewAwardsRepository(db)

	// user
	userRepo := mongo.NewUserRepository(db)
	userService := service.NewUserService(userRepo, awardsRepo, imageRepo)
	userHandler := handler.NewUserHandler(userService)

	// cards
	cardsRepo := mongo.NewCardsRepository(db)
	cardsService := service.NewCardsStaticService(cardsRepo, awardsRepo)
	cardsHandler := handler.NewCardsStaticHandler(cardsService, userService)

	// admin
	adminRepo := mongo.NewAdminsRepository(db)
	adminService := service.NewAdminsService(adminRepo)

	// auth
	authRepo := mongo.NewCredentialsRepository(db)
	authService := service.NewAuthService(authRepo, cfg.ModeratorUsername, cfg.ModeratorPassword)
	authHandler := handler.NewAuthHandler(authService, userService, adminService)

	// metrics
	metricService := service.NewServiceMetrics()
	metricsHandler := handler.NewMetricsHandler(metricService)

	go func() {
		for {
			users, err := userService.GetAll(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if err := cardsService.UpdateConstCards(context.Background(), users); err != nil {
				log.Fatal(err)
			}
			t, updatedUsers, err := cardsService.UpdateDailyCards(
				context.Background(),
				users,
				cfg.User.DailyCardsNum,
				cfg.User.UniqueGoals,
			)
			if err == model.ErrNoRandomCards {
				time.Sleep(20 * time.Second)
				continue
			}
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range updatedUsers {
				users[i].LastDailyCardsUpdate = t
				err := userService.Update(context.Background(), users[i])
				if err != nil {
					log.Fatal(err)
				}
			}
			time.Sleep(20 * time.Second)
		}
	}()

	router := gin.Default()

	router.Use(
		gin.Recovery(),
		cors.Default(),
	)

	docs.SwaggerInfo.BasePath = "/"

	api := router.Group("/api", metricsHandler.WithMetrics())
	apiUser := router.Group("/api", authHandler.WithAuth(model.RoleUser))
	apiAdmin := router.Group("/api", authHandler.WithAuth(model.RoleAdmin))
	apiModerator := router.Group("/api", authHandler.WithAuth(model.RoleModerator))
	{
		// auth
		api.POST("/auth/sign-in", authHandler.SignIn)
		api.POST("/auth/sign-up-user", authHandler.SignUpUser)
		apiModerator.POST("/auth/sign-up-admin", authHandler.SignUpAdmin)

		// cards
		apiAdmin.GET("/cards", cardsHandler.GetAllStatic)
		apiAdmin.POST("/cards", cardsHandler.CreateStatic)
		apiAdmin.DELETE("/cards", cardsHandler.DeleteStatic)
		apiAdmin.POST("/cards/done", cardsHandler.UpdateCard)
		apiAdmin.GET("/cards/:username", cardsHandler.GetUserCards)
		apiUser.GET("/cards/profile", cardsHandler.GetProfileCards)
		apiUser.POST("/cards/view", cardsHandler.ViewCard)

		// users
		apiAdmin.GET("/users/:username", userHandler.Get)
		apiUser.GET("/users/profile", userHandler.Profile)

		// images
		api.GET("/images/avatar", imageHandler.GetAvatarImages)
		api.GET("/images/prize", imageHandler.GetPrizeImages)
		api.GET("/images/card-background", imageHandler.GetCardsBackgrounds)
		apiAdmin.POST("/images/upload/avatar", imageHandler.UploadAvatarImage)
		apiAdmin.POST("/images/upload/prize", imageHandler.UploadPrizeImage)
		apiAdmin.POST("/images/upload/card-background", imageHandler.UploadCardsBackground)
		// static
		router.Static("/static/images", "./static/images")

		// metrics
		router.GET("/metrics", prometheusHandler())
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server := server.NewServer(router, cfg.ServerPort)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
