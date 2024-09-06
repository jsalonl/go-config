# GoConfig

<div style="text-align:center">
	<a target="_blank" href="https://github.com"><img alt="github" src="https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white"/></a>
	<a target="_blank" href="https://go.dev/"><img alt="spring" src="https://img.shields.io/badge/Go-007d9c?style=for-the-badge&logo=go&logoColor=white"/></a>
</div>

`GoConfig` is a lightweight Go library for reading and unmarshalling configuration files in various formats, defaults in JSON or YAML.
It supports environment variable substitution within the configuration files and allows custom unmarshalling functions

## Features

- Supports YAML and JSON formats by default, allows customs unmarshall functions.
- Parses configuration files into user-defined Go structs.
- Allows configuration files to be stored in a specified directory or defaults to a "config" directory.
- Replaces environment variables in the configuration file with their actual values.
- Supports multiple configuration files with different names, file formats, and directories.
- Supports loading environment variables from one or more `.env` files.

## Installation

To install `GoConfig`, use `go get`:

```sh
go get github.com/jsalonl/go-config/v2
```

## Usage GoConfig

Suppose you have the following configuration file in YAML format:

```yaml
app:
  name: ${APP_NAME}
  version: ${VERSION}
  log_level: 2
storage:
  postgres:
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
  mongo:
    master:
      name: MASTER_CONNECTION
      host: master-mongo.localhost
      port: 27017
      user: user
      password: password
      database: db
```

You can create a Go struct that matches the structure of the configuration file:

```go
package config

type Config struct {
    App     App     `yaml:"app"`
    Storage Storage `yaml:"storage"`
}

type App struct {
    Name     string `yaml:"name"`
    Version  string `yaml:"version"`
    LogLevel string `yaml:"log_level"`
}

type Storage struct {
    Postgres map[string]Postgres `yaml:"postgres"`
    Mongo    map[string]Mongo    `yaml:"mongo"`
}

type Postgres struct {
    Name     string `yaml:"name"`
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
}

type Mongo struct {
    Name     string `yaml:"name"`
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
}
```

Also you have a `.env` file with the following content:

```env
APP_NAME=MyApp
VERSION=1.0.0
```

Here is an example of how to use `GoConfig`:

```go
package main

import (
    "fmt"

    "github.com/jsalonl/go-config/v2"
)

func main() {
    gonConf := goconfig.NewGoConfig()
    
    // Load environment variables
    err := gonConf.LoadEnv()
    if err != nil {
        panic(fmt.Errorf("error loading environment variables: %v", err))
    }

    var config AppConfig

    err = gonConf.ParseConfig(&appCfg, "app")
    if err != nil {
    panic(fmt.Errorf("error reading configuration: %v", err))
    }

    fmt.Printf("App Name: %s\n", appCfg.App.Name)       // From environment variable
    fmt.Printf("App Version: %s\n", appCfg.App.Version) // From environment variable
    fmt.Printf("Postgres master: %s\n", appCfg.Storage.Postgres["master"].Host) // Use key to get the value
    fmt.Printf("Postgres slave: %s\n", appCfg.Storage.Postgres["slave"].Host) // Use key to get the value
    fmt.Printf("Mongo master: %s\n", appCfg.Storage.Mongo["master"].Host) // Use key to get the value
}
```

### Use custom unmarshalling functions

You can use custom unmarshalling functions to parse configuration values into Go types that are not supported by the default unmarshalling functions.

Here is an example of how to use a custom unmarshalling function to parse a configuration toml file:

```go
package main

import (
    "fmt"

    "github.com/jsalonl/go-config/v2"
)

func main() {
    gonConf := goconfig.NewGoConfig(unmarshallTOML)
    
    // Load environment variables
    err := gonConf.LoadEnv()
    if err != nil {
        panic(fmt.Errorf("error loading environment variables: %v", err))
    }

    var config AppConfig

    err = gonConf.ParseConfig(&appCfg, "app")
    if err != nil {
        panic(fmt.Errorf("error reading configuration: %v", err))
    }

    fmt.Printf("App Name: %s\n", appCfg.App.Name)       // From environment variable
    fmt.Printf("App Version: %s\n", appCfg.App.Version) // From environment variable
    fmt.Printf("Postgres master: %s\n", appCfg.Storage.Postgres["master"].Host) // Use key to get the value
    fmt.Printf("Postgres slave: %s\n", appCfg.Storage.Postgres["slave"].Host) // Use key to get the value
    fmt.Printf("Mongo master: %s\n", appCfg.Storage.Mongo["master"].Host) // Use key to get the value
}

func unmarshallTOML(structure interface{}, content []byte) error {
    err := toml.Unmarshal(content, structure)
    if err != nil {
        return fmt.Errorf("error unmarshalling configuration: %v", err)
    }

    return nil
}
```

### Environment Variables

You can use the method `LoadEnv` to load environment variables from one or more `.env` files.

Here is an example of a configuration file with environment variables:

```yaml
name: ${APP_NAME}
version: ${APP_VERSION}
```

## Usage LoadEnv

Here is an example of how to use `GoConfig`:

```go
package main

import (
    "fmt"
    "os"

    "github.com/jsalonl/go-config/v2"
)

func main() {
    gonConf := goconfig.NewGoConfig()

    // Load environment variables
    err := gonConf.LoadEnv()
    if err != nil {
        panic(fmt.Errorf("error loading environment variables: %v", err))
    }

    fmt.Printf("App Name: %s\n", os.Getenv("APP_NAME"))
    fmt.Printf("App Version: %s\n", os.Getenv("APP_VERSION"))
}
```

You can use multiple `.env` files by specifying the file names as arguments to the `LoadEnv` function.

```go
import (
    "fmt"
    "os"

    "github.com/jsalonl/go-config/v2"
)

func main() {
    gonConf := goconfig.NewGoConfig()

    // Load environment variables
    err := gonConf.LoadEnv("database.env", "mail.env")
    if err != nil {
        panic(fmt.Errorf("error loading environment variables: %v", err))
    }

    fmt.Printf("Database Host: %s\n", os.Getenv("DB_HOST"))
    fmt.Printf("Mail Server: %s\n", os.Getenv("MAIL_SERVER"))
}
```

## Sonar report

![Sonar report](https://i.imghippo.com/files/J9Mnn1724798103.png)

## License

This project is licensed under the MIT License.
This means you are free to use, modify, and distribute the software as you wish. See the [LICENSE](https://www.mit.edu/~amini/LICENSE.md) file for details.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request.

## Do you want to support me?

<a href="https://www.buymeacoffee.com/JoanSalomon" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-red.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>