package apis

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/fkssl/cert"
	"github.com/suisrc/fkssl/kube"
	"github.com/suisrc/fkssl/serve"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func QurayCertCmdApi(ctx *gin.Context) {
	co, ok := bingQuery(ctx)
	if !ok {
		return
	}
	if co.Token == "" {
		serve.Error(ctx, 400, "TOKEN-ERROR", "token is empty")
		return
	}
	cli, err := kube.GetClient()
	if err != nil {
		serve.Error(ctx, 500, "KUBE-ERROR", err.Error())
		return
	}
	if co.Domains == nil || len(co.Domains) == 0 {
		serve.Error(ctx, 400, "DOMAIN-ERROR", "domains is empty")
		return
	}
	// ==========================================================================
	info, err := cli.CoreV1().Secrets(kube.CurrentNamespace()).Get(ctx, fmt.Sprintf("%s%s-%s", PK, co.Key, "info"), metav1.GetOptions{})
	if err != nil {
		serve.Error(ctx, 500, "KUBE-INFO-ERROR", err.Error())
		return
	}
	secretKey := fmt.Sprintf("%s%s-", PK, co.Key)
	if bts, ok := info.Data["prefix"]; ok {
		secretKey = string(bts)
	}
	if len(co.Domains) == 1 {
		secretKey += co.Domains[0]
	} else {
		sort.Strings(co.Domains)
		md5Domains, _ := hashMd5([]byte(strings.Join(co.Domains, ",")))
		secretKey += md5Domains
	}

	domain, err := cli.CoreV1().Secrets(kube.CurrentNamespace()).Get(ctx, secretKey, metav1.GetOptions{})
	if err == nil {
		if ok, _ := cert.IsExpired(string(domain.Data["pem.crt"])); !ok {
			serve.Success(ctx, gin.H{
				"crt": string(domain.Data["pem.crt"]),
				"key": string(domain.Data["pem.key"]),
			})
			return
		}
		// 证书出现问题或者过期
	}
	// ==========================================================================
	if co.Kind != 1 {
		serve.Error(ctx, 400, "KIND-ERROR", "kind is error")
		return
	}
	// domain对应的cert不存在，重写生成cert
	dns := []string{} // 域名
	ips := []string{}
	reg, _ := regexp.Compile(`^(\d{1,3}\.){3}\d{1,3}$`)
	for _, domain := range co.Domains {
		// 正则表达式匹配IP， 暂时支持ipv4
		if ok := reg.MatchString(domain); ok {
			ips = append(ips, domain) // ipv4
		} else {
			dns = append(dns, domain) // 域名
		}
	}
	configBts, ok0 := info.Data["config"]
	caCrt, ok1 := info.Data["ca.crt"]
	caKey, ok2 := info.Data["ca.key"]
	if !ok0 || !ok1 || !ok2 {
		serve.Error(ctx, 200, "CA-NOT-FOUND", "CA证书不存在")
		return
	}
	config := cert.CertConfig{}
	if err := json.Unmarshal(configBts, &config); err != nil {
		serve.Error(ctx, 500, "KUBE-CONFIG-ERROR", err.Error())
		return
	}
	subCert, err := cert.CreateCert(config, co.CommonName, co.Profile, dns, ips, string(caCrt), string(caKey))
	if err != nil {
		serve.Error(ctx, 500, "KUBE-CERT-ERROR", err.Error())
		return
	}
	// ==========================================================================
	// 写入k8s secret
	domain = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretKey,
		},
		StringData: map[string]string{
			"pem.crt": subCert.Crt,
			"pem.key": subCert.Key,
			"domains": strings.Join(co.Domains, ","),
		},
	}
	if _, err := cli.CoreV1().Secrets(kube.CurrentNamespace()).Create(ctx, domain, metav1.CreateOptions{}); err != nil {
		serve.Error(ctx, 500, "KUBE-CREATE-ERROR", err.Error())
		return
	}
	// ==========================================================================
	serve.Success(ctx, gin.H{
		"crt": subCert.Crt,
		"key": subCert.Key,
	})
}
