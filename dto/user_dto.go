package dto

type UserCreateRequest struct {
	Owner             string            `json:"owner"`
	Name              string            `json:"name"`
	Id                string            `json:"id"`
	Type              string            `json:"type"`
	Password          string            `json:"password"`
	PasswordSalt      string            `json:"passwordSalt"`
	DisplayName       string            `json:"displayName"`
	FirstName         string            `json:"firstName"`
	LastName          string            `json:"lastName"`
	Email             string            `json:"email"`
	Phone             string            `json:"phone"`
	CountryCode       string            `json:"countryCode"`
	Region            string            `json:"region"`
	Location          string            `json:"location"`
	Address           []string          `json:"address"`
	Affiliation       string            `json:"affiliation"`
	Title             string            `json:"title"`
	Language          string            `json:"language"`
	Gender            string            `json:"gender"`
	Birthday          string            `json:"birthday"`
	Education         string            `json:"education"`
	SignupApplication string            `json:"signupApplication"`
	Groups            []string          `json:"groups"`
	Properties        map[string]string `json:"properties"`
}

type UserUpdateRequest struct {
	Id          string            `json:"id"`
	Owner       string            `json:"owner"`
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName"`
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	CountryCode string            `json:"countryCode"`
	Region      string            `json:"region"`
	Location    string            `json:"location"`
	Address     []string          `json:"address"`
	Affiliation string            `json:"affiliation"`
	Title       string            `json:"title"`
	Language    string            `json:"language"`
	Gender      string            `json:"gender"`
	Birthday    string            `json:"birthday"`
	Education   string            `json:"education"`
	Groups      []string          `json:"groups"`
	Properties  map[string]string `json:"properties"`
	IsAdmin     bool              `json:"isAdmin"`
	IsForbidden bool              `json:"isForbidden"`
	IsDeleted   bool              `json:"isDeleted"`
}

type UserQueryRequest struct {
	PaginationRequest
	Id       string `form:"id"`
	Email    string `form:"email"`
	Phone    string `form:"phone"`
	UserId   string `form:"userId"`
	IsOnline string `form:"isOnline"`
}

type UserBatchImportRequest struct {
	Owner string `json:"owner"`
}

type UserBatchExportRequest struct {
	Owner     string `json:"owner"`
	Field     string `json:"field"`
	Value     string `json:"value"`
	SortField string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
}

type SetPasswordRequest struct {
	UserOwner   string `json:"userOwner"`
	UserName    string `json:"userName"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
	Code        string `json:"code"`
}

type UserResponse struct {
	Owner            string        `json:"owner"`
	Name             string        `json:"name"`
	CreatedTime      string        `json:"createdTime"`
	UpdatedTime      string        `json:"updatedTime"`
	Id               string        `json:"id"`
	Type             string        `json:"type"`
	DisplayName      string        `json:"displayName"`
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	Avatar           string        `json:"avatar"`
	Email            string        `json:"email"`
	EmailVerified    bool          `json:"emailVerified"`
	Phone            string        `json:"phone"`
	CountryCode      string        `json:"countryCode"`
	Region           string        `json:"region"`
	Location         string        `json:"location"`
	Address          []string      `json:"address"`
	Affiliation      string        `json:"affiliation"`
	Title            string        `json:"title"`
	Language         string        `json:"language"`
	Gender           string        `json:"gender"`
	Birthday         string        `json:"birthday"`
	Education        string        `json:"education"`
	Score            int           `json:"score"`
	Ranking          int           `json:"ranking"`
	IsOnline         bool          `json:"isOnline"`
	IsAdmin          bool          `json:"isAdmin"`
	IsForbidden      bool          `json:"isForbidden"`
	IsDeleted        bool          `json:"isDeleted"`
	Groups           []string      `json:"groups"`
	Roles            []string      `json:"roles"`
	Permissions      []string      `json:"permissions"`
	MultiFactorAuths []interface{} `json:"multiFactorAuths"`
}

type UserListResponse struct {
	Users []*UserResponse `json:"users"`
	Total int64           `json:"total"`
}
