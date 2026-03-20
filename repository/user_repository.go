package repository

import (
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/core"
	"github.com/xorm-io/xorm"
)

type UserRepository interface {
	GetById(id string) (*object.User, error)
	GetByOwnerAndName(owner, name string) (*object.User, error)
	GetByEmail(owner, email string) (*object.User, error)
	GetByPhone(owner, phone string) (*object.User, error)
	GetByUserId(owner, userId string) (*object.User, error)
	GetByEmailOnly(email string) (*object.User, error)
	GetByPhoneOnly(phone string) (*object.User, error)
	GetByUserIdOnly(userId string) (*object.User, error)
	GetByFields(organization, username string) (*object.User, error)

	List(owner string, offset, limit int, field, value, sortField, sortOrder, groupName string) ([]*object.User, error)
	ListGlobal(offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error)
	ListByGroup(groupId string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error)
	ListSorted(owner, sorter string, limit int) ([]*object.User, error)

	Count(owner, field, value, groupName string) (int64, error)
	CountGlobal(field, value string) (int64, error)
	CountOnline(owner string, isOnline int) (int64, error)
	CountByGroup(groupId, field, value string) (int64, error)

	Create(user *object.User) (int64, error)
	CreateBatch(users []*object.User) (int64, error)
	Update(id string, user *object.User, columns []string, isAdmin bool) (int64, error)
	Delete(user *object.User) (int64, error)

	Exists(owner, name string) (bool, error)
	ExistsByEmail(owner, email string) (bool, error)
	ExistsByPhone(owner, phone string) (bool, error)
	ExistsByUserId(owner, userId string) (bool, error)

	GetWithFilter(owner string, cond builder.Cond) ([]*object.User, error)
}

type userRepository struct {
	engine *xorm.Engine
}

func NewUserRepository(engine *xorm.Engine) UserRepository {
	return &userRepository{engine: engine}
}

func (r *userRepository) GetById(id string) (*object.User, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return nil, err
	}
	return r.GetByOwnerAndName(owner, name)
}

func (r *userRepository) GetByOwnerAndName(owner, name string) (*object.User, error) {
	if owner == "" || name == "" {
		return nil, nil
	}
	user := &object.User{Owner: owner, Name: name}
	existed, err := r.engine.Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return user, nil
}

func (r *userRepository) GetByEmail(owner, email string) (*object.User, error) {
	if owner == "" || email == "" {
		return nil, nil
	}
	user := &object.User{Owner: owner, Email: email}
	existed, err := r.engine.Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return user, nil
}

func (r *userRepository) GetByPhone(owner, phone string) (*object.User, error) {
	if owner == "" || phone == "" {
		return nil, nil
	}
	user := &object.User{Owner: owner, Phone: phone}
	existed, err := r.engine.Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return user, nil
}

func (r *userRepository) GetByUserId(owner, userId string) (*object.User, error) {
	if owner == "" || userId == "" {
		return nil, nil
	}
	user := &object.User{Owner: owner, Id: userId}
	existed, err := r.engine.Get(user)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, nil
	}
	return user, nil
}

func (r *userRepository) GetByEmailOnly(email string) (*object.User, error) {
	return object.GetUserByEmailOnly(email)
}

func (r *userRepository) GetByPhoneOnly(phone string) (*object.User, error) {
	return object.GetUserByPhoneOnly(phone)
}

func (r *userRepository) GetByUserIdOnly(userId string) (*object.User, error) {
	return object.GetUserByUserIdOnly(userId)
}

func (r *userRepository) GetByFields(organization, username string) (*object.User, error) {
	return object.GetUserByFields(organization, username)
}

func (r *userRepository) List(owner string, offset, limit int, field, value, sortField, sortOrder, groupName string) ([]*object.User, error) {
	if groupName != "" {
		return r.ListByGroup(util.GetId(owner, groupName), offset, limit, field, value, sortField, sortOrder)
	}

	users := []*object.User{}
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

	err := session.Find(&users)
	return users, err
}

func (r *userRepository) ListGlobal(offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error) {
	users := []*object.User{}
	session := r.engine.NewSession()
	defer session.Close()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if field != "" && value != "" {
		session.Where(builder.Like{util.SnakeString(field), value})
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

	err := session.Find(&users)
	return users, err
}

func (r *userRepository) ListByGroup(groupId string, offset, limit int, field, value, sortField, sortOrder string) ([]*object.User, error) {
	return object.GetPaginationGroupUsers(groupId, offset, limit, field, value, sortField, sortOrder)
}

func (r *userRepository) ListSorted(owner, sorter string, limit int) ([]*object.User, error) {
	return object.GetSortedUsers(owner, sorter, limit)
}

func (r *userRepository) Count(owner, field, value, groupName string) (int64, error) {
	if groupName != "" {
		return r.CountByGroup(util.GetId(owner, groupName), field, value)
	}

	session := r.engine.NewSession()
	defer session.Close()

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	if field != "" && value != "" {
		session.And(builder.Like{util.SnakeString(field), value})
	}

	return session.Count(&object.User{})
}

func (r *userRepository) CountGlobal(field, value string) (int64, error) {
	session := r.engine.NewSession()
	defer session.Close()

	if field != "" && value != "" {
		session.Where(builder.Like{util.SnakeString(field), value})
	}

	return session.Count(&object.User{})
}

func (r *userRepository) CountOnline(owner string, isOnline int) (int64, error) {
	return r.engine.Where("is_online = ?", isOnline).Count(&object.User{Owner: owner})
}

func (r *userRepository) CountByGroup(groupId, field, value string) (int64, error) {
	return object.GetGroupUserCount(groupId, field, value)
}

func (r *userRepository) Create(user *object.User) (int64, error) {
	return r.engine.Insert(user)
}

func (r *userRepository) CreateBatch(users []*object.User) (int64, error) {
	if len(users) == 0 {
		return 0, nil
	}
	return r.engine.Insert(users)
}

func (r *userRepository) Update(id string, user *object.User, columns []string, isAdmin bool) (int64, error) {
	owner, name, err := util.GetOwnerAndNameFromIdWithError(id)
	if err != nil {
		return 0, err
	}

	session := r.engine.ID(core.PK{owner, name})

	if len(columns) > 0 {
		session.Cols(columns...)
	} else {
		session.AllCols()
	}

	affected, err := session.Update(user)
	return affected, err
}

func (r *userRepository) Delete(user *object.User) (int64, error) {
	return r.engine.ID(core.PK{user.Owner, user.Name}).Delete(&object.User{})
}

func (r *userRepository) Exists(owner, name string) (bool, error) {
	return r.engine.Get(&object.User{Owner: owner, Name: name})
}

func (r *userRepository) ExistsByEmail(owner, email string) (bool, error) {
	return r.engine.Get(&object.User{Owner: owner, Email: email})
}

func (r *userRepository) ExistsByPhone(owner, phone string) (bool, error) {
	return r.engine.Get(&object.User{Owner: owner, Phone: phone})
}

func (r *userRepository) ExistsByUserId(owner, userId string) (bool, error) {
	return r.engine.Get(&object.User{Owner: owner, Id: userId})
}

func (r *userRepository) GetWithFilter(owner string, cond builder.Cond) ([]*object.User, error) {
	users := []*object.User{}
	session := r.engine.Desc("created_time")
	if cond != nil {
		session = session.Where(cond)
	}
	err := session.Find(&users, &object.User{Owner: owner})
	return users, err
}
