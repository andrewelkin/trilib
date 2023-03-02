package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/andrewelkin/trilib/utils/logger"
	"io/fs"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/titanous/json5"

	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sts"
	"golang.org/x/crypto/ssh/terminal"
)

// Config is just a map of config variable-value pairs.
type Config struct {
	FileName     string                 // filename
	RO           bool                   // file is r/o
	lastModified time.Time              // last mod timestamp
	cfg          map[string]interface{} // config map

	protect       sync.Mutex
	igniterValues map[string]string
}

// WriteToLog writes current "PublicCfg" map to specified logger
func (c *Config) WriteToLog(l logger.Logger) {
	dataIndent, err := json.Marshal(c.cfg)
	if err != nil {
		logger.GetGlobalLogger().Errorf("TRT", "Writing config to log failed! err = %v\n", err)
		return
	}
	logger.GetGlobalLogger().Infof(
		"_TRT",
		`{"current_config":%v}`,
		string(dataIndent),
	)
}

func (c *Config) FromKeyAsString(key string) (string, error) {
	c.protect.Lock()
	defer c.protect.Unlock()
	res := c.GetValue(key)
	if res == nil {
		return "", fmt.Errorf("no section found in config: %v", key)
	}
	s, err := json.Marshal(res)
	return string(s), err
}

// FromKey creates a deep copy of the config section
func (c *Config) FromKey(key string) *Config {
	c.protect.Lock()
	defer c.protect.Unlock()
	res := c.GetValue(key)
	if res == nil {
		return nil
	}

	if newConf, ok := res.(*Config); ok {
		return &Config{
			cfg:           newConf.cfg,
			FileName:      newConf.FileName,
			igniterValues: newConf.igniterValues,
			lastModified:  newConf.lastModified,
		}
	}

	if resmap, ok := res.(map[string]interface{}); ok {
		return &Config{
			cfg:           resmap,
			FileName:      c.FileName,
			igniterValues: c.igniterValues,
		}
	}

	if v, ok := res.(string); ok {
		if strings.HasPrefix(v, "^") || strings.HasPrefix(v, "file://") || strings.HasPrefix(v, "rfile://") {
			ndx := 1
			ro := false
			if strings.HasPrefix(v, "file://") {
				ndx = 7
			} else if strings.HasPrefix(v, "rfile://") {
				ndx = 8
				ro = true

			}
			v = strings.TrimLeft(v[ndx:], " ")

			if !ro && !FlagExists(v) { // empty, create one
				var file, err = os.OpenFile(v, os.O_RDWR|os.O_CREATE, 0644)
				if err != nil {
					Throwf("file %s not exists and could not create an empty one", v)
				}

				file.WriteString("{}\n")
				file.Close()
			}
			cfg := &Config{
				cfg:           map[string]interface{}{},
				FileName:      v,
				RO:            ro,
				igniterValues: map[string]string{},
				lastModified:  time.Unix(3600000, 0), // some moment long ago
			}
			cfg.ReadConfigX()
			return cfg
		}
	}

	Throwf("error getting section %v", key)
	return nil
}
func (c *Config) SetCfg(cfg map[string]interface{}) *Config {
	c.cfg = cfg
	return c
}

func (c *Config) GetCfg() map[string]interface{} {
	return c.cfg
}

func (c *Config) GetValue(key string) interface{} {
	res, _ := c.cfg[key]
	return res
}

func (c *Config) GetBool(key string) bool {
	res, _ := c.cfg[key]
	rb := res.(bool)
	return rb
}

func (c *Config) GetFloat(key string) float64 {
	res, _ := c.cfg[key]
	rb := res.(float64)
	return rb
}

func (c *Config) GetFloatDefault(key string, dflt float64) float64 {
	res, ok := c.cfg[key]
	if !ok {
		return dflt
	}
	rb := res.(float64)
	return rb
}

func (c *Config) GetInt(key string) int64 {
	res, _ := c.cfg[key]
	rb := int64(res.(float64))
	return rb
}

func (c *Config) GetIntDefault(key string, dflt int64) int64 {
	res, ok := c.cfg[key]
	if !ok {
		return dflt
	}
	if rb, ok := res.(int); ok {
		return int64(rb)
	}
	if rb, ok := res.(float64); ok {
		return int64(rb)
	}
	if rb, ok := res.(float32); ok {
		return int64(rb)
	}

	return dflt
}

