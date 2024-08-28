# GoConfig

<div align="center">
	<a target="_blank" href="https://github.com"><img alt="github" src="https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white"/></a>
	<a target="_blank" href="https://go.dev/"><img alt="spring" src="https://img.shields.io/badge/Go-007d9c?style=for-the-badge&logo=go&logoColor=white"/></a>
</div>


`GoConfig` is a lightweight Go library for reading and unmarshalling configuration files in various formats, 
including YAML and JSON. It supports environment variable substitution within the configuration files.

## Features

- Supports YAML and JSON formats.
- Parses configuration files into user-defined Go structs.
- Allows configuration files to be stored in a specified directory or defaults to a "config" directory.
- Replaces environment variables in the configuration file with their actual values.

## Installation

To install `GoConfig`, use `go get`:

```sh
go get github.com/jsalonl/go_config
```

## Usage GoConfig

Here is an example of how to use `GoConfig`:

### Go File

```go
package main

import (
    "fmt"

    "github.com/jsalonl/go-config"
)

type AppConfig struct {
    Name    string `yaml:"name"`
    Version string `yaml:"version"`
}

func main() {
    var config AppConfig

    err := goconfig.NewConfig(&config, "app")
    if err != nil {
        panic(err)
    }

    fmt.Printf("App Name: %s\n", config.Name)
    fmt.Printf("App Version: %s\n", config.Version)
}
```

### Configuration File

The configuration file should be named `app.yaml` or `app.json` and stored in the `config` directory.

Here is an example of a configuration file in YAML format:

```yaml
name: MyApp
version: 1.0.0
```

### Environment Variables

You can use environment variables in the configuration file by enclosing them in `${}`.

Here is an example of a configuration file with environment variables:

```yaml
name: ${APP_NAME}
version: ${APP_VERSION}
```

## Usage LoadEnv

Here is an example of how to use `GoConfig`:

### Go File

```go
package main

import (
    "fmt"
	"os"

    "github.com/jsalonl/go-config"
)

func main() {
    err := goconfig.LoadEnv()
	if err != nil {
        panic(err)
    }

    fmt.Printf("App Name: %s\n", os.Getenv("APP_NAME"))
    fmt.Printf("App Version: %s\n", os.Getenv("APP_VERSION"))
}
```

### .env File

The `.env` file should be stored in the root directory of the project.

Here is an example of a configuration file in YAML format:

```env
APP_NAME=MyApp
APP_VERSION=1.0.0
```

You can use multiple `.env` files by specifying the file names as arguments to the `LoadEnv` function.

```go
package main

import (
    "fmt"
	"os"

    "github.com/jsalonl/go-config"
)

func main() {
    err := goconfig.LoadEnv("database.env", "mail.env")
	if err != nil {
        panic(err)
    }

    fmt.Printf("Database Host: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("Mail Server: %s\n", os.Getenv("MAIL_SERVER"))
}
```

Alternative, you can use the `LoadEnv` with `NewConfig` function to load the environment variables and configuration file.

```go
package main

import (
    "fmt"

    "github.com/jsalonl/go-config"
)

type AppConfig struct {
    Name    string `yaml:"name"`
    Version string `yaml:"version"`
}

func main() {
	err := goconfig.LoadEnv()
	if err != nil {
		panic(err)
	}
    
	var config AppConfig
    err = goconfig.NewConfig(&config, "app")
    if err != nil {
        panic(err)
    }

    fmt.Printf("App Name: %s\n", config.Name)
    fmt.Printf("App Version: %s\n", config.Version)
}
```

## Sonar report

![Sonar report](https://i.imghippo.com/files/J9Mnn1724798103.png)

## License

This project is licensed under the MIT License. 
This means you are free to use, modify, and distribute the software as you wish. See the [LICENSE](https://www.mit.edu/~amini/LICENSE.md) file for details.

