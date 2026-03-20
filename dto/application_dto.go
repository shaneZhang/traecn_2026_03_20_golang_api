package dto

type ApplicationCreateRequest struct {
	Owner                   string   `json:"owner"`
	Name                    string   `json:"name"`
	DisplayName             string   `json:"displayName"`
	Category                string   `json:"category"`
	Type                    string   `json:"type"`
	Logo                    string   `json:"logo"`
	HomepageUrl             string   `json:"homepageUrl"`
	Description             string   `json:"description"`
	Organization            string   `json:"organization"`
	Cert                    string   `json:"cert"`
	EnablePassword          bool     `json:"enablePassword"`
	EnableSignUp            bool     `json:"enableSignUp"`
	EnableCodeSignin        bool     `json:"enableCodeSignin"`
	EnableWebAuthn          bool     `json:"enableWebAuthn"`
	RedirectUris            []string `json:"redirectUris"`
	TokenFormat             string   `json:"tokenFormat"`
	ExpireInHours           float64  `json:"expireInHours"`
	RefreshExpireInHours    float64  `json:"refreshExpireInHours"`
	SignupUrl               string   `json:"signupUrl"`
	SigninUrl               string   `json:"signinUrl"`
	ForgetUrl               string   `json:"forgetUrl"`
}

type ApplicationUpdateRequest struct {
	Owner                   string   `json:"owner"`
	Name                    string   `json:"name"`
	DisplayName             string   `json:"displayName"`
	Category                string   `json:"category"`
	Type                    string   `json:"type"`
	Logo                    string   `json:"logo"`
	HomepageUrl             string   `json:"homepageUrl"`
	Description             string   `json:"description"`
	Organization            string   `json:"organization"`
	Cert                    string   `json:"cert"`
	EnablePassword          bool     `json:"enablePassword"`
	EnableSignUp            bool     `json:"enableSignUp"`
	EnableCodeSignin        bool     `json:"enableCodeSignin"`
	EnableWebAuthn          bool     `json:"enableWebAuthn"`
	RedirectUris            []string `json:"redirectUris"`
	TokenFormat             string   `json:"tokenFormat"`
	ExpireInHours           float64  `json:"expireInHours"`
	RefreshExpireInHours    float64  `json:"refreshExpireInHours"`
	SignupUrl               string   `json:"signupUrl"`
	SigninUrl               string   `json:"signinUrl"`
	ForgetUrl               string   `json:"forgetUrl"`
	IpWhitelist             string   `json:"ipWhitelist"`
}

type ApplicationQueryRequest struct {
	PaginationRequest
	Id            string `form:"id"`
	Organization  string `form:"organization"`
	WithKey       string `form:"withKey"`
}

type ApplicationResponse struct {
	Owner                   string   `json:"owner"`
	Name                    string   `json:"name"`
	CreatedTime             string   `json:"createdTime"`
	DisplayName             string   `json:"displayName"`
	Category                string   `json:"category"`
	Type                    string   `json:"type"`
	Logo                    string   `json:"logo"`
	HomepageUrl             string   `json:"homepageUrl"`
	Description             string   `json:"description"`
	Organization            string   `json:"organization"`
	Cert                    string   `json:"cert"`
	ClientId                string   `json:"clientId"`
	EnablePassword          bool     `json:"enablePassword"`
	EnableSignUp            bool     `json:"enableSignUp"`
	EnableCodeSignin        bool     `json:"enableCodeSignin"`
	EnableWebAuthn          bool     `json:"enableWebAuthn"`
	RedirectUris            []string `json:"redirectUris"`
	TokenFormat             string   `json:"tokenFormat"`
	ExpireInHours           float64  `json:"expireInHours"`
	RefreshExpireInHours    float64  `json:"refreshExpireInHours"`
	SignupUrl               string   `json:"signupUrl"`
	SigninUrl               string   `json:"signinUrl"`
	ForgetUrl               string   `json:"forgetUrl"`
}

type OAuthGrantRequest struct {
	ClientId     string   `json:"clientId"`
	RedirectUri  string   `json:"redirectUri"`
	Scope        string   `json:"scope"`
	State        string   `json:"state"`
	ResponseType string   `json:"responseType"`
	CodeChallenge string  `json:"codeChallenge"`
	CodeChallengeMethod string `json:"codeChallengeMethod"`
}

type OAuthTokenRequest struct {
	GrantType    string `json:"grantType"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirectUri"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	CodeVerifier string `json:"codeVerifier"`
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Scope        string `json:"scope"`
}

type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	IdToken      string `json:"id_token"`
}
