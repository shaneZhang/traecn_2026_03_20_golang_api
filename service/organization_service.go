package service

import (
	"fmt"

	"github.com/casdoor/casdoor/conf"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/repository"
)

type OrganizationService interface {
	GetOrganization(id string) (*object.Organization, error)
	GetOrganizationByUser(user *object.User) (*object.Organization, error)

	GetOrganizations(owner string, names ...string) ([]*object.Organization, error)
	GetOrganizationsByFields(owner string, fields ...string) ([]*object.Organization, error)
	GetPaginationOrganizations(owner, name string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Organization, int64, error)

	GetOrganizationCount(owner, name, field, value string) (int64, error)

	CreateOrganization(org *object.Organization, lang string) (bool, error)
	UpdateOrganization(id string, org *object.Organization, isGlobalAdmin bool, lang string) (bool, error)
	DeleteOrganization(org *object.Organization) (bool, error)

	GetDefaultApplication(id string) (*object.Application, error)
	CheckIpWhitelist(ipWhitelist, lang string) error

	GetMaskedOrganization(org *object.Organization) (*object.Organization, error)
	GetMaskedOrganizations(orgs []*object.Organization) ([]*object.Organization, error)
}

type GroupService interface {
	GetGroup(id string) (*object.Group, error)

	GetGroups(owner string) ([]*object.Group, error)
	GetGlobalGroups() ([]*object.Group, error)
	GetPaginationGroups(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Group, int64, error)

	GetGroupCount(owner, field, value string) (int64, error)

	CreateGroup(group *object.Group) (bool, error)
	CreateGroups(groups []*object.Group) (bool, error)
	UpdateGroup(id string, group *object.Group) (bool, error)
	DeleteGroup(group *object.Group) (bool, error)

	ExtendGroupWithUsers(group *object.Group) error
	ExtendGroupsWithUsers(groups []*object.Group) error
	ConvertToTreeData(groups []*object.Group, parentId string) []*object.Group
}

type organizationService struct {
	orgRepo   repository.OrganizationRepository
	groupRepo repository.GroupRepository
	appRepo   repository.ApplicationRepository
}

func NewOrganizationService(orgRepo repository.OrganizationRepository, groupRepo repository.GroupRepository, appRepo repository.ApplicationRepository) OrganizationService {
	return &organizationService{
		orgRepo:   orgRepo,
		groupRepo: groupRepo,
		appRepo:   appRepo,
	}
}

func (s *organizationService) GetOrganization(id string) (*object.Organization, error) {
	return s.orgRepo.GetById(id)
}

func (s *organizationService) GetOrganizationByUser(user *object.User) (*object.Organization, error) {
	return s.orgRepo.GetByUser(user)
}

func (s *organizationService) GetOrganizations(owner string, names ...string) ([]*object.Organization, error) {
	return s.orgRepo.List(owner, names...)
}

func (s *organizationService) GetOrganizationsByFields(owner string, fields ...string) ([]*object.Organization, error) {
	return s.orgRepo.ListByFields(owner, fields...)
}

