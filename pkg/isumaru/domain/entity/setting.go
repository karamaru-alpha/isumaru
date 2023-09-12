package entity

type Setting struct {
	Seconds            int32
	MainServerAddress  string
	AccessLogPath      string
	MysqlServerAddress string
	SlowQueryLogPath   string
}
