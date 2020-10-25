package config

import (
	"github.com/spf13/viper"
	"testing"
)

func TestProvideChildConfig(t *testing.T) {
	v := viper.New()
	v.Set("foo", "bar")
	v.MergeConfigMap(map[string]interface{}{"foo": "baz"})
	if v.GetString("foo") != "bar" {
		t.Errorf("want %s, got %s", "bar", v.GetString("foo"))
	}
	v.MergeConfigMap(map[string]interface{}{"quuz": "baz"})
	if v.GetString("quuz") != "baz" {
		t.Errorf("want %s, got %s", "baz", v.GetString("quuz"))
	}
}
