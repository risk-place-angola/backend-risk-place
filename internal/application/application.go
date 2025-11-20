package application

import (
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/alert"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/emergencycontact"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/locationsharing"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/myalerts"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/report"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/risk"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/saferoute"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/safetysettings"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/user"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type Application struct {
	UserUseCase               *user.UserUseCase
	AlertUseCase              *alert.AlertUseCase
	ReportUseCase             *report.ReportUseCase
	RiskUseCase               *risk.RiskUseCase
	LocationSharingUseCase    *locationsharing.LocationSharingUseCase
	SafeRouteUseCase          *saferoute.SafeRouteUseCase
	EmergencyContactUseCase   *emergencycontact.EmergencyContactUseCase
	EmergencyAlertUseCase     *emergencycontact.EmergencyAlertUseCase
	MyAlertsUseCase           *myalerts.MyAlertsUseCase
	SafetySettingsUseCase     *safetysettings.SafetySettingsUseCase
	ReportVerificationService domainService.ReportVerificationService
}

func NewUserApplication(
	userRepo domainrepository.UserRepository,
	roleRepo domainrepository.RoleRepository,
	alertRepo domainrepository.AlertRepository,
	riskTypeRepo domainrepository.RiskTypesRepository,
	riskTopicRepo domainrepository.RiskTopicsRepository,
	reportRepo domainrepository.ReportRepository,
	locationSharingRepo domainrepository.LocationSharingRepository,
	anonymousSessionRepo domainrepository.AnonymousSessionRepository,
	safeRouteRepo domainrepository.SafeRouteRepository,
	emergencyContactRepo domainrepository.EmergencyContactRepository,
	safetySettingsRepo domainrepository.SafetySettingsRepository,

	token port.TokenGenerator,
	hasher port.PasswordHasher,
	emailService port.EmailService,
	smsNotifier port.NotifierSMSService,
	config *config.Config,
	locationStore port.LocationStore,
	geoService port.GeolocationService,
	eventDispatcher port.EventDispatcher,
	migrationService domainService.AnonymousMigrationService,
	verificationService domainService.VerificationService,
	storageService port.StorageService,
) *Application {
	return &Application{
		UserUseCase: user.NewUserUseCase(
			userRepo,
			roleRepo,
			token,
			hasher,
			config,
			migrationService,
			verificationService,
		),
		AlertUseCase: alert.NewAlertUseCase(
			locationStore,
			geoService,
			alertRepo,
			riskTypeRepo,
			eventDispatcher,
		),
		ReportUseCase: report.NewReportUseCase(
			reportRepo,
			eventDispatcher,
			geoService,
			riskTypeRepo,
			locationStore,
		),
		RiskUseCase: risk.NewRiskUseCase(
			riskTypeRepo,
			riskTopicRepo,
			storageService,
		),
		LocationSharingUseCase: locationsharing.NewLocationSharingUseCase(
			locationSharingRepo,
			userRepo,
			anonymousSessionRepo,
			geoService,
			config,
		),
		SafeRouteUseCase: saferoute.NewSafeRouteUseCase(
			safeRouteRepo,
			userRepo,
		),
		EmergencyContactUseCase: emergencycontact.NewEmergencyContactUseCase(
			emergencyContactRepo,
		),
		EmergencyAlertUseCase: emergencycontact.NewEmergencyAlertUseCase(
			emergencyContactRepo,
			userRepo,
			smsNotifier,
		),
		MyAlertsUseCase: myalerts.NewMyAlertsUseCase(
			alertRepo,
			riskTypeRepo,
			riskTopicRepo,
		),
		SafetySettingsUseCase: safetysettings.NewSafetySettingsUseCase(
			safetySettingsRepo,
		),
	}
}