func (c *Config) GetBoolDefault(key string, dflt bool) bool {
	res, ok := c.cfg[key]
	if !ok {
		return dflt
	}
	rb := res.(bool)
	return rb
}

func (c *Config) GetString(key string) *string {
	res, _ := c.cfg[key]
	if res == nil {
		return nil
	}
	rs := res.(string)

	// secret value from igniter
	if strings.HasPrefix(rs, "#") {
		if c.igniterValues != nil { // if igniter config is used
			v := rs[1:]
			tailVal := ""
			ndx := strings.Index(v, "%")
			if ndx > 0 {
				tailVal = v[ndx+1:]
				v = v[0:ndx]
			}
			if val, hasVal := c.igniterValues[v]; hasVal {
				res := val + tailVal
				return &res
			}
			Throwf("igniter syntax used (%q), but not value saved for key: %q", "#", rs[1:])
		}
		Throwf("igniter syntax used (%q), but no igniter secrets stored", "#")
	}

	// replace value from environment
	if strings.HasPrefix(rs, "$") {
		envStr := os.Getenv(rs[1:])
		return &envStr
	}

	// replace value with environment variable, windows %xx% style
	if strings.HasPrefix(rs, "%") {
		ndx := strings.Index(rs[1:], "%")
		if ndx < 0 {
			envStr := os.Getenv(rs[1:])
			return &envStr
		}

		envStr := os.Getenv(rs[1 : ndx+1])
		d := envStr + rs[ndx+2:]
		return &d
	}
	return &rs
}

func (c *Config) GetStringDefault(key string, defaultVal string) *string {
	if r := c.GetString(key); r != nil {
		return r
	}
	return &defaultVal
}

func (c *Config) GetStringList(key string) []string {
	rawList := c.GetString(key)
	return strings.Split(*rawList, ",")
}

