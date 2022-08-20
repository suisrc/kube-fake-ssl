package apis

import (
	"crypto/md5"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/fkssl/cert"
	"github.com/suisrc/fkssl/kube"
	"github.com/suisrc/fkssl/serve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateCaCmdApi(ctx *gin.Context) {
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
	infoKey := fmt.Sprintf("%s%s-%s", PK, co.Key, "info")
	info, err := cli.CoreV1().Secrets(kube.CurrentNamespace()).Get(ctx, infoKey, metav1.GetOptions{})
	if err != nil {
		serve.Error(ctx, 500, "KUBE-INFO-ERROR", err.Error())
		return
	}
	if tkn, ok := info.Data["token"]; ok && string(tkn) != co.Token { // 必须存在，不存在，不可访问
		serve.Error(ctx, 400, "TOKEN-ERROR", "token error")
		return
	}
	// ==========================================================================

	config := cert.CertConfig{}
	if err := ctx.ShouldBindJSON(&config); err != nil {
		serve.Error(ctx, 400, "PARAM-BODY-ERROR", err.Error())
		return
	}

	if configStr, ok := info.Data["config"]; ok { // 配置已经存在
		//=======================================================================
		config2 := cert.CertConfig{}
		if err := json.Unmarshal(configStr, &config2); err != nil {
			serve.Error(ctx, 500, "KUBE-CONFIG-ERROR", err.Error())
			return
		}
		// 合并配置, 更新配置
		if update := config2.Merge(config); update {
			info.Data["config"] = []byte(config2.String())
			info, err = cli.CoreV1().Secrets(kube.CurrentNamespace()).Update(ctx, info, metav1.UpdateOptions{})
			if err != nil {
				serve.Error(ctx, 500, "KUBE-UPDATE-ERROR", err.Error())
				return
			}
		}
		config = config2
	} else { // 配置不存在
		info.Data["config"] = []byte(config.String())
		info, err = cli.CoreV1().Secrets(kube.CurrentNamespace()).Update(ctx, info, metav1.UpdateOptions{})
		if err != nil {
			serve.Error(ctx, 500, "KUBE-UPDATE-ERROR", err.Error())
			return
		}
	}
	// ==========================================================================
	if crt, ok := info.Data["ca.crt"]; ok {
		serve.Success(ctx, string(crt))
		return // 证书已经存在，立即返回
	}
	// 证书不存在，需要重写构建证书
	ca, err := cert.CreateCA(config)
	if err != nil {
		serve.Error(ctx, 500, "CA-CREATE-ERROR", err.Error())
		return
	}
	// ==========================================================================
	info.Data["ca.crt"] = []byte(ca.Crt)
	info.Data["ca.key"] = []byte(ca.Key)
	// 求ca.Key的md5值
	md5CaKey, _ := hashMd5([]byte(ca.Key))
	info.Data["prefix"] = []byte(fmt.Sprintf("%s%s-%s-", PK, co.Key, md5CaKey[:8]))
	_, err = cli.CoreV1().Secrets(kube.CurrentNamespace()).Update(ctx, info, metav1.UpdateOptions{})
	if err != nil {
		serve.Error(ctx, 500, "KUBE-UPDATE-ERROR", err.Error())
		return
	}
	serve.Success(ctx, ca.Crt)
}

// ==============================================================================

func bingQuery(ctx *gin.Context) (*SslQueryCO, bool) {
	co := &SslQueryCO{}
	if err := ctx.ShouldBindQuery(co); err != nil {
		serve.Error(ctx, 400, "PARAM-QUERY-ERROR", err.Error())
		return nil, false
	}
	if co.Key == "" {
		serve.Error(ctx, 400, "KEY-ERROR", "key is empty")
		return nil, false
	}
	return co, true
}

// MD5Hash MD5哈希值
func hashMd5(b []byte) (string, error) {
	h := md5.New()
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
