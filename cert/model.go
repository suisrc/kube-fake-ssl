package cert

import "encoding/json"

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
	update := false
	if bb.SignKey.Size > 0 && bb.SignKey.Size != aa.SignKey.Size {
		aa.SignKey.Size = bb.SignKey.Size
		update = true
	}
	for bKey, bVal := range bb.SignProfiles {
		aa.SignProfiles[bKey] = bVal
		update = true
	}
	return update
}

// String
func (aa *CertConfig) String() string {
	str, _ := json.Marshal(aa)
	return string(str)
}
