package template

type AuthConfig struct {
	Enabled       bool
	LoginURL      string
	LogoutURL     string
	RegisterURL   string
	UserKey       string
	PermissionsKey string
}

type AuthHelper struct {
	config AuthConfig
}

func NewAuthHelper(config AuthConfig) *AuthHelper {
	return &AuthHelper{config: config}
}

func (a *AuthHelper) PrepareData(data interface{}) interface{} {
	if data == nil {
		data = make(map[string]interface{})
	}
	
	if m, ok := data.(map[string]interface{}); ok {
		m["auth"] = map[string]interface{}{
			"is_authenticated": a.IsAuthenticated(m),
			"user":             a.GetUser(m),
			"login_url":        a.config.LoginURL,
			"logout_url":       a.config.LogoutURL,
			"register_url":     a.config.RegisterURL,
		}
		return m
	}
	
	return data
}

func (a *AuthHelper) IsAuthenticated(data interface{}) bool {
	if m, ok := data.(map[string]interface{}); ok {
		user, exists := m[a.config.UserKey]
		return exists && user != nil
	}
	return false
}

func (a *AuthHelper) GetUser(data interface{}) interface{} {
	if m, ok := data.(map[string]interface{}); ok {
		return m[a.config.UserKey]
	}
	return nil
}

func (a *AuthHelper) HasPermission(data interface{}, permission string) bool {
	if m, ok := data.(map[string]interface{}); ok {
		if perms, exists := m[a.config.PermissionsKey]; exists {
			if list, ok := perms.([]string); ok {
				for _, p := range list {
					if p == permission {
						return true
					}
				}
			}
		}
	}
	return false
}