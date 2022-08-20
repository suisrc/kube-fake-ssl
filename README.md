# 说明

该工具只能运行在kubernetes集群中，不能独立运行
该工具使用secret作为存储，运行时候，需要secret读写权限

secret所有的key，统一使用fkc-开头 -> fake cert; 因此在该应用空间内的所有内容尽量避免使用改内容。
PS: 证书过期，如果被调用，会自动删除过期的证书，生产新的证书覆盖，否则过期证书不会被删除 

## 接口列表

### 接口测试
```rest
GET http://127.0.0.1/ping
```

### 健康检查
```rest
GET http://127.0.0.1/healthz
```

### 创建CA
如果已经存在，返回之前令牌的基本信息+crt&key
```rest
POST http://127.0.0.1/api/crt/v1/ca/init?token=&name=tst

{
  "CN": "Kubernetes",
  "key": {
    "size": 2048
  },
  "name": {
      "C": "CN",
      "ST": "Liaoning",
      "L": "Dalian",
      "O": "Kubernetes",
      "OU": "CA"
  },
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profile1": {
      "expiry": "876000h",
      "name": {
          "C": "CN",
          "ST": "Liaoning",
          "L": "Dalian",
          "O": "Kubernetes",
          "OU": "CA"
      }
    }
  }
}
```
{
    success: true,
    data: {
        cfg: {... POST body}
        crt: {... PEM}
        key: {... PEM}
    },
    traceId: "123456"
}

### 获取CA
PS: 获取ca内容时候，无需令牌，该内容可以理解为公共内容， 无需健全
```rest
GET http://127.0.0.1/api/crt/v1/ca/init?name=tst
```
{
    success: true,
    data: {crt.pem},
    traceId: "123456"
}

### 获取PEM
PS: domain,domains二选一，domains使用md5存储，不如domain直观; kind=1(如果没有，新增)
```rest
GET http://127.0.0.1/api/crt/v1/cert?token=&name=tst&domain=&domains=[h1,h2,h3]&profile=&kind=1
```
{
    success: true,
    data: {
        crt: {crt.pem}
        key: {key.pem}
    },
    traceId: "123456"
}

### 其他
删除CA： 暂时不支持，需要人工从kubernetes中的secret中删除
注销PEM：暂时不支持

## 错误
{
    success: false,
    errorCode: "ERROR-CODE",
    errorMessage: "异常说明",
    traceId: "123456"
}
