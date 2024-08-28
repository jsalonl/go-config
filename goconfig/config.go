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
	supportedExtensions = []string{"yaml", "yml", "json"}
	setEnvFunc          = os.Setenv
)

const (
	// ErrVariableNotFound is the error message for a missing environment variable.
	ErrVariableNotFound = "environment variable %s not found"
	// ErrReadingFile is the error message for a file reading error.
	ErrReadingFile = "error reading file: %w"
	// ErrUnmarshalling is the error message for an unmarshalling error.
	ErrUnmarshalling = "error unmarshalling configuration: %v"
	// ErrUnsupportedExt is the error message for an unsupported extension.
	ErrUnsupportedExt = "unsupported extension %s"
	// ErrReadingConfig is the error message for a configuration reading error.
	ErrReadingConfig = "error reading configuration: %w"
	// ErrOpeningEnvFile is the error message for an error opening a .env file.
	ErrOpeningEnvFile = "error opening .env file: %w"
	// ErrInvalidEnvFormat is the error message for an invalid .env format.
	ErrInvalidEnvFormat = "invalid .env format on line: %s"
	// envVarPattern is the regex pattern to match lines in a .env file.
	envVarPattern = `^\s*([\w.-]+)\s*=\s*(.*)?\s*$`
)

// NewConfig reads a configuration file from a directory and unmarshalls it into a structure.
// If no directory is provided, it will use the default directory "config".
func NewConfig(structure interface{}, configName string, directoryName ...string) error {
	content, err := readConfig(configName, directoryName...)
	if err != nil {
		return err
	}

	if err := unmarshallGeneric(structure, content); err != nil {
		return err
	}

	return nil
}

// readConfig reads a configuration file from a directory.
// If no directory is provided, it will use the default directory "config".
func readConfig(configName string, directoryName ...string) ([]byte, error) {
	dir := "config"
	if len(directoryName) > 0 {
		dir = directoryName[0]
	}
	content, err := read(configName, dir)
	if err != nil {
		return nil, fmt.Errorf(ErrReadingConfig, err)
	}

	return content, nil
}

// unmarshallGeneric unmarshalls the content into the structure.
// Supported formats are YAML and JSON (JSON is a subset of YAML).
func unmarshallGeneric(structure interface{}, content []byte) error {
	err := yaml.Unmarshal(content, structure)
	if err != nil {
		return fmt.Errorf(ErrUnmarshalling, err)
	}

	return nil
}

// read reads a file from a directory and returns its content and extension.
// If no file is found, it returns an error.
func read(profile string, basePath string) ([]byte, error) {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		name, extension, found := strings.Cut(file.Name(), ".")
		if !found {
			continue
		}

		if !slices.Contains(supportedExtensions, extension) {
			continue
		}

		if strings.EqualFold(name, profile) {
			content, err := os.ReadFile(path.Join(basePath, file.Name()))
			if err != nil {
				return nil, fmt.Errorf(ErrReadingFile, err)
			}

			contentStr, _ := replaceEnvVariables(string(content))

			return []byte(contentStr), nil
		}
	}

	return nil, fmt.Errorf(ErrUnsupportedExt, profile)
}

// replaceEnvVariables replaces the environment variables in the content using the format ${ENV_VAR}.
// If the environment variable is not found, it will panic returning the name of the variable.
func replaceEnvVariables(content string) (string, error) {
	re := regexp.MustCompile(`\${(\w+)}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		envVar := re.FindStringSubmatch(match)[1]
		env := os.Getenv(envVar)
		if env == "" {
			panic(fmt.Errorf(ErrVariableNotFound, envVar))
		}

		return env
	}), nil
}

// LoadEnv loads the environment variables from one or more .env files.
// If no directory is provided, it will look in the current working directory.
func LoadEnv(envFiles ...string) error {
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

// openFile abstracts the logic of opening a file and returning a file handle.
func openFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf(ErrOpeningEnvFile, err)
	}

	return file, nil
}

// parseEnvFile reads and parses the .env file, setting the environment variables.
func parseEnvFile(scanner *bufio.Scanner) error {
	re := regexp.MustCompile(envVarPattern)

	for scanner.Scan() {
		line := scanner.Text()
		if isCommentOrEmpty(line) {
			continue
		}

		if err := setEnvVarFromLine(re, line); err != nil {
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
	if !re.MatchString(line) {
		return fmt.Errorf(ErrInvalidEnvFormat, line)
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf(ErrInvalidEnvFormat, line)
	}

	key, value := parts[0], parts[1]

	return setEnvFunc(key, value)
}
