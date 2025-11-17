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
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/device"
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

	UserHandler   *handler.UserHandler
	WSHandler     *websocket.WSHandler
	AlertHandler  *handler.AlertHandler
	ReportHandler *handler.ReportHandler
	RiskHandler   *handler.RiskHandler
	DeviceHandler *handler.DeviceHandler

	UserApp *application.Application

	Hub            *websocket.Hub
	AuthMiddleware *middleware.AuthMiddleware
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	logger.LoggerInit(cfg.AppEnv)

	database := db.NewPostgresConnection(cfg)
	rdb := redis.NewRedis(cfg)
	twilioSMS := twilio.NewTwilio(cfg.TwilioConfig)
	firebaseApp := fcm.NewFirebaseApp(cfg.FirebaseConfig)

	locationStore := location.NewRedisLocationStore(rdb)

	userRepoPG := postgres.NewUserRepoPG(database)
	roleRepoPG := postgres.NewRoleRepoPG(database)
	alertRepoPG := postgres.NewAlertRepoPG(database)
	riskTypeRepoPG := postgres.NewRiskTypeRepoPG(database)
	riskTopicRepoPG := postgres.NewRiskTopicRepoPG(database)
	reportRepoPG := postgres.NewReportRepoPG(database, locationStore)
	anonymousSessionRepoPG := postgres.NewAnonymousSessionRepository(database)

	emailService := notifier.NewSmtpEmailService(cfg)
	tokenService := service.NewJwtTokenService(cfg)
	hashService := service.NewBcryptHasher()
	geoService := domainService.NewGeolocationService()

	dispatcher := event.NewEventDispatcher()

	hub := websocket.NewHub(locationStore, geoService)
	go hub.Run()

	notifierFCM := notifier.NewFCMNotifier(firebaseApp)
	notifierSMS := notifier.NewSMSNotifier(twilioSMS, cfg.TwilioConfig)

	eventlistener.RegisterEventListeners(
		dispatcher,
		hub,
		userRepoPG,
		anonymousSessionRepoPG,
		notifierFCM,
		notifierSMS,
	)

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

	authMW := middleware.NewAuthMiddleware(cfg)
	optionalAuthMW := middleware.NewOptionalAuthMiddleware(authMW)

	registerDeviceUC := device.NewRegisterDeviceUseCase(anonymousSessionRepoPG)
	updateDeviceLocationUC := device.NewUpdateDeviceLocationUseCase(anonymousSessionRepoPG, locationStore)

	userHandler := handler.NewUserHandler(userApp)
	alertHandler := handler.NewAlertHandler(userApp)
	wsHandler := websocket.NewWSHandler(hub, *authMW, optionalAuthMW)
	reportHandler := handler.NewReportHandler(userApp)
	riskHandler := handler.NewRiskHandler(userApp)
	deviceHandler := handler.NewDeviceHandler(registerDeviceUC, updateDeviceLocationUC)

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
		DeviceHandler:  deviceHandler,
	}, nil
}
