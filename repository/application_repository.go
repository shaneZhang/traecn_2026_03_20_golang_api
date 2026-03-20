package repository

import (
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/core"
	"github.com/xorm-io/xorm"
)

type ApplicationRepository interface {
	GetById(id string) (*object.Application, error)
	GetByOwnerAndName(owner, name string) (*object.Application, error)
	GetByClientId(clientId string) (*object.Application, error)
	GetByOrganizationName(organization string) (*object.Application, error)
	GetByUser(user *object.User) (*object.Application, error)
	GetByUserId(userId string) (*object.Application, error)
	GetDefaultByOrganization(id string) (*object.Application, error)

	List(owner string) ([]*object.Application, error)
	ListByOrganization(owner, organization string) ([]*object.Application, error)
	ListPagination(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, error)
	ListPaginationByOrganization(owner, organization string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, error)

	Count(owner, field, value string) (int64, error)
	CountByOrganization(owner, organization, field, value string) (int64, error)

	Create(app *object.Application) (int64, error)
	Update(id string, app *object.Application) (int64, error)
	Delete(app *object.Application) (int64, error)

	Exists(owner, name string) (bool, error)
	ExistsByClientId(clientId string) (bool, error)
}

type applicationRepository struct {
	engine *xorm.Engine
}

func NewApplicationRepository(engine *xorm.Engine) ApplicationRepository {
	return &applicationRepository{engine: engine}
}

func (r *applicationRepository) GetById(id string) (*object.Application, error) {
	return object.GetApplication(id)
}

func (r *applicationRepository) GetByOwnerAndName(owner, name string) (*object.Application, error) {
	if owner == "" || name == "" {
		return nil, nil
	}
	return object.GetApplication(util.GetId(owner, name))
}

func (r *applicationRepository) GetByClientId(clientId string) (*object.Application, error) {
	return object.GetApplicationByClientId(clientId)
}

func (r *applicationRepository) GetByOrganizationName(organization string) (*object.Application, error) {
	return object.GetApplicationByOrganizationName(organization)
}

func (r *applicationRepository) GetByUser(user *object.User) (*object.Application, error) {
	return object.GetApplicationByUser(user)
}

func (r *applicationRepository) GetByUserId(userId string) (*object.Application, error) {
	return object.GetApplicationByUserId(userId)
}

func (r *applicationRepository) GetDefaultByOrganization(id string) (*object.Application, error) {
	return object.GetDefaultApplication(id)
}

func (r *applicationRepository) List(owner string) ([]*object.Application, error) {
	return object.GetApplications(owner)
}

func (r *applicationRepository) ListByOrganization(owner, organization string) ([]*object.Application, error) {
	return object.GetOrganizationApplications(owner, organization)
}

func (r *applicationRepository) ListPagination(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, error) {
	applications := []*object.Application{}
	session := r.engine.NewSession()
	defer session.Close()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	if field != "" && value != "" {
		session.And(builder.Like{util.SnakeString(field), value})
	}

	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}

	orderQuery := util.SnakeString(sortField)
	if sortOrder == "ascend" {
		session.Asc(orderQuery)
	} else {
		session.Desc(orderQuery)
	}

	err := session.Find(&applications)
	return applications, err
}

func (r *applicationRepository) ListPaginationByOrganization(owner, organization string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Application, error) {
	applications := []*object.Application{}
	session := r.engine.NewSession()
	defer session.Close()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	session.And("organization = ? OR is_shared = ?", organization, true)

	if field != "" && value != "" {
		session.And(builder.Like{util.SnakeString(field), value})
	}

	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}

	orderQuery := util.SnakeString(sortField)
	if sortOrder == "ascend" {
		session.Asc(orderQuery)
	} else {
		session.Desc(orderQuery)
	}

	err := session.Find(&applications)
	return applications, err
}

func (r *applicationRepository) Count(owner, field, value string) (int64, error) {
	return object.GetApplicationCount(owner, field, value)
}

func (r *applicationRepository) CountByOrganization(owner, organization, field, value string) (int64, error) {
	return object.GetOrganizationApplicationCount(owner, organization, field, value)
}

func (r *applicationRepository) Create(app *object.Application) (int64, error) {
	return r.engine.Insert(app)
}

func (r *applicationRepository) Update(id string, app *object.Application) (int64, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return 0, err
	}

	affected, err := r.engine.ID(core.PK{owner, name}).AllCols().Update(app)
	return affected, err
}

func (r *applicationRepository) Delete(app *object.Application) (int64, error) {
	return r.engine.ID(core.PK{app.Owner, app.Name}).Delete(&object.Application{})
}

func (r *applicationRepository) Exists(owner, name string) (bool, error) {
	return r.engine.Get(&object.Application{Owner: owner, Name: name})
}

func (r *applicationRepository) ExistsByClientId(clientId string) (bool, error) {
	return r.engine.Where("client_id = ?", clientId).Exist(&object.Application{})
}