func (s *organizationService) GetPaginationOrganizations(owner, name string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Organization, int64, error) {
	if offset == -1 || limit == -1 {
		orgs, err := s.orgRepo.List(owner)
		if err != nil {
			return nil, 0, err
		}
		return orgs, int64(len(orgs)), nil
	}

	count, err := s.orgRepo.Count(owner, name, field, value)
	if err != nil {
		return nil, 0, err
	}

	orgs, err := s.orgRepo.ListPagination(owner, name, offset, limit, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	return orgs, count, nil
}

func (s *organizationService) GetOrganizationCount(owner, name, field, value string) (int64, error) {
	return s.orgRepo.Count(owner, name, field, value)
}

func (s *organizationService) CreateOrganization(org *object.Organization, lang string) (bool, error) {
	count, err := s.orgRepo.Count("", "", "", "")
	if err != nil {
		return false, err
	}

	if err := checkQuotaForOrganization(int(count)); err != nil {
		return false, err
	}

	if err := s.CheckIpWhitelist(org.IpWhitelist, lang); err != nil {
		return false, err
	}

	if org.BalanceCurrency == "" {
		org.BalanceCurrency = "USD"
	}

	affected, err := s.orgRepo.Create(org)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *organizationService) UpdateOrganization(id string, org *object.Organization, isGlobalAdmin bool, lang string) (bool, error) {
	if err := s.CheckIpWhitelist(org.IpWhitelist, lang); err != nil {
		return false, err
	}

	if org.BalanceCurrency == "" {
		org.BalanceCurrency = "USD"
	}

	affected, err := s.orgRepo.Update(id, org, isGlobalAdmin)
	if err != nil {
		return false, err
	}

	return affected > 0, nil
}

func (s *organizationService) DeleteOrganization(org *object.Organization) (bool, error) {
	affected, err := s.orgRepo.Delete(org)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *organizationService) GetDefaultApplication(id string) (*object.Application, error) {
	return s.appRepo.GetDefaultByOrganization(id)
}

func (s *organizationService) CheckIpWhitelist(ipWhitelist, lang string) error {
	return object.CheckIpWhitelist(ipWhitelist, lang)
}

func (s *organizationService) GetMaskedOrganization(org *object.Organization) (*object.Organization, error) {
	return object.GetMaskedOrganization(org)
}

func (s *organizationService) GetMaskedOrganizations(orgs []*object.Organization) ([]*object.Organization, error) {
	return object.GetMaskedOrganizations(orgs)
}

type groupService struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

func NewGroupService(groupRepo repository.GroupRepository, userRepo repository.UserRepository) GroupService {
	return &groupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (s *groupService) GetGroup(id string) (*object.Group, error) {
	return s.groupRepo.GetById(id)
}

func (s *groupService) GetGroups(owner string) ([]*object.Group, error) {
	return s.groupRepo.List(owner)
}

func (s *groupService) GetGlobalGroups() ([]*object.Group, error) {
	return s.groupRepo.ListGlobal()
}

func (s *groupService) GetPaginationGroups(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Group, int64, error) {
	if offset == -1 || limit == -1 {
		groups, err := s.groupRepo.List(owner)
		if err != nil {
			return nil, 0, err
		}
		return groups, int64(len(groups)), nil
	}

	count, err := s.groupRepo.Count(owner, field, value)
	if err != nil {
		return nil, 0, err
	}

	groups, err := s.groupRepo.ListPagination(owner, offset, limit, field, value, sortField, sortOrder)
	if err != nil {
		return nil, 0, err
	}

	return groups, count, nil
}

func (s *groupService) GetGroupCount(owner, field, value string) (int64, error) {
	return s.groupRepo.Count(owner, field, value)
}

func (s *groupService) CreateGroup(group *object.Group) (bool, error) {
	affected, err := s.groupRepo.Create(group)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *groupService) CreateGroups(groups []*object.Group) (bool, error) {
	affected, err := s.groupRepo.CreateBatch(groups)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *groupService) UpdateGroup(id string, group *object.Group) (bool, error) {
	affected, err := s.groupRepo.Update(id, group)
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *groupService) DeleteGroup(group *object.Group) (bool, error) {
	return object.DeleteGroup(group)
}

func (s *groupService) ExtendGroupWithUsers(group *object.Group) error {
	return object.ExtendGroupWithUsers(group)
}

func (s *groupService) ExtendGroupsWithUsers(groups []*object.Group) error {
	return object.ExtendGroupsWithUsers(groups)
}

func (s *groupService) ConvertToTreeData(groups []*object.Group, parentId string) []*object.Group {
	return object.ConvertToTreeData(groups, parentId)
}

func checkQuotaForOrganization(count int) error {
	quota := conf.GetConfigQuota().Organization
	if quota == -1 {
		return nil
	}
	if count >= quota {
		return fmt.Errorf("organization quota is exceeded")
	}
	return nil
}
