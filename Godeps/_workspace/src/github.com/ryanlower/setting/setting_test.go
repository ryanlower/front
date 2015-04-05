package setting

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleConfig struct {
	AAA string `default:"A"`
	BBB string `env:"BBB"`
	CCC string `env:"CCC" default:"C"`
}

type NestedConfig struct {
	Service struct {
		Port string `default:"1234"`
		Auth struct {
			Key    string `default:"key"`
			Secret string `env:"SERVICE_SECRET"`
		}
	}
}

func TestSettingFromDefault(t *testing.T) {
	config := new(SimpleConfig)
	Load(config)

	assert.Equal(t, config.AAA, "A")
}

func TestSettingFromEnv(t *testing.T) {
	os.Setenv("BBB", "three b's")
	defer os.Unsetenv("BBB")

	config := new(SimpleConfig)
	Load(config)

	assert.Equal(t, config.BBB, "three b's")
}

func TestSettingFromEnvDefaultFallback(t *testing.T) {
	config := new(SimpleConfig)
	Load(config)

	assert.Equal(t, config.CCC, "C")
}

func TestSettingNestedConfig(t *testing.T) {
	os.Setenv("SERVICE_SECRET", "super secret")
	defer os.Unsetenv("SERVICE_SECRET")

	config := new(NestedConfig)
	Load(config)

	assert.Equal(t, config.Service.Port, "1234")
	assert.Equal(t, config.Service.Auth.Key, "key")
	assert.Equal(t, config.Service.Auth.Secret, "super secret")
}
