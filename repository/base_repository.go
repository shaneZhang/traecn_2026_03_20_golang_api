package repository

import (
	"fmt"

	"github.com/casdoor/casdoor/util"
	"github.com/xorm-io/builder"
	"github.com/xorm-io/xorm"
)

type BaseRepository interface {
	GetSession(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session
	GetSessionForUser(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session
}

type baseRepository struct {
	engine *xorm.Engine
}

func NewBaseRepository(engine *xorm.Engine) BaseRepository {
	return &baseRepository{engine: engine}
}

func (r *baseRepository) GetSession(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := r.engine.NewSession()
	defer session.Close()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	if field != "" && value != "" {
		if util.FilterField(field) {
			session.And(builder.Like{util.SnakeString(field), value})
		}
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

	return session
}

func (r *baseRepository) GetSessionForUser(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := r.engine.NewSession()
	defer session.Close()

	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}

	if owner != "" {
		session.Where("owner = ?", owner)
	}

	if field != "" && value != "" {
		if util.FilterField(field) {
			if offset != -1 {
				field = fmt.Sprintf("a.%s", field)
			}
			session.And(builder.Like{util.SnakeString(field), value})
		}
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

	return session
}
