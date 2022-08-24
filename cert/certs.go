package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"
)

// 构建的证书是虚假的，默认有效期5年
func CreateCA(config CertConfig /*, notAfter time.Time*/) (SignResult, error) {

	profile := config.CaProfile
	subject := pkix.Name{ //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		CommonName:         config.CommonName,
		Country:            StrToArray(profile.SubjectName.Country),
		Province:           StrToArray(profile.SubjectName.Province),
		Locality:           StrToArray(profile.SubjectName.Locality),
		Organization:       StrToArray(profile.SubjectName.Organization),
		OrganizationalUnit: StrToArray(profile.SubjectName.OrganizationUnit),
	}
	// 过期时间
	notAfter := time.Now()
	if profile.Expiry == "" {
		notAfter = notAfter.Add(5 * 365 * 24 * time.Hour)
	} else if strings.HasSuffix(profile.Expiry, "h") {
		expiry, _ := strconv.Atoi(profile.Expiry[:len(profile.Expiry)-1])
		notAfter = notAfter.Add(time.Duration(expiry) * time.Hour)
	} else if strings.HasSuffix(profile.Expiry, "d") {
		expiry, _ := strconv.Atoi(profile.Expiry[:len(profile.Expiry)-1])
		notAfter = notAfter.Add(time.Duration(expiry) * 24 * time.Hour)
	} else {
		return SignResult{}, fmt.Errorf("invalid profile: ca, expiry: %s", profile.Expiry)
	}
	// 加密方式
	algorithm := x509.SHA256WithRSA
	if config.SignKey.Size >= 4096 {
		algorithm = x509.SHA512WithRSA
	} else if config.SignKey.Size >= 2048 {
		algorithm = x509.SHA384WithRSA
	}
	pkey, _ := rsa.GenerateKey(rand.Reader, config.SignKey.Size) //生成一对具有指定字位数的RSA密钥

	sermax := new(big.Int).Lsh(big.NewInt(1), 128) //把 1 左移 128 位，返回给 big.Int
	serial, _ := rand.Int(rand.Reader, sermax)     //返回在 [0, max) 区间均匀随机分布的一个随机值
	pder := x509.Certificate{
		SerialNumber: serial, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		IsCA:         true,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		SignatureAlgorithm:    algorithm, // SignatureAlgorithm 签名算法
	}

	//CreateCertificate基于模板创建一个新的证书
	//第二个第三个参数相同，则证书是自签名的
	//返回的切片是DER编码的证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &pder, &pder, &pkey.PublicKey, pkey)
	if err != nil {
		return SignResult{}, nil
	}
	crtBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pkey)})

	return SignResult{Crt: string(crtBytes), Key: string(keyBytes)}, nil
}

// CreateCertificate 创建一个证书，默认有效期2年
func CreateCert(config CertConfig, commonName, profileKey string, dns []string, ips []net.IP, caCrtPem, caKeyPem string) (SignResult, error) {
	profile, ok := config.SignProfiles[profileKey]
	if !ok {
		profile, ok = config.SignProfiles["default"]
	}
	if !ok {
		return SignResult{}, fmt.Errorf("no profile: %s", profileKey)
	}

	caCrtBlk, _ := pem.Decode([]byte(caCrtPem))
	if caCrtBlk == nil {
		return SignResult{}, fmt.Errorf("invalid ca.crt, pem")
	}
	caCrt, err := x509.ParseCertificate(caCrtBlk.Bytes)
	if err != nil {
		return SignResult{}, fmt.Errorf("invalid ca.crt, bytes")
	}
	caKeyBlk, _ := pem.Decode([]byte(caKeyPem))
	if caKeyBlk == nil {
		return SignResult{}, fmt.Errorf("invalid ca.key, pem")
	}
	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlk.Bytes)
	if err != nil {
		return SignResult{}, fmt.Errorf("invalid ca.key, bytes")
	}

	if commonName == "" {
		if len(dns) == 1 {
			commonName = dns[1]
		} else if len(ips) == 1 {
			commonName = ips[1]
		} else {
			commonName = config.CommonName
		}
	}

	subject := pkix.Name{ //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
		CommonName:         commonName,
		Country:            StrToArray(profile.SubjectName.Country),
		Province:           StrToArray(profile.SubjectName.Province),
		Locality:           StrToArray(profile.SubjectName.Locality),
		Organization:       StrToArray(profile.SubjectName.Organization),
		OrganizationalUnit: StrToArray(profile.SubjectName.OrganizationUnit),
	}
	// 过期时间
	notAfter := time.Now()
	if profile.Expiry == "" {
		notAfter = notAfter.Add(2 * 365 * 24 * time.Hour)
	} else if strings.HasSuffix(profile.Expiry, "h") {
		expiry, _ := strconv.Atoi(profile.Expiry[:len(profile.Expiry)-1])
		notAfter = notAfter.Add(time.Duration(expiry) * time.Hour)
	} else if strings.HasSuffix(profile.Expiry, "d") {
		expiry, _ := strconv.Atoi(profile.Expiry[:len(profile.Expiry)-1])
		notAfter = notAfter.Add(time.Duration(expiry) * 24 * time.Hour)
	} else {
		return SignResult{}, fmt.Errorf("invalid profile: %s, expiry: %s", profileKey, profile.Expiry)
	}
	// 加密方式
	algorithm := x509.SHA256WithRSA
	if config.SignKey.Size >= 4096 {
		algorithm = x509.SHA512WithRSA
	} else if config.SignKey.Size >= 2048 {
		algorithm = x509.SHA384WithRSA
	}
	pkey, _ := rsa.GenerateKey(rand.Reader, config.SignKey.Size) //生成一对具有指定字位数的RSA密钥

	sermax := new(big.Int).Lsh(big.NewInt(1), 128) //把 1 左移 128 位，返回给 big.Int
	serial, _ := rand.Int(rand.Reader, sermax)     //返回在 [0, max) 区间均匀随机分布的一个随机值
	pder := x509.Certificate{
		SerialNumber: serial,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		SignatureAlgorithm: algorithm,
		DNSNames:           dns,
		IPAddresses:        ips,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &pder, caCrt, &pkey.PublicKey, caKey)
	if err != nil {
		return SignResult{}, nil
	}
	crtBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pkey)})

	return SignResult{Crt: string(crtBytes), Key: string(keyBytes)}, nil
}

// IsExpired 判定正式否使过期
func IsExpired(caCrtPem string) (bool, error) {
	caCrtBlk, _ := pem.Decode([]byte(caCrtPem))
	if caCrtBlk == nil {
		return true, fmt.Errorf("invalid ca.crt, pem")
	}
	caCrt, err := x509.ParseCertificate(caCrtBlk.Bytes)
	if err != nil {
		return true, fmt.Errorf("invalid ca.crt, bytes")
	}

	if time.Now().After(caCrt.NotAfter) {
		return true, nil
	}
	return false, nil
}
