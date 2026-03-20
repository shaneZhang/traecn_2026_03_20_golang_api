package repository

import (
	"github.com/casdoor/casdoor/object"
	"github.com/xorm-io/xorm"
)

type MfaRepository interface {
	GetMfaUtil(mfaType string, config *object.MfaProps) object.MfaInterface
	GetAllMfaProps(user *object.User, masked bool) []*object.MfaProps
	GetMfaPropsByType(user *object.User, mfaType string, masked bool) *object.MfaProps
}

type mfaRepository struct {
	engine *xorm.Engine
}

func NewMfaRepository(engine *xorm.Engine) MfaRepository {
	return &mfaRepository{engine: engine}
}

func (r *mfaRepository) GetMfaUtil(mfaType string, config *object.MfaProps) object.MfaInterface {
	return object.GetMfaUtil(mfaType, config)
}

func (r *mfaRepository) GetAllMfaProps(user *object.User, masked bool) []*object.MfaProps {
	return object.GetAllMfaProps(user, masked)
}

func (r *mfaRepository) GetMfaPropsByType(user *object.User, mfaType string, masked bool) *object.MfaProps {
	return user.GetMfaProps(mfaType, masked)
}
