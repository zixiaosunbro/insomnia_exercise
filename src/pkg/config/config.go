package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"insomnia/src/pkg/config/impl"
	"insomnia/src/pkg/utils"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	appConfig unsafe.Pointer
	ENV       = os.Getenv("CUSTOM_RUNTIME_ENV")
	InDocker  = os.Getenv("IN_DOCKER") != ""
	conf      Configer
)

type Configer = impl.Configer

type DBConf struct {
	Driver        string
	Master        string
	Slave         string
	MaxLifetime   time.Duration `config:"max_lifetime"`
	MaxOpenConns  int           `config:"max_open_conns"`
	MaxIdleConns  int           `config:"max_idle_conns"`
	EnableLog     bool          `config:"enable_log"`
	EnableMetrics bool          `config:"enable_metrics"`
	AutoCommit    bool          `config:"auto_commit"`
}

type RedisConf struct {
	Uri             string
	Name            string
	Auth            string        `config:"auth"`
	ConnectTimeout  time.Duration `config:"connect_timeout"`
	ReadTimeout     time.Duration `config:"read_timeout"`
	WriteTimeout    time.Duration `config:"write_timeout"`
	PoolMaxActive   int           `config:"pool_max_active"`
	PoolIdleTimeout time.Duration `config:"pool_idle_timeout"`
	MaxRetries      int           `config:"max_retries"`
	DB              int           `config:"db"`
}

type AppConfig struct {
	AppID   string `config:"app_id"`
	Plugins struct {
		DB         bool
		Redis      bool
		Q          bool
		Async      bool
		Kafka      bool
		Tablestore bool
	}
	DB    map[string]*DBConf
	Redis map[string]*RedisConf
}

type ViperWrap struct {
	*viper.Viper
	mutex sync.RWMutex
}

func (w *ViperWrap) MergeFromFile(file string) error {
	cfg := viper.New()
	cfg.SetConfigFile(file)
	err := cfg.ReadInConfig()
	if err != nil {
		log.Infof("read config file error: %v", err)
		return err
	}
	w.mutex.Lock()
	defer w.mutex.Unlock()
	err = w.Viper.MergeConfigMap(cfg.AllSettings())
	return err
}
func (w *ViperWrap) Unmarshal(v interface{}) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Viper.Unmarshal(v)
}

func NewViperWrap(cfg *viper.Viper) *ViperWrap {
	w := new(ViperWrap)
	w.Viper = cfg
	return w
}

func GetConfig() *AppConfig {
	p := atomic.LoadPointer(&appConfig)
	return (*AppConfig)(p)
}

func SetConfig() {
	var cfg AppConfig
	err := conf.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unmarshal config error: %v", err)
		return
	}
	atomic.StorePointer(&appConfig, unsafe.Pointer(&cfg))
}

func init() {
	cfgInit()
	// app.yaml start
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("app.yaml not found: " + err.Error())
		} else {
			panic("app.yaml read error: " + err.Error())
		}
	}
	conf = NewViperWrap(viper.GetViper())
	configFile := fmt.Sprintf("config/%s.yaml", ENV)
	if utils.JudgeFileExist(configFile) {
		err = conf.MergeFromFile(configFile)
		if err != nil {
			panic("merge config error: " + err.Error())
		}
	}
	SetConfig()
}
