package cert

type CertConfig struct {
	CommonName   string                 `json:"CN"`
	SignKey      SignKey                `json:"key"`
	CaProfile    SignProfile            `json:"CA"`
	SignProfiles map[string]SignProfile `json:"profiles"`
}

type SignKey struct {
	Size int `json:"size"`
}

type SubjectName struct {
	Country          string `json:"C"`
	Province         string `json:"ST"`
	Locality         string `json:"L"`
	Organization     string `json:"O"`
	OrganizationUnit string `json:"OU"`
}

type SignProfile struct {
	Expiry      string       `json:"expiry"`
	CommonName  string       `json:"CN"`
	SignKey     SignKey      `json:"key"`
	SubjectName *SubjectName `json:"name"`
}

//===========================================================================

type SignResult struct {
	Crt string `json:"crt"`
	Key string `json:"key"`
}

//===========================================================================

// 合并配置
func (aa *CertConfig) Merge(bb CertConfig) bool {
	return false
}

// String
func (aa *CertConfig) String() string {
	return ""
}
