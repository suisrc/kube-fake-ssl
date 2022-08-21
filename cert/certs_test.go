package cert_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suisrc/fkssl/cert"
)

func TestCert(t *testing.T) {
	// 读取bin/config.json文件内容给cert.CertConfig对象
	bts, _ := ioutil.ReadFile("../bin/config.json")
	cfg := cert.CertConfig{}
	json.Unmarshal(bts, &cfg)

	// 生成证书
	ca, err := cert.CreateCA(cfg)
	assert.Nil(t, err)

	ct, err := cert.CreateCert(cfg, "dev1", "", []string{"sso.dev1.com"}, nil, ca.Crt, ca.Key)
	assert.Nil(t, err)

	// 保存证书
	ioutil.WriteFile("../bin/ca.crt", []byte(ca.Crt), 0644)
	ioutil.WriteFile("../bin/ca.key", []byte(ca.Key), 0644)
	ioutil.WriteFile("../bin/dev1.crt", []byte(ct.Crt), 0644)
	ioutil.WriteFile("../bin/dev1.key", []byte(ct.Key), 0644)

}
