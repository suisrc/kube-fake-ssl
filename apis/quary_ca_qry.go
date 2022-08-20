package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/fkssl/kube"
	"github.com/suisrc/fkssl/serve"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func QuaryCaQryApi(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		serve.Error(ctx, 400, "KEY-ERROR", "key is empty")
		return
	}
	cli, err := kube.GetClient()
	if err != nil {
		serve.Error(ctx, 200, "KUBE-ERROR", err.Error())
		return
	}
	infoKey := fmt.Sprintf("%s%s-%s", PK, key, "info")
	info, err := cli.CoreV1().Secrets(kube.CurrentNamespace()).Get(ctx, infoKey, metav1.GetOptions{})
	if err != nil {
		serve.Error(ctx, 200, "KUBE-INFO-ERROR", err.Error())
		return
	}
	crt, ok := info.Data["ca.crt"]
	if !ok {
		serve.Error(ctx, 500, "CA-NOT-FOUND", "CA证书不存在")
		return
	}
	serve.Success(ctx, string(crt))
}

func QuaryCaQryTxtApi(ctx *gin.Context) {
	key := ctx.Query("key")
	if key == "" {
		ctx.Status(400)
		ctx.Writer.WriteString("key is empty")
		ctx.Abort()
		return
	}
	cli, err := kube.GetClient()
	if err != nil {
		ctx.Status(500)
		ctx.Writer.WriteString(err.Error())
		ctx.Abort()
		return
	}
	infoKey := fmt.Sprintf("%s%s-%s", PK, key, "info")
	info, err := cli.CoreV1().Secrets(kube.CurrentNamespace()).Get(ctx, infoKey, metav1.GetOptions{})
	if err != nil {
		ctx.Status(500)
		ctx.Writer.WriteString(err.Error())
		ctx.Abort()
		return
	}
	crt, ok := info.Data["ca.crt"]
	if !ok {
		ctx.Status(500)
		ctx.Writer.WriteString("CA证书不存在")
		ctx.Abort()
		return
	}
	ctx.Writer.Write(crt)
	ctx.Abort()
}
