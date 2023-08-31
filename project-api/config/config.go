package config //viper库配置文件

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"test.com/project-common/logs"
)

var C = InitConfig()

type Config struct {
	viper      *viper.Viper
	SC         *ServerConfig //.
	GC         *GrpcConfig   //.
	EtcdConfig *EtcdConfig
}

type ServerConfig struct {
	Name string
	Addr string
}

type GrpcConfig struct {
	Name string
	Addr string
}
type EtcdConfig struct {
	Addrs []string
}

func InitConfig() *Config {
	config := &Config{viper: viper.New()}

	config.viper.SetConfigName("config")
	config.viper.SetConfigType("yaml")
	config.viper.AddConfigPath("/etc/ms_project/user")
	workDir, _ := os.Getwd()
	config.viper.AddConfigPath(workDir + "/config")

	err := config.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}

	config.ReadServerConfig() //读取server配置
	config.InitZapLog()
	config.ReadEtcdConfig()
	return config
}

func (c *Config) ReadServerConfig() { //读取server配置
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) InitZapLog() {
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	if err := c.viper.UnmarshalKey("etcd.addrs", &addrs); err != nil {
		log.Fatalln(err)
	}
	ec.Addrs = addrs
	c.EtcdConfig = ec
}
