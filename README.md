[![GoDoc](https://godoc.org/github.com/golobby/dotenv?status.svg)](https://godoc.org/github.com/golobby/dotenv)
[![CI](https://github.com/golobby/dotenv/actions/workflows/ci.yml/badge.svg)](https://github.com/golobby/dotenv/actions/workflows/ci.yml)
[![CodeQL](https://github.com/golobby/dotenv/workflows/CodeQL/badge.svg)](https://github.com/golobby/dotenv/actions?query=workflow%3ACodeQL)
[![Go Report Card](https://goreportcard.com/badge/github.com/golobby/dotenv)](https://goreportcard.com/report/github.com/golobby/dotenv)
[![Coverage Status](https://coveralls.io/repos/github/golobby/dotenv/badge.svg?v=1)](https://coveralls.io/github/golobby/dotenv)

# DotEnv
GoLobby DotEnv is a lightweight package for loading dot env (.env) files into structs for Go projects

## Documentation
### Supported Versions
It requires Go `v1.16` or newer versions.

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

IPS=192.168.0.1,192.168.0.2
IDS=10,11,12,13

DB_NAME=shop
DB_PORT=3306
DB_USER=root
DB_PASS=secret
```

Sample `.go` file:

```go
type Config struct {
    Debug bool      `env:"DEBUG"`
    App struct {
        Name string `env:"APP_NAME"`
        Port int16  `env:"APP_PORT"`
    }
    Database struct {
        Name string `env:"DB_NAME"`
        Port int16  `env:"DB_PORT"`
        User string `env:"DB_USER"`
        Pass string `env:"DB_PASS"`
    }
    IPs []string `env:"IPS"`
	IDs []int64  `env:"IDS"`
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

### DotEnv Syntax
The following snippet shows a valid dot env file.

```env
String = Hello Dot Env # Comment

# Quotes
Quote1="Quoted message!"
Quote2="You can use ' here"
Quote3='You can use " here'
Quote4="You can use # here"

# Booleans
Bool1 = true
Bool2 = 1 # true
Bool3 = false
Bool4 = 0 # false

# Arrays
Ints    = 1,2, 3, 4 , 5 # []int{1, 2, 3, 4, 5}
Strings = a,b, c, d , e # []string{"a", "b", "c", "d", "e"}
Floats  = 3.14,9.8, 6.9 # []float32{3.14, 9.8, 6.9}

```

## See Also
* [GoLobby/Config](https://github.com/golobby/config):
  A lightweight yet powerful configuration management for Go projects
* [GoLobby/Env](https://github.com/golobby/env):
  A lightweight package for loading OS environment variables into structs for Go projects

## License
GoLobby DotEnv is released under the [MIT License](http://opensource.org/licenses/mit-license.php).
