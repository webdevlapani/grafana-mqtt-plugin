package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ini "gopkg.in/ini.v1"
)

const ConfigFile = "config.ini"

type Cfg struct {
	Raw    *ini.File

	Service string
	Zone    string
	Env     string

	S3BucketPrefix string

	AWSAccountId string
	AWSAccessKey string
	AWSSecretKey string

	Storage map[string]int
}

func applyEnvVariableOverrides(file *ini.File) error {
	for _, section := range file.Sections() {
		for _, key := range section.Keys() {
			envKey := envKey(section.Name(), key.Name())
			envValue := os.Getenv(envKey)

			if len(envValue) > 0 {
				key.SetValue(envValue)
			}
		}
	}

	return nil
}

func envKey(sectionName string, keyName string) string {
	sN := strings.ToUpper(strings.Replace(sectionName, ".", "_", -1))
	sN = strings.Replace(sN, "-", "_", -1)
	kN := strings.ToUpper(strings.Replace(keyName, ".", "_", -1))
	envKey := fmt.Sprintf("GF_%s_%s", sN, kN)
	return envKey
}

func (cfg *Cfg) loadIni() (*ini.File, error) {
	var err error

	home, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}

	// check if home is correct
	if !pathExists(filepath.Join(home, "config.ini")) {
		// try down one path
		if !pathExists(filepath.Join(home, "../config.ini")) {
			return nil, fmt.Errorf("Could not find %q\n", filepath.Join(home, "config.ini"))
		}
               	home = filepath.Join(home, "../")
	}

	// load config
	parsedFile, err := ini.Load(filepath.Join(home, "config.ini"))
	if err != nil {
		return nil, err
	}

	parsedFile.BlockMode = false

	// apply environment overrides
	err = applyEnvVariableOverrides(parsedFile)
	if err != nil {
		return nil, err
	}

	return parsedFile, nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func loadConfig() *Cfg {
	cfg := &Cfg{
		Raw:    ini.Empty(),
        }

	iniFile, err := cfg.loadIni()
	if err != nil {
		fmt.Printf("Could not load %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}

	cfg.Raw = iniFile

	cfg.Service, err = valueAsString(iniFile.Section(""), "service", "Cloud")
	if err != nil {
		fmt.Printf("service invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}
	cfg.Zone, err = valueAsString(iniFile.Section(""), "zone", "Landing Zone")
	if err != nil {
		fmt.Printf("zone invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}
	cfg.Env, err = valueAsString(iniFile.Section(""), "environment", "alpha")
	if err != nil {
		fmt.Printf("enviroment invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}

	cfg.S3BucketPrefix, err = valueAsString(iniFile.Section("s3"), "bucket_prefix", "")
	if err != nil {
		fmt.Printf("aws.bucket_prefix invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}
	if cfg.S3BucketPrefix == "" {
		fmt.Printf("aws.bucket_prefix required in %q\n", ConfigFile)
		os.Exit(1)
	}

	cfg.AWSAccountId, err = valueAsString(iniFile.Section("aws"), "account_id", "")
	if err != nil {
		fmt.Printf("aws.account_id invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}
	if cfg.AWSAccountId == "" {
		fmt.Printf("aws.account_id required in %q\n", ConfigFile)
		os.Exit(1)
	}

	cfg.AWSAccessKey, err = valueAsString(iniFile.Section("aws"), "access_key", "")
	if err != nil {
		fmt.Printf("aws.access_key invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}

	cfg.AWSSecretKey, err = valueAsString(iniFile.Section("aws"), "secret_key", "")
	if err != nil {
		fmt.Printf("aws.secret_key invalid in %q: %v\n", ConfigFile, err)
		os.Exit(1)
	}

	return cfg
}

func valueAsString(section *ini.Section, keyName string, defaultValue string) (value string, err error) {
	defer func() {
		if err_ := recover(); err_ != nil {
			err = errors.New("Invalid value for key '" + keyName + "' in configuration file")
		}
	}()

	return section.Key(keyName).MustString(defaultValue), nil
}
