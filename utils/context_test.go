package utils

import (
	"context"
	"github.com/andrewelkin/trilib/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/nats.go"
	"testing"
	"time"
)

func TestGetOrCreateGlobalContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	globalLogger, natsLogger := logger.NewMockLogger(ctrl), logger.NewMockLogger(ctrl)
	globalLogger.EXPECT().AddOutput(gomock.Any(), gomock.Any(), logger.LogLevel(0), false, false)
	globalLogger.EXPECT().Infof("*", "Creating global context with default logger level %v and namespace %v", gomock.Any(), gomock.Any())
	globalLogger.EXPECT().Debugf("*", "Global logger created")
	ctx := CreateMockedContext(ctrl, nil, nil, nil, nil, []*logger.MockLogger{globalLogger, globalLogger, natsLogger})

	ctx.Logger.Debugf("*", "Global logger created")
	time.Sleep(time.Millisecond * 100)
}

func CreateMockedContext(ctrl *gomock.Controller, natsFilter, natsExclude, gFilter, gExclude *string, loggers []*logger.MockLogger) (ctx *ContextWithCancel) {
	frog := "frog"
	debug := "debug"
	empty := ""
	url := nats.DefaultURL

	mockedNatsConfig := NewMockIConfig(ctrl)
	mockedNatsConfig.EXPECT().GetStringDefault("subject", "default-logger-subject").Return(&frog)
	mockedNatsConfig.EXPECT().GetStringDefault("url", nats.DefaultURL).Return(&url)
	mockedNatsConfig.EXPECT().GetString("filter").Return(natsFilter)
	mockedNatsConfig.EXPECT().GetString("exclude").Return(natsExclude)
	mockedNatsConfig.EXPECT().GetStringDefault("logLevel", "debug").Return(&debug)
	mockedNatsConfig.EXPECT().GetBoolDefault("ansicodes", false).Return(false)
	mockedNatsConfig.EXPECT().GetStringDefault("natsSeed", "").Return(&empty)

	mockedOutputsConfig := NewMockIConfig(ctrl)
	mockedOutputsConfig.EXPECT().GetCfg().Return(map[string]interface{}{
		"nats": mockedNatsConfig,
	}).AnyTimes()
	mockedOutputsConfig.EXPECT().FromKey("nats").Return(mockedNatsConfig).AnyTimes()

	mockedOutputs := NewMockIConfig(ctrl)
	loglevel := "debug"
	mockedOutputs.EXPECT().GetString("loglevel").Return(&loglevel).AnyTimes()
	mockedOutputs.EXPECT().GetString("filter").Return(gFilter).AnyTimes()
	mockedOutputs.EXPECT().GetString("exclude").Return(gExclude).AnyTimes()
	mockedOutputs.EXPECT().FromKey("outputs").Return(mockedOutputsConfig).AnyTimes()

	config := NewMockIConfig(ctrl)
	config.EXPECT().FromKey("logger").Return(mockedOutputs).AnyTimes()

	ctx = GetOrCreateGlobalContext(config,
		func(context.Context, logger.LogLevel, logger.FilterFunc) logger.Logger {
			logger := loggers[0]
			loggers = loggers[1:]
			return logger
		}, logger.FilterUnderscore)
	return
}
