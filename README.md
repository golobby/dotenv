[![GoDoc](https://godoc.org/github.com/golobby/dotenv?status.svg)](https://godoc.org/github.com/golobby/dotenv)
[![CI](https://github.com/golobby/dotenv/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/dotenv/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/dotenv/workflows/CodeQL/badge.svg)](https://github.com/golobby/dotenv/actions?query=workflow%3ACodeQL)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/dotenv)](https://goreportcard.com/report/github.com/golobby/dotenv)
[![Coverage Status](https://coveralls.io/repos/github/golobby/dotenv/badge.svg)](https://coveralls.io/github/golobby/dotenv)

# DotEnv
GoLobby DotEnv is a lightweight package for loading OS environment variables into structs for Go projects.

## Documentation
### Supported Versions
It requires Go `v1.11` or newer versions.

### Installation
To install this package run the following command in the root of your project.

```bash
go get github.com/golobby/dotenv
```

### Usage Example

Sample `.env` file:

```env
DEBUG=true

APP_NAME=MyApp
APP_PORT=8585

DB_NAME=shop
DB_PORT=3306
DB_USER=root
DB_PASS=secret
```

Sample `.go` file:

```go
type Config struct {
    Debug bool      `dotenv:"DEBUG"`
    App struct {
        Name string `dotenv:"APP_NAME"`
        Port int16  `dotenv:"APP_PORT"`
    }
    Database struct {
        Name string `dotenv:"DB_NAME"`
        Port int16  `dotenv:"DB_PORT"`
        User string `dotenv:"DB_USER"`
        Pass string `dotenv:"DB_PASS"`
    }
}

config := Config{}
file, err := os.Open(".env")

err = dotenv.NewDecoder(file).Decode(&config)

// Use `config` struct in your app!
```

### Usage Tips
* The `Decode()` function gets a pointer of a struct.
* It ignores the fields that have no related environment variables in the file.
* It supports nested structs and struct pointers.

### Field Types
GoLobby DotEnv uses the [GoLobby Cast](https://github.com/golobby/cast) package to cast environment variables to related struct field types.
Here you can see the supported types:

https://github.com/golobby/cast#supported-types

## See Also
* [GoLobby/Config](https://github.com/golobby/config):
  A lightweight yet powerful config package for Go projects
* [GoLobby/Env](https://github.com/golobby/env):
  A lightweight package for loading OS environment variables into structs for Go projects

## License
GoLobby DotEnv is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
