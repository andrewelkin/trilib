package utils

import (
	"context"
	"github.com/andrewelkin/trilib/utils/logger"
	"github.com/nats-io/nats.go"
	"os"
	"strings"
	"sync"
)

type ContextWithCancel struct {
	context.Context
	Cancel context.CancelFunc
	Logger logger.Logger
	Wg     sync.WaitGroup // wait group
}

var globalContext *ContextWithCancel
var GIL sync.Mutex

/*
delete me later -- unused AE 3/17/2023

func SetGlobalLogger(logger logger.Logger) logger.Logger {
	if globalContext != nil {
		globalContext.Logger = logger
	} else {
		logger.Warnf("*", "Unable to set global logger")
	}
	return logger
}
*/

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func GetGlobalContext() *ContextWithCancel {
	return globalContext
}

func loggerFromConfig(cfg IConfig, defaultFilter logger.FilterFunc) logger.FilterFunc {
	filter := logger.And(
		logger.FilterOrDefault(cfg.GetString("filter"), defaultFilter),
		logger.Not(logger.FilterOrDefault(cfg.GetString("exclude"), logger.FilterMatchNone)),
	)
	return filter
}

// GetOrCreateGlobalContext sets a new global context with logging and cancel
// Expects a config, which is normally would be a "Logging" section
func GetOrCreateGlobalContext(gconfig IConfig) *ContextWithCancel {
	GIL.Lock()
	defer GIL.Unlock()

	if globalContext != nil {
		return globalContext
	}
	ctx, cancel := context.WithCancel(context.Background())
	var config IConfig

	if gconfig != nil {
		config = gconfig.FromKey("logger")
	}

	var logLevel logger.LogLevel
	logLevel = logger.LogLevelDebug
	logNameSpace := "*"

	var outputsCfg IConfig
	if config != nil {
		if tmp := config.GetString("loglevel"); tmp != nil {
			logLevel = logger.ParseLogLevel(*tmp, logger.LogLevelDebug)
		}

		// allow the default stdout log namespace filter to be overridden by the "filter" config field
		logger.LogDefaultFilter = loggerFromConfig(config, logger.FilterUnderscore)
		outputsCfg = config.FromKey("outputs")
	}

	// create the strategies global logger (default writes to STDOUT)
	globalLogger := logger.GetOrCreateGlobalLoggerEx(ctx, logLevel, logger.FilterUnderscore)

	// add additional writers, if configured
	if outputsCfg != nil {
		allOutputCfg := outputsCfg.GetCfg()
		for outputType := range allOutputCfg {

			cfg := outputsCfg.FromKey(outputType)
			if cfg == nil {
				panic("failed to parse output config: " + outputType)
			}

			switch strings.ToLower(outputType) {

			case "nats_publisher", "natspublisher", "nats":
				subject := cfg.GetStringDefault("subject", "default-logger-subject")
				url := cfg.GetStringDefault("url", nats.DefaultURL)
				filter := logger.And(
					logger.FilterOrDefault(cfg.GetString("filter"), logger.FilterMatchAll),
					logger.Not(logger.FilterOrDefault(cfg.GetString("exclude"), logger.FilterMatchNone)),
				)

				rawLevel := cfg.GetStringDefault("logLevel", "debug")
				ansi := cfg.GetBoolDefault("ansicodes", false)

				if len(*subject) == 0 {
					panic("failed to create nats logger : empty publishing subject")
				}

				var nkeyOpt nats.Option
				nSeedFile := *cfg.GetStringDefault("natsSeed", "")
				if len(nSeedFile) > 0 {
					var err error
					nkeyOpt, err = nats.NkeyOptionFromSeed(nSeedFile)
					if err != nil {
						panic(err)
					}
				}

				nc, err := nats.Connect(*url, nkeyOpt)
				if err != nil {
					panic(err)
				}
				globalLogger.AddOutput(filter, logger.NewNatsLogger(*subject, nc), logger.ParseLogLevel(*rawLevel, logger.LogLevelDebug), ansi, false)

			case "filewriter", "file":
				rawLevel, path, prefix, suffix, filter, skipRepeating :=
					cfg.GetStringDefault("logLevel", "debug"),
					cfg.GetStringDefault("path", "/tmp/test_logs"),
					cfg.GetStringDefault("filePrefix", ""),
					cfg.GetStringDefault("fileSuffix", ".log"),
					loggerFromConfig(cfg, logger.FilterMatchAll),
					cfg.GetBoolDefault("skipRepeating", true)

				fileWriter, err := logger.NewFileWriter(*path, prefix, suffix, skipRepeating)
				if err != nil {
					panic("failed to create file writer: " + err.Error())
				}

				globalLogger.Infof(logNameSpace, "Adding log file output; filter=%s path=%s level=%s", filter, *path, *rawLevel)
				globalLogger.AddOutput(
					filter,
					fileWriter,
					logger.ParseLogLevel(*rawLevel, logger.LogLevelDebug), false, true)

			case "jsonstream", "jsonout", "prod":
				rawLevel, filter :=
					cfg.GetStringDefault("logLevel", "info"),
					loggerFromConfig(cfg, logger.FilterMatchAll)

				globalLogger.Infof(logNameSpace, "Adding log json output; filter=%s level=%s", filter, *rawLevel)
				globalLogger.AddOutput(
					filter,
					os.Stderr,
					logger.ParseLogLevel(*rawLevel, logger.LogLevelDebug),
					false,
					true,
					logger.NewJsonFormatter(),
				)

			default:
				panic("unknown log output type: " + outputType)
			}
		}
	}

	globalLogger.Infof(logNameSpace, "Creating global context with default logger level %v and namespace %v", logLevel, logNameSpace)
	globalContext = &ContextWithCancel{
		Context: ctx,
		Cancel:  cancel,
		Logger:  globalLogger,
	}

	return globalContext
}
