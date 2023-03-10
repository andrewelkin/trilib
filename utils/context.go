package utils

import (
	"context"
	"github.com/andrewelkin/trilib/utils/logger"
	"regexp"
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

func SetGlobalLogger(logger logger.Logger) logger.Logger {
	if globalContext != nil {
		globalContext.Logger = logger
	} else {
		logger.Warnf("*", "Unable to set global logger")
	}
	return logger
}

// -----------------------------------------------------------------------------------------[AE: 2023-03-1]-----------/

func GetGlobalContext() *ContextWithCancel {
	return globalContext
}

// GetOrCreateGlobalContext sets a new global context with logging and cancel
// Expects a config, which is normally would be a "Logging" section
func GetOrCreateGlobalContext(gconfig *Vconfig) *ContextWithCancel {
	if globalContext != nil {
		return globalContext
	}
	ctx, cancel := context.WithCancel(context.Background())
	var config *Vconfig

	if gconfig != nil {
		config = gconfig.FromKey("logger")
	}

	var logLevel logger.LogLevel
	logLevel = logger.LogLevelDebug
	logNameSpace := "*"

	var outputsCfg *Vconfig
	if config != nil {
		if tmp := config.GetString("loglevel"); tmp != nil {
			logLevel = logger.ParseLogLevel(*tmp, logger.LogLevelDebug)
		}

		// allow the default stdout log namespace filter to be overridden by the "filter" config field
		logger.LogDefaultFilter = regexp.MustCompile(*config.GetStringDefault("filter", logger.LogDefaultFilter.String()))
		outputsCfg = config.FromKey("outputs")
	}

	// create the main global logger (default writes to STDOUT)
	globalLogger := logger.GetOrCreateGlobalLogger(ctx, logLevel)

	// add additional writers, if configured
	if outputsCfg != nil {
		allOutputCfg := outputsCfg.GetCfg()
		for outputType := range allOutputCfg {

			cfg := outputsCfg.FromKey(outputType)
			if cfg == nil {
				panic("failed to parse output config: " + outputType)
			}

			switch strings.ToLower(outputType) {
			case "filewriter":
				rawLevel, path, prefix, suffix, rawFilter, skipRepeating :=
					cfg.GetStringDefault("logLevel", "debug"),
					cfg.GetStringDefault("path", "/tmp/test_logs"),
					cfg.GetStringDefault("filePrefix", ""),
					cfg.GetStringDefault("fileSuffix", ".log"),
					cfg.GetStringDefault("filter", logger.FilterMatchAll.String()),
					cfg.GetBoolDefault("skipRepeating", true)

				fileWriter, err := logger.NewFileWriter(*path, prefix, suffix, skipRepeating)
				if err != nil {
					panic("failed to create file writer: " + err.Error())
				}
				filter, err := regexp.Compile(*rawFilter)
				if err != nil {
					panic("failed to compile log filter regexp: " + err.Error())
				}

				globalLogger.Infof(logNameSpace, "Adding log file output; filter=%s path=%s level=%s", *rawFilter, *path, *rawLevel)
				globalLogger.AddOutput(
					filter,
					fileWriter,
					logger.ParseLogLevel(*rawLevel, logger.LogLevelDebug),
				)
				continue

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
