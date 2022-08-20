package apis

const PK = "fkc-"

type SslQueryCO struct {
	Token      string   `form:"token"`
	Key        string   `form:"key"`
	Kind       int      `form:"kind"`
	CommonName string   `form:"cn"`
	Profile    string   `form:"profile"`
	Domains    []string `form:"domain"`
}
