package config

type SQLDBSetting struct {
	SqlDSN              string
	SqlMaxOpenConns     int
	SqlMaxIdleConns     int
	SqlConnsMaxLifetime int
}

func (s *SQLDBSetting) DSN() string {
	return s.SqlDSN
}

func (s *SQLDBSetting) MaxOpenConns() int {
	return s.SqlMaxOpenConns
}

func (s *SQLDBSetting) MaxIdleConns() int {
	return s.SqlMaxIdleConns
}

func (s *SQLDBSetting) ConnsMaxLifetime() int {
	return s.SqlConnsMaxLifetime
}
