package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

func ProvideChildConfig(nodes... string) (*viper.Viper, error) {
	var err error
	v := viper.New()
	for _, n := range nodes {
		vSettings, ok := viper.AllSettings()[n].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s settings not found", n)
		}
		err = v.MergeConfigMap(vSettings)
		if err != nil {
			return nil, fmt.Errorf("config not merged: %w", err)
		}
	}
	appSettings, ok := viper.AllSettings()["app"].(map[string]interface{})
	if !ok {
		return nil, errors.New("app settings not found")
	}
	globalSettings, ok := viper.AllSettings()["global"].(map[string]interface{})
	if !ok {
		return nil, errors.New("global settings not found")
	}
	err = v.MergeConfigMap(globalSettings)
	if err != nil {
		return nil, err
	}
	err = v.MergeConfigMap(appSettings)
	if err != nil {
		return nil, err
	}
	return v, nil
}
