package utils

import (
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"path/filepath"
	"reflect"
	"strings"
)

// Vconfig is just a wrapper around viper Config
type Vconfig struct {
	*viper.Viper
}

// FromKey creates a deep copy of the Vconfig section
// this fixes viper bug when calling Sub() for hcl format
func (c *Vconfig) FromKey(key string) *Vconfig {
	key = strings.ToLower(key)

	subv := viper.New()
	data := c.Get(key)
	if data == nil {
		return nil
	}

	if reflect.TypeOf(data).Kind() == reflect.Map { // never happens in fact, always []map[string]interface{}
		subv.MergeConfigMap(cast.ToStringMap(data))
	} else if reflect.TypeOf(data).Kind() == reflect.Slice {
		sl := cast.ToSlice(data)
		if len(sl) != 1 {
			Throwf("error getting section %v", key)
		}
		subv.MergeConfigMap(cast.ToStringMap(sl[0]))
	} else {
		Throwf("error getting section %v", key)
	}

	return &Vconfig{
		Viper: subv,
	}
}

func (c *Vconfig) GetCfg() map[string]interface{} {
	return c.AllSettings()
}

func (c *Vconfig) GetValue(key string) interface{} {
	key = strings.ToLower(key)
	return c.Get(key)
}

func (c *Vconfig) GetFloat(key string) float64 {
	key = strings.ToLower(key)
	return c.GetFloat64(key)
}

func (c *Vconfig) GetFloatDefault(key string, dflt float64) float64 {
	key = strings.ToLower(key)
	c.SetDefault(key, dflt)
	return c.GetFloat64(key)
}

func (c *Vconfig) GetInt(key string) int64 {
	key = strings.ToLower(key)
	return c.GetInt64(key)
}

func (c *Vconfig) GetIntDefault(key string, dflt int64) int64 {
	key = strings.ToLower(key)
	c.SetDefault(key, dflt)
	return c.GetInt64(key)
}

func (c *Vconfig) GetBoolDefault(key string, dflt bool) bool {
	key = strings.ToLower(key)
	c.SetDefault(key, dflt)
	return c.GetBool(key)
}

func (c *Vconfig) GetString(key string) *string {
	key = strings.ToLower(key)
	rs := c.Viper.GetString(key)
	return &rs
}

func (c *Vconfig) GetStringDefault(key string, defaultVal string) *string {
	key = strings.ToLower(key)
	c.SetDefault(key, defaultVal)
	rs := c.Viper.GetString(key)
	return &rs
}

// ReadConfig reads json Vconfig. Throws if something is wrong
func (c *Vconfig) ReadConfig(filename string) *Vconfig {

	if c.Viper == nil {
		c.Viper = viper.GetViper()
	}

	onlyName := strings.TrimRight(strings.Replace(filepath.Base(filename), filepath.Ext(filename), "", 1), ".")
	c.SetConfigName(onlyName) // name of config file (without extension) -- what a f innovation!
	c.SetConfigType(strings.ToLower(strings.TrimLeft(filepath.Ext(filename), ".")))
	c.AddConfigPath(filepath.Dir(filename))
	err := c.ReadInConfig()
	if err != nil {
		Throwf("fatal error reading config file: %w", err)
	}

	return c
}

// WriteConfigX writes a section of a Vconfig back to file
func (c *Vconfig) WriteConfigX() error {
	Throwf("WriteConfigX not implemented yet")
	return nil
}
