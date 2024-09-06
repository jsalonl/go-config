# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v2.0.0] - 2024-09-06

### Changed

- Refactored `GoConfig` to use an interface and `goConfig` implementation.
- Added `unmarshallFunc` parameter to `NewGoConfig` to allow custom unmarshalling functions.
- Renamed `NewConfig` to `ParseConfig` and `readConfig` to `read` for clarity.
- Updated error handling to use `fmt.Errorf` with the new `formatError` constant.
- Improved regex patterns for environment variables and `.env` file parsing.
- Simplified environment variable replacement logic in `replaceEnvVariables`.
- Updated `unmarshallYAML` to handle errors using the new error format.

### Removed

- Deprecated `ErrVariableNotFound`, `ErrReadingFile`, `ErrUnmarshalling`, `ErrUnsupportedExt`, `ErrReadingConfig`,
- `ErrOpeningEnvFile`, and `ErrInvalidEnvFormat` constants in favor of a consolidated error handling approach.

### Fixed

- Corrected the error message formatting in `openFile` and `read`.
- Improved the handling of unsupported file extensions in the `read` function.
- Enhanced the `.env` file parsing to handle invalid formats more gracefully.

### Other Changes

- **Error Handling:** All error constants are now defined using the `errors` package. This change consolidates error
  messages and simplifies error management across the package. The following errors are now managed via `errors.New`:
  - `ErrUnmarshalling`
  - `ErrVariableNotFound`
  - `ErrUnsupportedExt`
  - `ErrReadingFile`
  - `ErrOpenDir`
  - `ErrOpeningEnvFile`
  - `ErrInvalidEnvFormat`

## [v1.0.0] - 2024-08-27

### Added

- Initial release of `GoConfig`.
- Support for reading YAML and JSON configuration files.
- Environment variable substitution in configuration files.
- Support .env file for environment variables.
- SonarQube integration and code quality improvements.