func (c *Config) GetListOfStrings(key string) []string {
	var result []string
	st := c.GetValue(key)
	if st != nil {
		for _, s := range st.([]interface{}) {
			result = append(result, s.(string))
		}
	}
	return result
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func getUserInput(prompt string) string {
	fmt.Printf("%s: ", prompt)
	var a string
	fmt.Scanln(&a)
	return a
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func replaceIncludes(c *Config) *Config {

	var onceMore = true
	for onceMore {
		onceMore = false

		m := c.GetCfg()

		for k, v := range m {
			s, ok := v.(string)
			if ok && (strings.HasPrefix(s, "^") || strings.HasPrefix(s, "file://") || strings.HasPrefix(s, "rfile://")) {
				ndx := 1
				ro := false
				if strings.HasPrefix(s, "file://") {
					ndx = 7
				} else if strings.HasPrefix(s, "rfile://") {
					ndx = 8
					ro = true
				}
				s = strings.TrimLeft(s[ndx:], " ")
				cfg := &Config{
					FileName: s,
					RO:       ro,
				}
				cfg.ReadConfigX()
				m[k] = cfg
				// can I modify the map inside the loop?  not guaranteed
				// https://stackoverflow.com/questions/68639365/modifying-map-while-iterating-over-it-in-go
				onceMore = true
				break
			}
		}
	}

	return c
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/
func replaceInteractiveFields(c *Config) *Config {

	secretKeys := make(map[string]int)
	var keys, vars []string
	m := c.GetCfg()
	for k, v := range m {
		s, ok := v.(string)
		if ok {
			if strings.HasPrefix(s, "!") {
				secretKeys[k] = 0
				keys = append(keys, k)
				vars = append(vars, s[1:])
			}
			if strings.HasPrefix(s, "**") {
				secretKeys[k] = 1
				keys = append(keys, k)
				vars = append(vars, s[2:])
			}
			if strings.HasPrefix(s, "$") {
				secretKeys[k] = 2
				keys = append(keys, k)
				vars = append(vars, s[1:])
			}
			if strings.HasPrefix(s, "%") {
				secretKeys[k] = 3
				keys = append(keys, k)
				vars = append(vars, s[1:])
			}
		}
	}

	for i := 0; i < len(keys); i++ {
		k, v := keys[i], vars[i]

		var val string

		switch secretKeys[k] {
		case 0:
			val = getUserInput(v)
			fmt.Printf("Assigning %s <-- %s\n", k, fmt.Sprintf("'%v'", val))
		case 1:
			ndx := strings.Index(v, "%")
			prmpt := ""
			rest := ""
			if ndx < 0 {
				prmpt = v
			} else {
				prmpt = v[0:ndx]
				rest = v[ndx+1:]
			}
			fmt.Printf("%s: ", prmpt)
			raw, _ := terminal.ReadPassword(int(syscall.Stdin))
			fmt.Println()
			val = string(raw) + rest
		case 2:
			val = os.Getenv(v)
		case 3:
			ndx := strings.Index(v, "%")
			if ndx < 0 {
				val = os.Getenv(v)
			} else {
				envStr := os.Getenv(v[0:ndx])
				val = envStr + v[ndx+1:]
			}
		}
		m[keys[i]] = val
	}

	var cfgs []string
	for k, v := range m {
		_, ok := v.(map[string]interface{})
		if ok {
			cfgs = append(cfgs, k)
		}
	}

	for _, k := range cfgs {
		m1 := c.FromKey(k)
		replaceInteractiveFields(m1)
	}
	return c
}

// returns true if the file modified since last ReadFile
func (c *Config) ModifiedQ() bool {
	fileStats, _ := os.Stat(c.FileName)
	return c.lastModified != fileStats.ModTime()
}

// ReadConfig reads json config. Throws if something is wrong
func (c *Config) ReadConfig(filename string) *Config {

	absolutePath, err := filepath.Abs(filename)
	if err != nil {
		Throw("Error converting relative to absolute path before trying to read file " + filename)
	}

	c.FileName = absolutePath
	return c.ReadConfigX()
}

// config for tests
func (c *Config) BuildFromBytes(config []byte) *Config {
	err := json5.Unmarshal(config, &c.cfg)
	if err != nil {
		Throw("Error unmarshalling bytes: " + err.Error())
	}
	return c
}

// WriteConfigX writes a section of a config back to file
func (c *Config) WriteConfigX() error {

	if c.RO {
		return fmt.Errorf("config is read-only, can't save to %s", c.FileName)
	}
	data, err := json.MarshalIndent(c.cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal config for saving")
	}
	err = ioutil.WriteFile(c.FileName, data, fs.ModePerm)
	return err
}

// ReadConfigX reads json config from the stored filename. Throws if something is wrong
func (c *Config) ReadConfigX() *Config {

	fileBytes, err := ioutil.ReadFile(c.FileName)
	if err != nil {
		Throw("Error reading file " + c.FileName + ": " + err.Error())
	}
	fileStats, _ := os.Stat(c.FileName)
	c.lastModified = fileStats.ModTime()
	c.protect.Lock()

	err = json5.Unmarshal(fileBytes, &c.cfg)
	c.protect.Unlock()
	if err != nil {

		Throw("Error in json file " + c.FileName + " error:" + err.Error())
	}

	replaceIncludes(c)
	replaceInteractiveFields(c)
	if igniterConfig := c.FromKey("igniter"); igniterConfig != nil {
		loadIgniterValues(c, igniterConfig)
	}
	return c
}

func loadIgniterValues(main *Config, igniter *Config) {
	resourceID := igniter.GetString("resource_id")
	rawToken := igniter.GetString("token")
	region := igniter.GetStringDefault("region", "us-east-1")

	payload, err := base64.StdEncoding.DecodeString(*rawToken)
	if err != nil {
		Throwf("failed to decode b64 encoded payload: %v", err)
	}

	creds := new(sts.Credentials)
	if err := json.NewDecoder(bytes.NewBuffer(payload)).Decode(creds); err != nil {
		Throwf("failed to parse credentials: %v", err)
	}

	cred := credentials.NewStaticCredentials(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)
	session, err := session.NewSession(&aws.Config{Credentials: cred, Region: region})
	if err != nil {
		Throwf("failed to authenticate with AWS: %v", err)
	}

	secretsManager := secretsmanager.New(session)
	rawOutput, err := secretsManager.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId:     resourceID,
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	})
	if rawOutput.SecretString == nil {
		Throwf("failed to get secret (no value stored): %v", err.Error())
	}

	igniterValues := make(map[string]string)
	if err := json.Unmarshal([]byte(*rawOutput.SecretString), &igniterValues); err != nil {
		Throwf("failed to decode key response: %v", err)
	}
	main.igniterValues = igniterValues
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/
