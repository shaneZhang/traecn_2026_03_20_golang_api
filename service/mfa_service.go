package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/casdoor/casdoor/dto"
	"github.com/casdoor/casdoor/object"
	"github.com/casdoor/casdoor/repository"
	"github.com/casdoor/casdoor/util"
	"github.com/google/uuid"
)

type MfaService interface {
	SetupInitiate(owner, name, mfaType string) (*dto.MfaSetupInitiateResponse, error)
	SetupVerify(mfaType, passcode, secret, dest, countryCode string) error
	SetupEnable(owner, name, mfaType, secret, dest, countryCode, recoveryCodes string) error
	DeleteMfa(owner, name string) ([]*object.MfaProps, error)
	SetPreferredMfa(owner, name, mfaType string) ([]*object.MfaProps, error)
	
	GetAllMfaProps(user *object.User, masked bool) []*object.MfaProps
	GetMfaPropsByType(user *object.User, mfaType string, masked bool) *object.MfaProps
	VerifyMfa(mfaType, passcode string, user *object.User) error
	RecoverMfa(user *object.User, recoveryCode string) error
}

type mfaService struct {
	mfaRepo  repository.MfaRepository
	userRepo repository.UserRepository
	orgRepo  repository.OrganizationRepository
}

func NewMfaService(mfaRepo repository.MfaRepository, userRepo repository.UserRepository, orgRepo repository.OrganizationRepository) MfaService {
	return &mfaService{
		mfaRepo:  mfaRepo,
		userRepo: userRepo,
		orgRepo:  orgRepo,
	}
}

func (s *mfaService) SetupInitiate(owner, name, mfaType string) (*dto.MfaSetupInitiateResponse, error) {
	userId := util.GetId(owner, name)
	if len(userId) == 0 {
		return nil, errors.New(http.StatusText(http.StatusBadRequest))
	}
	
	mfaUtil := s.mfaRepo.GetMfaUtil(mfaType, nil)
	if mfaUtil == nil {
		return nil, errors.New("invalid auth type")
	}
	
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user doesn't exist")
	}
	
	org, err := s.orgRepo.GetByUser(user)
	if err != nil {
		return nil, err
	}
	
	issuer := ""
	if org != nil && org.DisplayName != "" {
		issuer = org.DisplayName
	} else if org != nil {
		issuer = org.Name
	}
	
	mfaProps, err := mfaUtil.Initiate(user.GetId(), issuer)
	if err != nil {
		return nil, err
	}
	
	recoveryCode := uuid.NewString()
	mfaProps.RecoveryCodes = []string{recoveryCode}
	if org != nil {
		mfaProps.MfaRememberInHours = org.MfaRememberInHours
	}
	
	return &dto.MfaSetupInitiateResponse{
		MfaPropsResponse: dto.MfaPropsResponse{
			Enabled:            mfaProps.Enabled,
			IsPreferred:        mfaProps.IsPreferred,
			MfaType:            mfaProps.MfaType,
			Secret:             mfaProps.Secret,
			CountryCode:        mfaProps.CountryCode,
			URL:                mfaProps.URL,
			RecoveryCodes:      mfaProps.RecoveryCodes,
			MfaRememberInHours: mfaProps.MfaRememberInHours,
		},
	}, nil
}

func (s *mfaService) SetupVerify(mfaType, passcode, secret, dest, countryCode string) error {
	if mfaType == "" || passcode == "" {
		return errors.New("missing auth type or passcode")
	}
	
	config := &object.MfaProps{
		MfaType: mfaType,
	}
	
	switch mfaType {
	case object.TotpType:
		if secret == "" {
			return errors.New("totp secret is missing")
		}
		config.Secret = secret
	case object.SmsType:
		if dest == "" {
			return errors.New("destination is missing")
		}
		config.Secret = dest
		if countryCode == "" {
			return errors.New("country code is missing")
		}
		config.CountryCode = countryCode
	case object.EmailType:
		if dest == "" {
			return errors.New("destination is missing")
		}
		config.Secret = dest
	case object.RadiusType:
		if dest == "" {
			return errors.New("RADIUS username is missing")
		}
		config.Secret = dest
		if secret == "" {
			return errors.New("RADIUS provider is missing")
		}
		config.URL = secret
	case object.PushType:
		if dest == "" {
			return errors.New("push notification receiver is missing")
		}
		config.Secret = dest
		if secret == "" {
			return errors.New("push notification provider is missing")
		}
		config.URL = secret
	}
	
	mfaUtil := s.mfaRepo.GetMfaUtil(mfaType, config)
	if mfaUtil == nil {
		return errors.New("invalid multi-factor authentication type")
	}
	
	return mfaUtil.SetupVerify(passcode)
}

