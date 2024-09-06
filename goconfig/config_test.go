package goconfig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jsalonl/go-config/v2/goconfig"
	"github.com/stretchr/testify/assert"
)

const (
	configFileYaml = "app.yaml"
)

func TestParseConfig(t *testing.T) {
	config := goconfig.NewGoConfig()

	assert.NotNil(t, config)
}

func TestParseConfigWithCustomUnmarshall(t *testing.T) {
	customUnmarshall := func(structure interface{}, content []byte) error {
		return nil
	}
	config := goconfig.NewGoConfig(customUnmarshall)

	assert.NotNil(t, config)
}

func TestLoadEnvSuccess(t *testing.T) {
	content := `APP_NAME=TestApp
APP_VERSION=1.0
`
	createEnvFile(t, content)
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err := config.LoadEnv()
	assert.NoError(t, err)

	assert.Equal(t, "TestApp", os.Getenv("APP_NAME"))
	assert.Equal(t, "1.0", os.Getenv("APP_VERSION"))

	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
	removeEnvFile(t)
}

func TestLoadEnvSuccessWithComments(t *testing.T) {
	content := `APP_NAME=TestApp
# This is a comment
APP_VERSION=1.0
`
	createEnvFile(t, content)
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err := config.LoadEnv()
	assert.NoError(t, err)

	assert.Equal(t, "TestApp", os.Getenv("APP_NAME"))
	assert.Equal(t, "1.0", os.Getenv("APP_VERSION"))

	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
	removeEnvFile(t)
}

func TestLoadEnvFailOpenDir(t *testing.T) {
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err := config.LoadEnv("nonexistent")
	assert.Error(t, err)
}

func TestLoadEnvFailMatchString(t *testing.T) {
	content := `APP_NAME:=TestApp
APP_VERSION=1.0
`
	createEnvFile(t, content)
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err := config.LoadEnv()
	assert.Error(t, err)
	assert.ErrorIs(t, err, goconfig.ErrInvalidEnvFormat)

	removeEnvFile(t)
}

func TestLoadEnvFailSplitN(t *testing.T) {
	content := `APP_NAME:TestApp
APP_VERSION:1.0
`
	createEnvFile(t, content)
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err := config.LoadEnv()
	assert.Error(t, err)
	assert.ErrorIs(t, err, goconfig.ErrInvalidEnvFormat)

	removeEnvFile(t)
}

func TestParseConfigSuccessYAML(t *testing.T) {
	content := `App:
  name: AppName
  version: 1.0
  log_level: 2
storage:
  master:
    name: MASTER_CONNECTION
    host: master-pg.localhost
    port: 5432
    user: user
    password: password
    database: db
  slave:
    name: SLAVE_CONNECTION
    host: slave-pg.localhost
    port: 5432
    user: user
    password: password
    database: db
`
	dir, file := createConfigFile(t, content)

	var yamlCfg AppConfig
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err := config.ParseConfig(&yamlCfg, "App", dir)
	assert.NoError(t, err)
	assert.Equal(t, "AppName", yamlCfg.App.Name)
	assert.Equal(t, "1.0", yamlCfg.App.Version)

	_ = os.Remove(filepath.Join(dir, file))
}

func TestParseConfigSuccessWithEnvVariables(t *testing.T) {
	err := os.Setenv("APP_NAME", "TestApp")
	assert.NoError(t, err)

	content := `App:
  name: ${APP_NAME}
  version: 1.0
  log_level: 2
storage:
  master:
    name: MASTER_CONNECTION
    host: master-pg.localhost
    port: 5432
    user: user
    password: password
    database: db
  slave:
    name: SLAVE_CONNECTION
    host: slave-pg.localhost
    port: 5432
    user: user
    password: password
    database: db
`
	dir, file := createConfigFile(t, content)

	var yamlCfg AppConfig
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err = config.ParseConfig(&yamlCfg, "App", dir)
	assert.NoError(t, err)
	assert.Equal(t, "TestApp", yamlCfg.App.Name)
	assert.Equal(t, "1.0", yamlCfg.App.Version)

	_ = os.Unsetenv("APP_NAME")
	_ = os.Remove(filepath.Join(dir, file))
}

