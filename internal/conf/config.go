package conf

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type GlobalConfiguration struct {
	// Credential
	Credential *Credential

	// Watch config
	Watcher *Watcher

	// ECS Instance Info
	ECS *ECS

	// Security Group Info
	SecurityGroup *SecurityGroup `split_words:"true"`
}

type Credential struct {
	Type            *string `json:"type,omitempty"`
	AccessKeyId     *string `json:"access_key_id,omitempty" split_words:"true"`
	AccessKeySecret *string `json:"access_key_secret,omitempty" split_words:"true"`
}

type Watcher struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	Interval *int64  `json:"interval,omitempty"`
	Path     *string `json:"path,omitempty"`
}

type ECS struct {
	InstanceId *string `json:"instance_id,omitempty" split_words:"true"`
	RegionId   *string `json:"region_id,omitempty" split_words:"true"`
	Endpoint   *string `json:"endpoint,omitempty"` // See: https://api.aliyun.com/product/Ecs
}

type SecurityGroup struct {
	Id *string `json:"id,omitempty"`
}

var (
	DefaultPrefix = "ALIYUN_SGMGR"
)

func NewConfig() *GlobalConfiguration {
	return &GlobalConfiguration{
		Credential:    &Credential{},
		Watcher:       &Watcher{},
		ECS:           &ECS{},
		SecurityGroup: &SecurityGroup{},
	}
}

// load .env file to environment variables
func LoadFile(filename string) error {
	var err error
	if filename != "" {
		err = godotenv.Overload(filename)
	} else {
		err = godotenv.Load()
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func LoadGlobalFromEnv() (config *GlobalConfiguration, err error) {
	config = NewConfig()
	err = envconfig.Process(DefaultPrefix, config)
	return config, err
}
