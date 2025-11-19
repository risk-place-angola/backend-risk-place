package application

import (
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/alert"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/locationsharing"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/report"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/risk"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/saferoute"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/user"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type Application struct {
	UserUseCase            *user.UserUseCase
	AlertUseCase           *alert.AlertUseCase
	ReportUseCase          *report.ReportUseCase
	RiskUseCase            *risk.RiskUseCase
	LocationSharingUseCase *locationsharing.LocationSharingUseCase
	SafeRouteUseCase       *saferoute.SafeRouteUseCase
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

	token port.TokenGenerator,
	hasher port.PasswordHasher,
	emailService port.EmailService,
	config *config.Config,
	locationStore port.LocationStore,
	geoService port.GeolocationService,
	eventDispatcher port.EventDispatcher,
) *Application {
	return &Application{
		UserUseCase: user.NewUserUseCase(
			userRepo,
			roleRepo,
			token,
			hasher,
			emailService,
			config,
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
	}
}
