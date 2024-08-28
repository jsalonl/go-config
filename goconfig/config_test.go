package goconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	configFileYaml = "app.yaml"
	configFileJson = "app.json"
)

func TestNewConfigSuccessYAML(t *testing.T) {
	dir := t.TempDir()
	content := `name: TestApp
version: 1.0`
	err := os.WriteFile(filepath.Join(dir, configFileYaml), []byte(content), 0644)
	assert.NoError(t, err)
	var yamlCfg struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	err = NewConfig(&yamlCfg, "app", dir)
	assert.NoError(t, err)
	assert.Equal(t, "TestApp", yamlCfg.Name)
	assert.Equal(t, "1.0", yamlCfg.Version)
}

func TestNewConfigSuccessMultipleDir(t *testing.T) {
	dir := t.TempDir()
	content := `name: TestApp
version: 1.0`
	err := os.WriteFile(filepath.Join(dir, configFileYaml), []byte(content), 0644)
	assert.NoError(t, err)
	var yamlCfg struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	err = NewConfig(&yamlCfg, "app", dir, "configuration")
	assert.NoError(t, err)
	assert.Equal(t, "TestApp", yamlCfg.Name)
	assert.Equal(t, "1.0", yamlCfg.Version)
}

func TestNewConfigSuccessJSON(t *testing.T) {
	dir := t.TempDir()
	jsonContent := `{
	"name": "TestApp",
	"version": "1.0"
}`
	err := os.WriteFile(filepath.Join(dir, configFileJson), []byte(jsonContent), 0644)
	assert.NoError(t, err)
	var yamlCfg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	err = NewConfig(&yamlCfg, "app", dir)
	assert.NoError(t, err)
	assert.Equal(t, "TestApp", yamlCfg.Name)
	assert.Equal(t, "1.0", yamlCfg.Version)
}

func TestNewConfigSuccessReplaceEnvVariables(t *testing.T) {
	err := os.Setenv("TESTAPPNAME", "goConfig")
	require.NoError(t, err)
	defer func() {
		_ = os.Unsetenv("TESTAPPNAME")
	}()

	content := `name: ${TESTAPPNAME}
version: 1.0`
	replacedContent, err := replaceEnvVariables(content)
	assert.NoError(t, err)
	assert.Equal(t, "name: goConfig\nversion: 1.0", replacedContent)
}

func TestNewConfigErrorUnsupportedFileExtension(t *testing.T) {
	dir := t.TempDir()
	unsupportedContent := `name: TestApp
	version: 1.0`
	err := os.WriteFile(filepath.Join(dir, "app.toml"), []byte(unsupportedContent), 0644)
	assert.NoError(t, err)
	var yamlCfg struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	err = NewConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
}

func TestNewConfigErrorNoDirFound(t *testing.T) {
	var yamlCfg struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	err := NewConfig(&yamlCfg, "app")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading configuration: open config: no such file or directory")
}

func TestNewConfigErrorMultipleDirNotFound(t *testing.T) {
	var yamlCfg struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	err := NewConfig(&yamlCfg, "app", "configuration")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading configuration: "+
		"open configuration: no such file or directory")
}

func TestNewConfigErrorUnmarshallYAML(t *testing.T) {
	dir := t.TempDir()
	content := `
	name: TestApp
version: 1.0
`
	err := os.WriteFile(filepath.Join(dir, configFileYaml), []byte(content), 0644)
	assert.NoError(t, err)
	var yamlCfg struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	err = NewConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: line 2: found character that cannot start any token")
}

func TestNewConfigErrorUnmarshallJSON(t *testing.T) {
	dir := t.TempDir()
	jsonContent := `
{
	"name": "TestApp",}
	"version": "1.0"
}
`
	err := os.WriteFile(filepath.Join(dir, "app.json"), []byte(jsonContent), 0644)
	assert.NoError(t, err)
	var yamlCfg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	err = NewConfig(&yamlCfg, "app", dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yaml: line 4: found character that cannot start any token")
}

func TestNewConfigErrorFileWithoutExtension(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "appconfig"), []byte("dummy content"), 0644)
	assert.NoError(t, err)
	content, err := read("app", dir)
	assert.Nil(t, content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported extension app")
}

