package dto

type OrganizationCreateRequest struct {
	Owner               string   `json:"owner"`
	Name                string   `json:"name"`
	DisplayName         string   `json:"displayName"`
	WebsiteUrl          string   `json:"websiteUrl"`
	Logo                string   `json:"logo"`
	PasswordType        string   `json:"passwordType"`
	PasswordSalt        string   `json:"passwordSalt"`
	PasswordOptions     []string `json:"passwordOptions"`
	DefaultAvatar       string   `json:"defaultAvatar"`
	DefaultApplication  string   `json:"defaultApplication"`
	Tags                []string `json:"tags"`
	Languages           []string `json:"languages"`
	InitScore           int      `json:"initScore"`
	EnableSoftDeletion  bool     `json:"enableSoftDeletion"`
	IsProfilePublic     bool     `json:"isProfilePublic"`
	MfaItems            []string `json:"mfaItems"`
	AccountItems        []string `json:"accountItems"`
}

type OrganizationUpdateRequest struct {
	Owner               string   `json:"owner"`
	Name                string   `json:"name"`
	DisplayName         string   `json:"displayName"`
	WebsiteUrl          string   `json:"websiteUrl"`
	Logo                string   `json:"logo"`
	PasswordType        string   `json:"passwordType"`
	PasswordSalt        string   `json:"passwordSalt"`
	PasswordOptions     []string `json:"passwordOptions"`
	DefaultAvatar       string   `json:"defaultAvatar"`
	DefaultApplication  string   `json:"defaultApplication"`
	Tags                []string `json:"tags"`
	Languages           []string `json:"languages"`
	InitScore           int      `json:"initScore"`
	EnableSoftDeletion  bool     `json:"enableSoftDeletion"`
	IsProfilePublic     bool     `json:"isProfilePublic"`
	MfaItems            []string `json:"mfaItems"`
	AccountItems        []string `json:"accountItems"`
	IpWhitelist         string   `json:"ipWhitelist"`
}

type OrganizationQueryRequest struct {
	PaginationRequest
	Id               string `form:"id"`
	OrganizationName string `form:"organizationName"`
}

type OrganizationResponse struct {
	Owner               string   `json:"owner"`
	Name                string   `json:"name"`
	CreatedTime         string   `json:"createdTime"`
	DisplayName         string   `json:"displayName"`
	WebsiteUrl          string   `json:"websiteUrl"`
	Logo                string   `json:"logo"`
	PasswordType        string   `json:"passwordType"`
	DefaultAvatar       string   `json:"defaultAvatar"`
	DefaultApplication  string   `json:"defaultApplication"`
	Tags                []string `json:"tags"`
	Languages           []string `json:"languages"`
	InitScore           int      `json:"initScore"`
	EnableSoftDeletion  bool     `json:"enableSoftDeletion"`
	IsProfilePublic     bool     `json:"isProfilePublic"`
	MfaRememberInHours  int      `json:"mfaRememberInHours"`
}

type GroupCreateRequest struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Manager      string `json:"manager"`
	ContactEmail string `json:"contactEmail"`
	Type         string `json:"type"`
	ParentId     string `json:"parentId"`
	IsTopGroup   bool   `json:"isTopGroup"`
	IsEnabled    bool   `json:"isEnabled"`
}

type GroupUpdateRequest struct {
	Owner        string `json:"owner"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Manager      string `json:"manager"`
	ContactEmail string `json:"contactEmail"`
	Type         string `json:"type"`
	ParentId     string `json:"parentId"`
	IsTopGroup   bool   `json:"isTopGroup"`
	IsEnabled    bool   `json:"isEnabled"`
}

type GroupQueryRequest struct {
	PaginationRequest
	Id       string `form:"id"`
	WithTree string `form:"withTree"`
}

type GroupResponse struct {
	Owner        string   `json:"owner"`
	Name         string   `json:"name"`
	CreatedTime  string   `json:"createdTime"`
	DisplayName  string   `json:"displayName"`
	Manager      string   `json:"manager"`
	ContactEmail string   `json:"contactEmail"`
	Type         string   `json:"type"`
	ParentId     string   `json:"parentId"`
	ParentName   string   `json:"parentName"`
	IsTopGroup   bool     `json:"isTopGroup"`
	Users        []string `json:"users"`
	HaveChildren bool     `json:"haveChildren"`
	Children     []interface{} `json:"children"`
	IsEnabled    bool     `json:"isEnabled"`
}
