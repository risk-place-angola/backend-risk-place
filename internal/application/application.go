package application

import (
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/alert"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/report"
	"github.com/risk-place-angola/backend-risk-place/internal/application/usecase/user"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type Application struct {
	UserUseCase   *user.UserUseCase
	AlertUseCase  *alert.AlertUseCase
	ReportUseCase *report.ReportUseCase
}

func NewUserApplication(
	userRepo domainrepository.UserRepository,
	roleRepo domainrepository.RoleRepository,
	alertRepo domainrepository.AlertRepository,
	riskTypeRepo domainrepository.RiskTypesRepository,
	reportRepo domainrepository.ReportRepository,

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
	}
}
