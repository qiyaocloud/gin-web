package tests

// 初始化单元测试必须的环境, 有些方法无需启动应用先测试
func InitTestEnv() {
	// 初始化配置
	Config()

	// 初始化日志
	Logger()

	// 初始化数据库
	Redis()
	Mysql()
}
