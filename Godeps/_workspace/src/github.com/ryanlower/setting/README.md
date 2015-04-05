## setting

A simple utility to load your configuration from the environment (with optional defaults)

[![Circle CI](https://circleci.com/gh/ryanlower/setting.svg?style=svg)](https://circleci.com/gh/ryanlower/setting)

```go
type Config struct {
  Port string `env:"PORT" default:"3000"`
}

config := new(Config)
setting.Load(config)

// If the PORT environment variable is set, `config.Port` will
// now be its value (`os.Getenv("PORT")`). Otherwise it will
// be "3000"
```

Both `env` and `default` are optional:
```go
type Config struct {
  Name string `env:"NAME"`
  Location struct {
    City string `default:"San Francisco"`
  }
}

config := new(Config)
setting.Load(config)

// `config.Name` will be the NAME environment variable value if
// present, or the zero value if not
// `config.Location.City` will be "San Francisco"
```

### Limitations
Setting currently only supports string type fields! If you need another type open an issue or a pull request :smile:
