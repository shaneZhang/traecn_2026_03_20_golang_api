package repository

import (
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/core"
	"github.com/xorm-io/xorm"
)

type OrganizationRepository interface {
	GetById(id string) (*object.Organization, error)
	GetByOwnerAndName(owner, name string) (*object.Organization, error)
	GetByUser(user *object.User) (*object.Organization, error)

	List(owner string, names ...string) ([]*object.Organization, error)
	ListByFields(owner string, fields ...string) ([]*object.Organization, error)
	ListPagination(owner, name string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Organization, error)

	Count(owner, name, field, value string) (int64, error)

	Create(org *object.Organization) (int64, error)
	Update(id string, org *object.Organization, isGlobalAdmin bool) (int64, error)
	Delete(org *object.Organization) (int64, error)

	Exists(owner, name string) (bool, error)
}

type organizationRepository struct {
	engine *xorm.Engine
}

func NewOrganizationRepository(engine *xorm.Engine) OrganizationRepository {
	return &organizationRepository{engine: engine}
}

func (r *organizationRepository) GetById(id string) (*object.Organization, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return nil, err
	}
	return r.GetByOwnerAndName(owner, name)
}

func (r *organizationRepository) GetByOwnerAndName(owner, name string) (*object.Organization, error) {
	if owner == "" || name == "" {
		return nil, nil
	}
	org := &object.Organization{Owner: owner, Name: name}
	existed, err := r.engine.Get(org)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return org, nil
}

func (r *organizationRepository) GetByUser(user *object.User) (*object.Organization, error) {
	if user == nil {
		return nil, nil
	}
	return r.GetByOwnerAndName("admin", user.Owner)
}

func (r *organizationRepository) List(owner string, names ...string) ([]*object.Organization, error) {
	organizations := []*object.Organization{}
	var err error

	if len(names) > 0 {
		err = r.engine.Desc("created_time").Where(builder.In("name", names)).Find(&organizations)
	} else {
		err = r.engine.Desc("created_time").Find(&organizations, &object.Organization{Owner: owner})
	}

	return organizations, err
}

func (r *organizationRepository) ListByFields(owner string, fields ...string) ([]*object.Organization, error) {
	organizations := []*object.Organization{}
	err := r.engine.Desc("created_time").Cols(fields...).Find(&organizations, &object.Organization{Owner: owner})
	return organizations, err
}

func (r *organizationRepository) ListPagination(owner, name string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Organization, error) {
	organizations := []*object.Organization{}
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

	var err error
	if name != "" {
		err = session.Find(&organizations, &object.Organization{Name: name})
	} else {
		err = session.Find(&organizations)
	}

	return organizations, err
}

func (r *organizationRepository) Count(owner, name, field, value string) (int64, error) {
	session := r.engine.NewSession()
	defer session.Close()

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	if field != "" && value != "" {
		session.And(builder.Like{util.SnakeString(field), value})
	}

	return session.Count(&object.Organization{Name: name})
}

func (r *organizationRepository) Create(org *object.Organization) (int64, error) {
	return r.engine.Insert(org)
}

func (r *organizationRepository) Update(id string, org *object.Organization, isGlobalAdmin bool) (int64, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return 0, err
	}

	oldOrg, err := r.GetByOwnerAndName(owner, name)
	if err != nil {
		return 0, err
	}
	if oldOrg == nil {
		return 0, nil
	}

	if name == "built-in" {
		org.Name = name
	}

	affected, err := r.engine.ID(core.PK{owner, name}).AllCols().Update(org)
	return affected, err
}

func (r *organizationRepository) Delete(org *object.Organization) (int64, error) {
	if org.Name == "built-in" {
		return 0, nil
	}
	return r.engine.ID(core.PK{org.Owner, org.Name}).Delete(&object.Organization{})
}

func (r *organizationRepository) Exists(owner, name string) (bool, error) {
	return r.engine.Get(&object.Organization{Owner: owner, Name: name})
}

type GroupRepository interface {
	GetById(id string) (*object.Group, error)
	GetByOwnerAndName(owner, name string) (*object.Group, error)

	List(owner string) ([]*object.Group, error)
	ListGlobal() ([]*object.Group, error)
	ListPagination(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Group, error)

	Count(owner, field, value string) (int64, error)

	Create(group *object.Group) (int64, error)
	CreateBatch(groups []*object.Group) (int64, error)
	Update(id string, group *object.Group) (int64, error)
	Delete(group *object.Group) (int64, error)

	GetHaveChildrenMap(groups []*object.Group) (map[string]*object.Group, error)
}

type groupRepository struct {
	engine *xorm.Engine
}

func NewGroupRepository(engine *xorm.Engine) GroupRepository {
	return &groupRepository{engine: engine}
}

func (r *groupRepository) GetById(id string) (*object.Group, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return nil, err
	}
	return r.GetByOwnerAndName(owner, name)
}

func (r *groupRepository) GetByOwnerAndName(owner, name string) (*object.Group, error) {
	if owner == "" || name == "" {
		return nil, nil
	}
	group := &object.Group{Owner: owner, Name: name}
	existed, err := r.engine.Get(group)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return group, nil
}

func (r *groupRepository) List(owner string) ([]*object.Group, error) {
	groups := []*object.Group{}
	err := r.engine.Desc("created_time").Find(&groups, &object.Group{Owner: owner})
	return groups, err
}

func (r *groupRepository) ListGlobal() ([]*object.Group, error) {
	groups := []*object.Group{}
	err := r.engine.Desc("created_time").Find(&groups)
	return groups, err
}

func (r *groupRepository) ListPagination(owner string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.Group, error) {
	groups := []*object.Group{}
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

	err := session.Find(&groups)
	return groups, err
}

func (r *groupRepository) Count(owner, field, value string) (int64, error) {
	session := r.engine.NewSession()
	defer session.Close()

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	if field != "" && value != "" {
		session.And(builder.Like{util.SnakeString(field), value})
	}

	return session.Count(&object.Group{})
}

func (r *groupRepository) Create(group *object.Group) (int64, error) {
	return r.engine.Insert(group)
}

func (r *groupRepository) CreateBatch(groups []*object.Group) (int64, error) {
	if len(groups) == 0 {
		return 0, nil
	}
	return r.engine.Insert(groups)
}

func (r *groupRepository) Update(id string, group *object.Group) (int64, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return 0, err
	}

	affected, err := r.engine.ID(core.PK{owner, name}).AllCols().Update(group)
	return affected, err
}

func (r *groupRepository) Delete(group *object.Group) (int64, error) {
	return r.engine.ID(core.PK{group.Owner, group.Name}).Delete(&object.Group{})
}

func (r *groupRepository) GetHaveChildrenMap(groups []*object.Group) (map[string]*object.Group, error) {
	return object.GetGroupsHaveChildrenMap(groups)
}
