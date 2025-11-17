package bootstrap

import (
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/eventlistener"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/handler"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/notifier"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/repository/postgres"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/service"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/websocket"
	"github.com/risk-place-angola/backend-risk-place/internal/application"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/event"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/db"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/fcm"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/location"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/logger"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/redis"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/twilio"
)

type Container struct {
	Cfg *config.Config

	// Handlers
	UserHandler   *handler.UserHandler
	WSHandler     *websocket.WSHandler
	AlertHandler  *handler.AlertHandler
	ReportHandler *handler.ReportHandler
	RiskHandler   *handler.RiskHandler

	UserApp *application.Application

	Hub *websocket.Hub
	// Middlewares
	AuthMiddleware *middleware.AuthMiddleware
}

// NewContainer initializes the application container with all dependencies.
func NewContainer() (*Container, error) {
	cfg := config.Load()

	logger.LoggerInit(cfg.AppEnv)

	// Infra
	database := db.NewPostgresConnection(cfg)
	rdb := redis.NewRedis(cfg)
	twilioSMS := twilio.NewTwilio(cfg.TwilioConfig)
	firebaseApp := fcm.NewFirebaseApp(cfg.FirebaseConfig)

	// Location Store (usado por reports e websocket)
	locationStore := location.NewRedisLocationStore(rdb)

	// Repositories
	userRepoPG := postgres.NewUserRepoPG(database)
	roleRepoPG := postgres.NewRoleRepoPG(database)
	alertRepoPG := postgres.NewAlertRepoPG(database)
	riskTypeRepoPG := postgres.NewRiskTypeRepoPG(database)
	riskTopicRepoPG := postgres.NewRiskTopicRepoPG(database)
	reportRepoPG := postgres.NewReportRepoPG(database, locationStore)

	// Services (Adapters)
	emailService := notifier.NewSmtpEmailService(cfg)
	tokenService := service.NewJwtTokenService(cfg)
	hashService := service.NewBcryptHasher()
	geoService := domainService.NewGeolocationService()

	// Event Dispatcher
	dispatcher := event.NewEventDispatcher()

	hub := websocket.NewHub(locationStore, geoService)
	go hub.Run()

	// Notifier
	notifierFCM := notifier.NewFCMNotifier(firebaseApp)
	notifierSMS := notifier.NewSMSNotifier(twilioSMS, cfg.TwilioConfig)

	eventlistener.RegisterEventListeners(
		dispatcher,
		hub,
		userRepoPG,
		notifierFCM,
		notifierSMS,
	)

	// Application (usecases)
	userApp := application.NewUserApplication(
		userRepoPG,
		roleRepoPG,
		alertRepoPG,
		riskTypeRepoPG,
		riskTopicRepoPG,
		reportRepoPG,
		tokenService,
		hashService,
		emailService,
		&cfg,
		locationStore,
		geoService,
		dispatcher,
	)

	// Middlewares
	authMW := middleware.NewAuthMiddleware(cfg)

	// Handlers
	userHandler := handler.NewUserHandler(userApp)
	alertHandler := handler.NewAlertHandler(userApp)
	wsHandler := websocket.NewWSHandler(hub, *authMW)
	reportHandler := handler.NewReportHandler(userApp)
	riskHandler := handler.NewRiskHandler(userApp)

	return &Container{
		UserApp:        userApp,
		UserHandler:    userHandler,
		AuthMiddleware: authMW,
		WSHandler:      wsHandler,
		Hub:            hub,
		Cfg:            &cfg,
		AlertHandler:   alertHandler,
		ReportHandler:  reportHandler,
		RiskHandler:    riskHandler,
	}, nil
}
