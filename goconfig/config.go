package goconfig

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	excludeExtensions = []string{"go"}
	regexEnv          = regexp.MustCompile(`\${(\w+)}`)
	regexEnvFromFile  = regexp.MustCompile(`^\s*([\w.-]+)\s*=\s*(.*)?\s*$`)
)

const (
	formatError = "%w: %v"
)

// goConfig is the GoConfig implementation.
type goConfig struct {
	unmarshallFunc func(interface{}, []byte) error
}

// GoConfig is the interface that wraps the Read, LoadEnv and Unmarshall methods.
type GoConfig interface {
	// LoadEnv loads environment variables from a .env files.
	// If no files are provided, it will use the default file ".env".
	LoadEnv(envFiles ...string) error
	// ParseConfig reads a configuration file from a directory and unmarshalls it into a structure.
	// If no directory is provided, it will use the default directory "config".
	ParseConfig(structure interface{}, fileName string, directoryName ...string) error
}

// NewGoConfig creates a new GoConfig instance.
// It receives an optional unmarshalling function, if not provided it will default to unmarshallYAML.
func NewGoConfig(unmarshallFunc ...func(interface{}, []byte) error) GoConfig {
	var unmarshall func(interface{}, []byte) error
	if len(unmarshallFunc) > 0 {
		unmarshall = unmarshallFunc[0]
	} else {
		unmarshall = unmarshallYAML
	}

	return &goConfig{unmarshallFunc: unmarshall}
}

func (g goConfig) LoadEnv(envFiles ...string) error {
	dir := "."
	if len(envFiles) == 0 {
		envFiles = []string{".env"}
	}

	for _, envFile := range envFiles {
		filePath := path.Join(dir, envFile)
		file, err := openFile(filePath)
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(file)
		if err := parseEnvFile(scanner); err != nil {
			return err
		}

		_ = file.Close()
	}

	return nil
}

func (g goConfig) ParseConfig(structure interface{}, configName string, directoryName ...string) error {
	content, err := read(configName, directoryName...)
	if err != nil {
		return err
	}

	if err := g.unmarshallFunc(structure, content); err != nil {
		return err
	}

	return nil
}

// read reads a file from a directory and returns its content and extension.
// If no file is found, it returns an error.
func read(fileName string, basePath ...string) ([]byte, error) {
	dir := "config"
	if len(basePath) > 0 {
		dir = basePath[0]
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf(formatError, ErrOpenDir, basePath)
	}

	for _, file := range files {
		name, extension, found := strings.Cut(file.Name(), ".")
		if !found {
			continue
		}

		if slices.Contains(excludeExtensions, extension) {
			continue
		}

		if strings.EqualFold(name, fileName) {
			content, err := os.ReadFile(path.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf(formatError, ErrReadingFile, fileName)
			}

			contentStr := replaceEnvVariables(string(content))

			return []byte(contentStr), nil
		}
	}

	return nil, fmt.Errorf("%w: in profile %v", ErrUnsupportedExt, fileName)
}

// replaceEnvVariables replaces the environment variables in the content using the format ${ENV_VAR}.
// If the environment variable is not found, it will panic returning the name of the variable.
func replaceEnvVariables(content string) string {
	return regexEnv.ReplaceAllStringFunc(content, func(match string) string {
		envVar := regexEnv.FindStringSubmatch(match)[1]
		env := os.Getenv(envVar)
		if env == "" {
			panic(fmt.Errorf(formatError, ErrVariableNotFound, envVar))
		}

		return env
	})
}

// unmarshallYAML unmarshalls the content into the structure.
// Supported formats are YAML and JSON (JSON is a subset of YAML).
func unmarshallYAML(structure interface{}, content []byte) error {
	err := yaml.Unmarshal(content, structure)
	if err != nil {
		return fmt.Errorf(formatError, ErrUnmarshalling, err)
	}

	return nil
}

// openFile abstracts the logic of opening a file and returning a file handle.
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: in %v", ErrOpeningEnvFile, filePath)
	}

	return file, nil
}

// parseEnvFile reads and parses the .env file, setting the environment variables.
func parseEnvFile(scanner *bufio.Scanner) error {
	for scanner.Scan() {
		line := scanner.Text()
		if isCommentOrEmpty(line) {
			continue
		}

		if err := setEnvVarFromLine(regexEnvFromFile, line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading .env file: %w", err)
	}

	return nil
}

// isCommentOrEmpty checks if a line is a comment or empty.
func isCommentOrEmpty(line string) bool {
	return strings.HasPrefix(line, "#") || strings.TrimSpace(line) == ""
}

// setEnvVarFromLine parses a line and sets the corresponding environment variable.
func setEnvVarFromLine(re *regexp.Regexp, line string) error {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf(formatError, ErrInvalidEnvFormat, line)
	}

	if !re.MatchString(line) {
		return fmt.Errorf(formatError, ErrInvalidEnvFormat, line)
	}

	key, value := parts[0], parts[1]

	return os.Setenv(key, value)
}
