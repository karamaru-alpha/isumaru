package constant

const (
	DefaultAccessLogPath    = "pkg/agent/log/access.log"
	DefaultSlowQueryLogPath = "pkg/agent/log/slow-query.log"

	IsumaruSlowQueryLogDirFormat = "pkg/isumaru/log/%s/slowquerylog"
	IsumaruAccessLogDirFormat    = "pkg/isumaru/log/%s/accesslog"
	SlpConfigPath                = "config/slp.yaml"
	AlpConfigPath                = "config/alp.yaml"
)
