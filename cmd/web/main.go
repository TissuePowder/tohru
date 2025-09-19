package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"tohru/internal/env"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

type config struct {
	baseUrl  string
	httpPort string
}

type application struct {
	cfg    config
	logger *slog.Logger
	wg     sync.WaitGroup
}

func main() {

	logLvl := env.GetStr("LOG_LEVEL", "info")

	var lvl slog.Level
	err := (&lvl).UnmarshalText([]byte(logLvl))
	if err != nil {
		fmt.Printf("invalid LOG_LEVEL %q (using info)\n", logLvl)
		lvl = slog.LevelInfo
	}

	lvlvar := new(slog.LevelVar)
	lvlvar.Set(lvl)

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lvlvar,
		AddSource: lvl <= slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "source" {
				if src, ok := a.Value.Any().(*slog.Source); ok {
					return slog.String(a.Key, fmt.Sprintf("%s:%d", src.File, src.Line))
				}
			}
			return a
		},
	}))

	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(lvl)

	cfg := config{
		baseUrl:  env.GetStr("BASE_URL", "127.0.0.1"),
		httpPort: env.GetStr("HTTP_PORT", "8080"),
	}

	app := &application{
		cfg:    cfg,
		logger: logger,
	}

	err = app.serveHTTP()
	if err != nil {
		logger.Error(err.Error())
		debug.PrintStack()
	}

}
