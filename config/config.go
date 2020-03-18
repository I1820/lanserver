/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 30-01-2019
 * |
 * | File Name:     config.go
 * +===============================================
 */

package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	// Config holds all lanserver component configurations
	Config struct {
		Debug    bool
		Database Database `mapstructure:"database"`
		App      struct {
			Broker struct {
				Addr string
			}
		}
		Node struct {
			Broker struct {
				Addr string
			}
		}
	}

	// Database holds database configuration
	Database struct {
		URL  string `mapstructure:"url"`
		Name string `mapstructure:"name"`
	}
)

// New reads configuration with viper
func New() Config {
	var instance Config

	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadConfig(bytes.NewBufferString(Default)); err != nil {
		logrus.Fatalf("fatal error loading **default** config array: %s \n", err)
	}

	v.SetConfigName("config")

	if err := v.MergeInConfig(); err != nil {
		switch err.(type) {
		default:
			logrus.Fatalf("fatal error loading config file: %s \n", err)
		case viper.ConfigFileNotFoundError:
			logrus.Infof("no config file found. Using defaults and environment variables")
		}
	}

	v.SetEnvPrefix("i1820_lanserver")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.UnmarshalExact(&instance); err != nil {
		logrus.Infof("configuration: %s", err)
	}
	fmt.Printf("Following configuration is loaded:\n%+v\n", instance)

	return instance
}