func TestParseConfigFailEnvNotFound(t *testing.T) {
	content := `App:
  name: ${APP_NAME}
  version: 1.0
  log_level: 2
storage:
  master:
    name: MASTER_CONNECTION
    host: master-pg.localhost
    port: 5432
    user: user
    password: password
    database: db
  slave:
    name: SLAVE_CONNECTION
    host: slave-pg.localhost
    port: 5432
    user: user
    password: password
    database: db
`
	dir, file := createConfigFile(t, content)

	var yamlCfg AppConfig
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	assert.PanicsWithError(t, "environment variable not found: APP_NAME", func() {
		_ = config.ParseConfig(&yamlCfg, "App", dir)
	})

	_ = os.Remove(filepath.Join(dir, file))
}

func TestParseConfigFailUnsupportedFileExtension(t *testing.T) {
	dir := t.TempDir()
	unsupportedContent := `name: TestApp
	version: 1.0`
	err := os.WriteFile(filepath.Join(dir, "app.go"), []byte(unsupportedContent), 0644)
	assert.NoError(t, err)

	var yamlCfg AppConfig
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	err = config.ParseConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
}

func TestParseConfigFailNoDirFound(t *testing.T) {
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	var yamlCfg AppConfig
	err := config.ParseConfig(&yamlCfg, "app")
	assert.Error(t, err)
	assert.ErrorIs(t, err, goconfig.ErrOpenDir)
}

func TestParseConfigFailMultipleDirNotFound(t *testing.T) {
	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	var yamlCfg AppConfig
	err := config.ParseConfig(&yamlCfg, "app", "configuration")
	assert.Error(t, err)
	assert.ErrorIs(t, err, goconfig.ErrOpenDir)
}

func TestParseConfigFailUnmarshall(t *testing.T) {
	dir := t.TempDir()
	content := `
	name: TestApp
version: 1.0
`
	err := os.WriteFile(filepath.Join(dir, configFileYaml), []byte(content), 0644)
	assert.NoError(t, err)

	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	var yamlCfg AppConfig
	err = config.ParseConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
	assert.ErrorIs(t, err, goconfig.ErrUnmarshalling)
}

func TestParseConfigFailFileWithoutExtension(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "appconfig"), []byte("dummy content"), 0644)
	assert.NoError(t, err)

	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	var yamlCfg AppConfig
	err = config.ParseConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
	assert.ErrorIs(t, err, goconfig.ErrUnsupportedExt)
}

func TestParseConfigFailReadingFile(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, configFileYaml), []byte("dummy content"), 0000)
	assert.NoError(t, err)

	defer func(name string, mode os.FileMode) {
		_ = os.Chmod(name, mode)
	}(filepath.Join(dir, configFileYaml), 0644)

	config := goconfig.NewGoConfig()
	assert.NotNil(t, config)

	var yamlCfg AppConfig
	err = config.ParseConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
}

func createEnvFile(t *testing.T, content string) {
	dir := "."
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)
	assert.NoError(t, err)
}

func removeEnvFile(t *testing.T) {
	dir := "."
	err := os.Remove(filepath.Join(dir, ".env"))
	assert.NoError(t, err)
}

func createConfigFile(t *testing.T, content string) (string, string) {
	dir := t.TempDir()
	file := "App.yaml"
	err := os.WriteFile(filepath.Join(dir, file), []byte(content), 0644)
	assert.NoError(t, err)

	return dir, file
}

type AppConfig struct {
	App     App                `yaml:"App"`
	Storage map[string]Storage `yaml:"storage"`
}

type App struct {
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	LogLevel string `yaml:"log_level"`
}

type Storage struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
