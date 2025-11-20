package bootstrap

import (
	"context"
	"log/slog"

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
	"github.com/risk-place-angola/backend-risk-place/internal/infra/aws"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/aws/s3"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/db"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/fcm"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/location"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/logger"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/redis"

	"github.com/risk-place-angola/backend-risk-place/internal/infra/twilio"
)

type Container struct {
	Cfg *config.Config

	UserHandler             *handler.UserHandler
	WSHandler               *websocket.WSHandler
	AlertHandler            *handler.AlertHandler
	ReportHandler           *handler.ReportHandler
	RiskHandler             *handler.RiskHandler
	DeviceHandler           *handler.DeviceHandler
	LocationSharingHandler  *handler.LocationSharingHandler
	SafeRouteHandler        *handler.SafeRouteHandler
	EmergencyContactHandler *handler.EmergencyContactHandler
	MyAlertsHandler         *handler.MyAlertsHandler
	SafetySettingsHandler   *handler.SafetySettingsHandler
	NotificationHandler     *handler.NotificationHandler
	StorageHandler          *handler.StorageHandler

	UserApp *application.Application

	Hub                    *websocket.Hub
	AuthMiddleware         *middleware.AuthMiddleware
	OptionalAuthMiddleware *middleware.OptionalAuthMiddleware
}

func NewContainer() (*Container, error) {
	cfg := config.Load()

	logger.LoggerInit(cfg.AppEnv)

	awsConfig, err := aws.LoadDefaultConfigCredentials(context.Background())
	if err != nil {
		slog.Error("unable to load AWS SDK config", "error", err)
		panic(err)
	}

	cfg.AWSConfig.AwsConfig = *awsConfig

	database := db.NewPostgresConnection(cfg)
	rdb := redis.NewRedis(cfg)
	twilioSMS := twilio.NewTwilio(cfg.TwilioConfig)
	firebaseApp := fcm.NewFirebaseApp(cfg.FirebaseConfig)

	locationStore := location.NewRedisLocationStore(rdb)
	storageService := s3.NewS3StorageService(cfg.AWSConfig)

	userRepoPG := postgres.NewUserRepoPG(database)
	roleRepoPG := postgres.NewRoleRepoPG(database)
	alertRepoPG := postgres.NewAlertRepoPG(database)
	riskTypeRepoPG := postgres.NewRiskTypeRepoPG(database)
	riskTopicRepoPG := postgres.NewRiskTopicRepoPG(database)
	reportRepoPG := postgres.NewReportRepoPG(database, locationStore)
	anonymousSessionRepoPG := postgres.NewAnonymousSessionRepository(database)
	locationSharingRepoPG := postgres.NewLocationSharingRepository(database)
	safeRouteRepoPG := postgres.NewSafeRouteRepoPG(database)
	emergencyContactRepoPG := postgres.NewEmergencyContactRepository(database)
	safetySettingsRepoPG := postgres.NewSafetySettingsRepository(database)
	deviceMappingRepoPG := postgres.NewDeviceUserMappingRepository(database)
	migrationRepoPG := postgres.NewAnonymousMigrationRepository(database)

	emailService := notifier.NewSmtpEmailService(cfg)
	tokenService := service.NewJwtTokenService(cfg)
	hashService := service.NewBcryptHasher()
	geoService := domainService.NewGeolocationService()

	migrationService := service.NewAnonymousMigrationService(
		deviceMappingRepoPG,
		migrationRepoPG,
		anonymousSessionRepoPG,
		alertRepoPG,
		safetySettingsRepoPG,
		locationSharingRepoPG,
	)

	dispatcher := event.NewEventDispatcher()

	hub := websocket.NewHub(locationStore, geoService)
	go hub.Run()

	notifierFCM := notifier.NewFCMNotifier(firebaseApp)
	notifierSMS := notifier.NewSMSNotifier(twilioSMS, cfg.TwilioConfig)

	verificationService := service.NewVerificationService(
		rdb,
		notifierSMS,
		emailService,
		cfg.FrontendURL,
	)

	translationService := service.NewTranslationService()
	reportVerificationService := service.NewReportVerificationService(reportRepoPG)

	eventlistener.RegisterEventListeners(
		dispatcher,
		hub,
		userRepoPG,
		anonymousSessionRepoPG,
		notifierFCM,
		notifierSMS,
		translationService,
	)

	userApp := application.NewUserApplication(
		userRepoPG,
		roleRepoPG,
		alertRepoPG,
		riskTypeRepoPG,
		riskTopicRepoPG,
		reportRepoPG,
		locationSharingRepoPG,
		anonymousSessionRepoPG,
		safeRouteRepoPG,
		emergencyContactRepoPG,
		safetySettingsRepoPG,
		tokenService,
		hashService,
		emailService,
		notifierSMS,
		&cfg,
		locationStore,
		geoService,
		dispatcher,
		migrationService,
		verificationService,
		storageService,
	)

	authMW := middleware.NewAuthMiddleware(cfg)
	optionalAuthMW := middleware.NewOptionalAuthMiddleware(authMW)

	registerDeviceUC := device.NewRegisterDeviceUseCase(anonymousSessionRepoPG)
	updateDeviceLocationUC := device.NewUpdateDeviceLocationUseCase(anonymousSessionRepoPG, locationStore)

	userApp.ReportVerificationService = reportVerificationService

	userHandler := handler.NewUserHandler(userApp)
	alertHandler := handler.NewAlertHandler(userApp)
	wsHandler := websocket.NewWSHandler(hub, *authMW, optionalAuthMW)
	reportHandler := handler.NewReportHandler(userApp, reportRepoPG)
	riskHandler := handler.NewRiskHandler(userApp)
	deviceHandler := handler.NewDeviceHandler(registerDeviceUC, updateDeviceLocationUC)
	locationSharingHandler := handler.NewLocationSharingHandler(userApp)
	safeRouteHandler := handler.NewSafeRouteHandler(userApp)
	emergencyContactHandler := handler.NewEmergencyContactHandler(userApp)
	myAlertsHandler := handler.NewMyAlertsHandler(userApp)
	safetySettingsHandler := handler.NewSafetySettingsHandler(userApp)
	notificationHandler := handler.NewNotificationHandler(userApp)
	storageHandler := handler.NewStorageHandler(storageService, userApp)

	return &Container{
		UserApp:                 userApp,
		UserHandler:             userHandler,
		AuthMiddleware:          authMW,
		OptionalAuthMiddleware:  optionalAuthMW,
		WSHandler:               wsHandler,
		Hub:                     hub,
		Cfg:                     &cfg,
		AlertHandler:            alertHandler,
		ReportHandler:           reportHandler,
		RiskHandler:             riskHandler,
		DeviceHandler:           deviceHandler,
		LocationSharingHandler:  locationSharingHandler,
		SafeRouteHandler:        safeRouteHandler,
		EmergencyContactHandler: emergencyContactHandler,
		MyAlertsHandler:         myAlertsHandler,
		SafetySettingsHandler:   safetySettingsHandler,
		NotificationHandler:     notificationHandler,
		StorageHandler:          storageHandler,
	}, nil
}
