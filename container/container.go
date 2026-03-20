package container

import (
	"sync"

	"github.com/casdoor/casdoor/cache"
	"github.com/casdoor/casdoor/repository"
	"github.com/casdoor/casdoor/service"
	"github.com/xorm-io/xorm"
)

type Container struct {
	engine *xorm.Engine
	
	userRepo         repository.UserRepository
	orgRepo          repository.OrganizationRepository
	appRepo          repository.ApplicationRepository
	groupRepo        repository.GroupRepository
	mfaRepo          repository.MfaRepository
	
	userService         service.UserService
	orgService          service.OrganizationService
	groupService        service.GroupService
	appService          service.ApplicationService
	oauthService        service.OAuthService
	mfaService          service.MfaService
	
	userCache        cache.UserCache
	orgCache         cache.OrganizationCache
	appCache         cache.ApplicationCache
}

var (
	instance *Container
	once     sync.Once
)

func GetContainer(engine *xorm.Engine) *Container {
	once.Do(func() {
		instance = &Container{
			engine: engine,
		}
		instance.initRepositories()
		instance.initCaches()
		instance.initServices()
	})
	return instance
}

func (c *Container) initRepositories() {
	c.userRepo = repository.NewUserRepository(c.engine)
	c.orgRepo = repository.NewOrganizationRepository(c.engine)
	c.appRepo = repository.NewApplicationRepository(c.engine)
	c.groupRepo = repository.NewGroupRepository(c.engine)
	c.mfaRepo = repository.NewMfaRepository(c.engine)
}

func (c *Container) initCaches() {
	c.userCache = cache.NewUserCache()
	c.orgCache = cache.NewOrganizationCache()
	c.appCache = cache.NewApplicationCache()
}

func (c *Container) initServices() {
	c.userService = service.NewUserService(c.userRepo, c.orgRepo, c.appRepo)
	c.orgService = service.NewOrganizationService(c.orgRepo, c.groupRepo, c.appRepo)
	c.groupService = service.NewGroupService(c.groupRepo, c.userRepo)
	c.appService = service.NewApplicationService(c.appRepo)
	c.oauthService = service.NewOAuthService(c.appRepo, c.userRepo)
	c.mfaService = service.NewMfaService(c.mfaRepo, c.userRepo, c.orgRepo)
}

func (c *Container) GetUserService() service.UserService {
	return c.userService
}

func (c *Container) GetOrganizationService() service.OrganizationService {
	return c.orgService
}

func (c *Container) GetGroupService() service.GroupService {
	return c.groupService
}

func (c *Container) GetApplicationService() service.ApplicationService {
	return c.appService
}

func (c *Container) GetOAuthService() service.OAuthService {
	return c.oauthService
}

func (c *Container) GetMfaService() service.MfaService {
	return c.mfaService
}

func (c *Container) GetUserCache() cache.UserCache {
	return c.userCache
}

func (c *Container) GetOrganizationCache() cache.OrganizationCache {
	return c.orgCache
}

func (c *Container) GetApplicationCache() cache.ApplicationCache {
	return c.appCache
}

func (c *Container) GetUserRepository() repository.UserRepository {
	return c.userRepo
}

func (c *Container) GetOrganizationRepository() repository.OrganizationRepository {
	return c.orgRepo
}

func (c *Container) GetApplicationRepository() repository.ApplicationRepository {
	return c.appRepo
}

func (c *Container) GetGroupRepository() repository.GroupRepository {
	return c.groupRepo
}

func (c *Container) GetMfaRepository() repository.MfaRepository {
	return c.mfaRepo
}