func TestNewConfigErrorReadingFile(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, configFileYaml), []byte("dummy content"), 0000)
	assert.NoError(t, err)

	defer func(name string, mode os.FileMode) {
		_ = os.Chmod(name, mode)
	}(filepath.Join(dir, configFileYaml), 0644)

	content, err := read("app", dir)
	assert.Nil(t, content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading file")
}

func TestNewConfigErrorReplaceEnvVariablesPanicOnMissingEnvVar(t *testing.T) {
	content := `name: ${TESTAPPNAME}
version: 1.0`
	assert.PanicsWithError(t, "environment variable TESTAPPNAME not found", func() {
		_, _ = replaceEnvVariables(content)
	})
}

func TestLoadEnvSuccess(t *testing.T) {
	dir := "."
	content := `APP_NAME=TestApp
APP_VERSION=1.0
`
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)
	assert.NoError(t, err)

	err = LoadEnv()
	assert.NoError(t, err)

	assert.Equal(t, "TestApp", os.Getenv("APP_NAME"))
	assert.Equal(t, "1.0", os.Getenv("APP_VERSION"))

	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
	err = os.Remove(filepath.Join(dir, ".env"))
	assert.NoError(t, err)
}

func TestLoadEnvSuccessWithComment(t *testing.T) {
	dir := "."
	content := `APP_NAME=TestApp
# This is a comment
APP_VERSION=1.0
`
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)
	assert.NoError(t, err)

	err = LoadEnv()
	assert.NoError(t, err)

	assert.Equal(t, "TestApp", os.Getenv("APP_NAME"))
	assert.Equal(t, "1.0", os.Getenv("APP_VERSION"))

	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
	err = os.Remove(filepath.Join(dir, ".env"))
	assert.NoError(t, err)
}

func TestLoadEnvSuccessMultipleFiles(t *testing.T) {
	dir := "."
	content := `APP_NAME=TestApp
APP_VERSION=1.0
`
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)
	assert.NoError(t, err)

	content = `AUTHOR=John Doe
`

	err = os.WriteFile(filepath.Join(dir, ".env2"), []byte(content), 0644)
	assert.NoError(t, err)

	err = LoadEnv(".env", ".env2")
	assert.NoError(t, err)

	assert.Equal(t, "TestApp", os.Getenv("APP_NAME"))
	assert.Equal(t, "1.0", os.Getenv("APP_VERSION"))
	assert.Equal(t, "John Doe", os.Getenv("AUTHOR"))

	_ = os.Unsetenv("APP_NAME")
	_ = os.Unsetenv("APP_VERSION")
	_ = os.Unsetenv("AUTHOR")
	err = os.Remove(filepath.Join(dir, ".env"))
	assert.NoError(t, err)
	err = os.Remove(filepath.Join(dir, ".env2"))
	assert.NoError(t, err)
}

func TestLoadEnvErrorOpenFile(t *testing.T) {
	err := LoadEnv("nonexistentfile")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open nonexistentfile: no such file or directory")
}

func TestLoadEnvErrorInvalidFormat(t *testing.T) {
	dir := "."
	content := `APP_NAME:TestApp
APP_VERSION:1.0
`
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)
	assert.NoError(t, err)

	err = LoadEnv()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid .env format on line: APP_NAME:TestApp")

	err = os.Remove(filepath.Join(dir, ".env"))
	assert.NoError(t, err)
}

func TestSetEnvVarError(t *testing.T) {
	originalSetenv := setEnvFunc
	defer func() {
		setEnvFunc = originalSetenv
	}()

	setEnvFunc = func(key, value string) error {
		return fmt.Errorf("simulated error setting %s", key)
	}

	content := `APP_NAME=TestApp`
	dir := "."
	err := os.WriteFile(filepath.Join(dir, ".env"), []byte(content), 0644)
	require.NoError(t, err)

	err = LoadEnv()
	require.Error(t, err)

	err = os.Remove(filepath.Join(dir, ".env"))
	assert.NoError(t, err)
}