func (s *mfaService) SetupEnable(owner, name, mfaType, secret, dest, countryCode, recoveryCodes string) error {
	user, err := s.userRepo.GetById(util.GetId(owner, name))
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user doesn't exist")
	}
	
	config := &object.MfaProps{
		MfaType: mfaType,
	}
	
	switch mfaType {
	case object.TotpType:
		if secret == "" {
			return errors.New("totp secret is missing")
		}
		config.Secret = secret
	case object.EmailType:
		if user.Email == "" {
			if dest == "" {
				return errors.New("destination is missing")
			}
			user.Email = dest
		}
	case object.SmsType:
		if user.Phone == "" {
			if dest == "" {
				return errors.New("destination is missing")
			}
			user.Phone = dest
			if countryCode == "" {
				return errors.New("country code is missing")
			}
			user.CountryCode = countryCode
		}
	case object.RadiusType:
		if dest == "" {
			return errors.New("RADIUS username is missing")
		}
		config.Secret = dest
		if secret == "" {
			return errors.New("RADIUS provider is missing")
		}
		config.URL = secret
	case object.PushType:
		if dest == "" {
			return errors.New("push notification receiver is missing")
		}
		config.Secret = dest
		if secret == "" {
			return errors.New("push notification provider is missing")
		}
		config.URL = secret
	}
	
	if recoveryCodes == "" {
		return errors.New("recovery codes is missing")
	}
	config.RecoveryCodes = []string{recoveryCodes}
	
	mfaUtil := s.mfaRepo.GetMfaUtil(mfaType, config)
	if mfaUtil == nil {
		return errors.New("invalid multi-factor authentication type")
	}
	
	return mfaUtil.Enable(user)
}

func (s *mfaService) DeleteMfa(owner, name string) ([]*object.MfaProps, error) {
	userId := util.GetId(owner, name)
	
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user doesn't exist")
	}
	
	err = object.DisabledMultiFactorAuth(user)
	if err != nil {
		return nil, err
	}
	
	return s.mfaRepo.GetAllMfaProps(user, true), nil
}

func (s *mfaService) SetPreferredMfa(owner, name, mfaType string) ([]*object.MfaProps, error) {
	userId := util.GetId(owner, name)
	
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user doesn't exist")
	}
	
	err = object.SetPreferredMultiFactorAuth(user, mfaType)
	if err != nil {
		return nil, err
	}
	
	return s.mfaRepo.GetAllMfaProps(user, true), nil
}

func (s *mfaService) GetAllMfaProps(user *object.User, masked bool) []*object.MfaProps {
	return s.mfaRepo.GetAllMfaProps(user, masked)
}

func (s *mfaService) GetMfaPropsByType(user *object.User, mfaType string, masked bool) *object.MfaProps {
	return s.mfaRepo.GetMfaPropsByType(user, mfaType, masked)
}

func (s *mfaService) VerifyMfa(mfaType, passcode string, user *object.User) error {
	mfaProps := s.mfaRepo.GetMfaPropsByType(user, mfaType, false)
	if mfaProps == nil {
		return errors.New("MFA not configured for this type")
	}
	
	mfaUtil := s.mfaRepo.GetMfaUtil(mfaType, &object.MfaProps{
		MfaType:     mfaProps.MfaType,
		Secret:      mfaProps.Secret,
		CountryCode: mfaProps.CountryCode,
		URL:         mfaProps.URL,
	})
	if mfaUtil == nil {
		return errors.New("invalid MFA type")
	}
	
	return mfaUtil.Verify(passcode)
}

func (s *mfaService) RecoverMfa(user *object.User, recoveryCode string) error {
	return object.MfaRecover(user, recoveryCode)
}

func (s *mfaService) validateMfaConfig(mfaType, secret, dest, countryCode string) error {
	switch mfaType {
	case object.TotpType:
		if secret == "" {
			return errors.New("totp secret is missing")
		}
	case object.SmsType:
		if dest == "" {
			return errors.New("destination is missing")
		}
		if countryCode == "" {
			return errors.New("country code is missing")
		}
	case object.EmailType:
		if dest == "" {
			return errors.New("destination is missing")
		}
	case object.RadiusType:
		if dest == "" {
			return errors.New("RADIUS username is missing")
		}
		if secret == "" {
			return errors.New("RADIUS provider is missing")
		}
	case object.PushType:
		if dest == "" {
			return errors.New("push notification receiver is missing")
		}
		if secret == "" {
			return errors.New("push notification provider is missing")
		}
	default:
		return fmt.Errorf("invalid MFA type: %s", mfaType)
	}
	return nil
}
