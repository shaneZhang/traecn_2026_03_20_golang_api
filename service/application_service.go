package service

import (
	"errors"
	"fmt"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/dto"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/repository"
	"github.com/casdoor/casdoor/util"
)

type ApplicationService interface {
	GetApplication(id string) (*object.Application, error)
	GetApplicationByClientId(clientId string) (*object.Application, error)
	GetApplicationByUser(user *object.User) (*object.Application, error)
	GetApplicationByUserId(userId string) (*object.Application, error)

	GetApplications(owner string) ([]*object.Application, error)
	GetOrganizationApplications(owner, organization string) ([]*object.Application, error)
	GetPaginationApplications(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, int64, error)
	GetPaginationOrganizationApplications(owner, organization string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, int64, error)

	GetApplicationCount(owner, field, value string) (int64, error)
	GetOrganizationApplicationCount(owner, organization, field, value string) (int64, error)

	CreateApplication(app *object.Application, lang string) (bool, error)
	UpdateApplication(id string, app *object.Application, isGlobalAdmin bool, lang string) (bool, error)
	DeleteApplication(app *object.Application) (bool, error)

	GetMaskedApplication(app *object.Application, userId string) *object.Application
	GetMaskedApplications(apps []*object.Application, userId string) []*object.Application
	GetAllowedApplications(apps []*object.Application, userId, lang string) ([]*object.Application, error)

	CheckIpWhitelist(ipWhitelist, lang string) error
}

type applicationService struct {
	appRepo repository.ApplicationRepository
}

func NewApplicationService(appRepo repository.ApplicationRepository) ApplicationService {
	return &applicationService{
		appRepo: appRepo,
	}
}

func (s *applicationService) GetApplication(id string) (*object.Application, error) {
	return s.appRepo.GetById(id)
}

func (s *applicationService) GetApplicationByClientId(clientId string) (*object.Application, error) {
	return s.appRepo.GetByClientId(clientId)
}

func (s *applicationService) GetApplicationByUser(user *object.User) (*object.Application, error) {
	return s.appRepo.GetByUser(user)
}

func (s *applicationService) GetApplicationByUserId(userId string) (*object.Application, error) {
	return s.appRepo.GetByUserId(userId)
}

func (s *applicationService) GetApplications(owner string) ([]*object.Application, error) {
	return s.appRepo.List(owner)
}

func (s *applicationService) GetOrganizationApplications(owner, organization string) ([]*object.Application, error) {
	return s.appRepo.ListByOrganization(owner, organization)
}

