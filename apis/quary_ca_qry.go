package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/fkssl/kube"
	"github.com/suisrc/fkssl/serve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func QuaryCaQryApi(ctx *gin.Context) {
	co, ok := bingQuery(ctx)
	if !ok {
		return
	}
	cli, err := kube.GetClient()
	if err != nil {
		serve.Error(ctx, 500, "KUBE-ERROR", err.Error())
		return
	}
	infoKey := fmt.Sprintf("%s%s-%s", PK, co.Key, "info")
	info, err := cli.CoreV1().Secrets("").Get(ctx, infoKey, metav1.GetOptions{})
	if err != nil {
		serve.Error(ctx, 500, "KUBE-INFO-ERROR", err.Error())
		return
	}
	crt, ok := info.StringData["ca.crt"]
	if !ok {
		serve.Error(ctx, 200, "CA-NOT-FOUND", "CA证书不存在")
		return
	}
	serve.Success(ctx, crt)
}
