package dto

type MfaSetupInitiateRequest struct {
	Owner   string `json:"owner" form:"owner"`
	Name    string `json:"name" form:"name"`
	MfaType string `json:"mfaType" form:"mfaType"`
}

type MfaSetupVerifyRequest struct {
	MfaType     string `json:"mfaType" form:"mfaType"`
	Passcode    string `json:"passcode" form:"passcode"`
	Secret      string `json:"secret" form:"secret"`
	Dest        string `json:"dest" form:"dest"`
	CountryCode string `json:"countryCode" form:"countryCode"`
}

type MfaSetupEnableRequest struct {
	Owner        string `json:"owner" form:"owner"`
	Name         string `json:"name" form:"name"`
	MfaType      string `json:"mfaType" form:"mfaType"`
	Secret       string `json:"secret" form:"secret"`
	Dest         string `json:"dest" form:"dest"`
	CountryCode  string `json:"countryCode" form:"countryCode"`
	RecoveryCodes string `json:"recoveryCodes" form:"recoveryCodes"`
}

type MfaDeleteRequest struct {
	Owner string `json:"owner" form:"owner"`
	Name  string `json:"name" form:"name"`
}

type MfaSetPreferredRequest struct {
	Owner   string `json:"owner" form:"owner"`
	Name    string `json:"name" form:"name"`
	MfaType string `json:"mfaType" form:"mfaType"`
}

type MfaVerifyRequest struct {
	MfaType  string `json:"mfaType" form:"mfaType"`
	Passcode string `json:"passcode" form:"passcode"`
}

type MfaPropsResponse struct {
	Enabled            bool     `json:"enabled"`
	IsPreferred        bool     `json:"isPreferred"`
	MfaType            string   `json:"mfaType"`
	Secret             string   `json:"secret,omitempty"`
	CountryCode        string   `json:"countryCode,omitempty"`
	URL                string   `json:"url,omitempty"`
	RecoveryCodes      []string `json:"recoveryCodes,omitempty"`
	MfaRememberInHours int      `json:"mfaRememberInHours"`
}

type MfaSetupInitiateResponse struct {
	MfaPropsResponse
}