func (s *applicationService) GetPaginationApplications(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, int64, error) {
	if offset == -1 || limit == -1 {
		apps, err := s.appRepo.List(owner)
		if err != nil {
			return nil, 0, err
		}
		return apps, int64(len(apps)), nil
	}

	count, err := s.appRepo.Count(owner, field, value)
	if err != nil {
		return nil, 0, err
	}

	apps, err := s.appRepo.ListPagination(owner, offset, limit, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	return apps, count, nil
}

func (s *applicationService) GetPaginationOrganizationApplications(owner, organization string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, int64, error) {
	if offset == -1 || limit == -1 {
		apps, err := s.appRepo.ListByOrganization(owner, organization)
		if err != nil {
			return nil, 0, err
		}
		return apps, int64(len(apps)), nil
	}

	count, err := s.appRepo.CountByOrganization(owner, organization, field, value)
	if err != nil {
		return nil, 0, err
	}

	apps, err := s.appRepo.ListPaginationByOrganization(owner, organization, offset, limit, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	return apps, count, nil
}

func (s *applicationService) GetApplicationCount(owner, field, value string) (int64, error) {
	return s.appRepo.Count(owner, field, value)
}

func (s *applicationService) GetOrganizationApplicationCount(owner, organization, field, value string) (int64, error) {
	return s.appRepo.CountByOrganization(owner, organization, field, value)
}

func (s *applicationService) CreateApplication(app *object.Application, lang string) (bool, error) {
	count, err := s.appRepo.Count("", "", "")
	if err != nil {
		return false, err
	}

	if err := checkQuotaForApplication(int(count)); err != nil {
		return false, err
	}

	if err := s.CheckIpWhitelist(app.IpWhitelist, lang); err != nil {
		return false, err
	}

	affected, err := s.appRepo.Create(app)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *applicationService) UpdateApplication(id string, app *object.Application, isGlobalAdmin bool, lang string) (bool, error) {
	if err := s.CheckIpWhitelist(app.IpWhitelist, lang); err != nil {
		return false, err
	}

	affected, err := object.UpdateApplication(id, app, isGlobalAdmin, lang)
	if err != nil {
		return false, err
	}

	return affected, nil
}

func (s *applicationService) DeleteApplication(app *object.Application) (bool, error) {
	affected, err := s.appRepo.Delete(app)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *applicationService) GetMaskedApplication(app *object.Application, userId string) *object.Application {
	return object.GetMaskedApplication(app, userId)
}

func (s *applicationService) GetMaskedApplications(apps []*object.Application, userId string) []*object.Application {
	return object.GetMaskedApplications(apps, userId)
}

func (s *applicationService) GetAllowedApplications(apps []*object.Application, userId, lang string) ([]*object.Application, error) {
	return object.GetAllowedApplications(apps, userId, lang)
}

func (s *applicationService) CheckIpWhitelist(ipWhitelist, lang string) error {
	return object.CheckIpWhitelist(ipWhitelist, lang)
}

func checkQuotaForApplication(count int) error {
	quota := conf.GetConfigQuota().Application
	if quota == -1 {
		return nil
	}
	if count >= quota {
		return fmt.Errorf("application quota is exceeded")
	}
	return nil
}

type OAuthService interface {
	HandleOAuthGrant(req *dto.OAuthGrantRequest) error
	HandleOAuthToken(req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error)
}

type oauthService struct {
	appRepo  repository.ApplicationRepository
	userRepo repository.UserRepository
}

func NewOAuthService(appRepo repository.ApplicationRepository, userRepo repository.UserRepository) OAuthService {
	return &oauthService{
		appRepo:  appRepo,
		userRepo: userRepo,
	}
}

func (s *oauthService) HandleOAuthGrant(req *dto.OAuthGrantRequest) error {
	if req.ClientId == "" {
		return errors.New("client_id is required")
	}
	if req.RedirectUri == "" {
		return errors.New("redirect_uri is required")
	}

	app, err := s.appRepo.GetByClientId(req.ClientId)
	if err != nil {
		return err
	}
	if app == nil {
		return errors.New("invalid client_id")
	}

	validRedirect := false
	for _, uri := range app.RedirectUris {
		if uri == req.RedirectUri {
			validRedirect = true
			break
		}
	}
	if !validRedirect {
		return errors.New("invalid redirect_uri")
	}

	return nil
}

func (s *oauthService) HandleOAuthToken(req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error) {
	if req.GrantType == "" {
		return nil, errors.New("grant_type is required")
	}

	switch req.GrantType {
	case "authorization_code":
		return s.handleAuthorizationCodeGrant(req)
	case "refresh_token":
		return s.handleRefreshTokenGrant(req)
	case "password":
		return s.handlePasswordGrant(req)
	case "client_credentials":
		return s.handleClientCredentialsGrant(req)
	default:
		return nil, errors.New("unsupported grant_type")
	}
}

func (s *oauthService) handleAuthorizationCodeGrant(req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *oauthService) handleRefreshTokenGrant(req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *oauthService) handlePasswordGrant(req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *oauthService) handleClientCredentialsGrant(req *dto.OAuthTokenRequest) (*dto.OAuthTokenResponse, error) {
	return nil, errors.New("not implemented")
}

func GetOwnerAndNameFromIdWithError(id string) (string, string, error) {
	return util.GetOwnerAndNameFromIdWithError(id)
}
