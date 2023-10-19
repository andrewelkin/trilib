package utils

import (
	"context"
	"github.com/andrewelkin/trilib/utils/logger"
	"github.com/golang/mock/gomock"
	"github.com/nats-io/nats.go"
	"strings"
	"testing"
	"time"
)

func TestGetOrCreateGlobalContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	acceptString := "accept"   // The logger will only accept this namespace
	excludeString := "exclude" // The logger will exclude this namespace
	underscoreString := "_testNS"

	tests := []struct {
		Actor    func(*ContextWithCancel)
		Result   string
		GFilter  *string
		GExclude *string
	}{
		{func(ctx *ContextWithCancel) {
			ctx.Logger.Debugf(acceptString, "test")
			ctx.Logger.Debugf(excludeString, "test")
			ctx.Logger.Debugf(underscoreString, "test")
		}, "accept test\nexclude test\n", nil, nil},
		{func(ctx *ContextWithCancel) {
			ctx.Logger.Debugf(acceptString, "test")
			ctx.Logger.Debugf(excludeString, "test")
			ctx.Logger.Debugf(underscoreString, "test")
		}, "accept test\n", &acceptString, nil},
		{func(ctx *ContextWithCancel) {
			ctx.Logger.Debugf(acceptString, "test")
			ctx.Logger.Debugf(excludeString, "test")
			ctx.Logger.Debugf(underscoreString, "test")

		}, "", &acceptString, &acceptString},
	}
	for _, test := range tests {
		buf := new(strings.Builder)
		logger.SetDefaultScreenIO(buf)
		globalContext = nil
		var gfilter logger.FilterFunc

		globalLogger, natsLogger := logger.NewMockLogger(ctrl), logger.NewMockLogger(ctrl)
		globalLogger.EXPECT().AddOutput(gomock.Any(), gomock.Any(), logger.LogLevel(0), false, false).Do(
			func(filter logger.FilterFunc, output interface{}, minLevel logger.LogLevel, ansi bool, trailCR bool, opts ...interface{}) {
				gfilter = filter
			})
		globalLogger.EXPECT().Infof("*", "Creating global context with default logger level %v and namespace %v", gomock.Any(), gomock.Any())
		globalLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).Do(func(namespace, format string, a ...interface{}) {
			if gfilter(namespace) {
				buf.WriteString(namespace + " " + format + "\n")
			}
		}).AnyTimes()
		ctx := CreateMockedContext(ctrl, test.GFilter, test.GExclude, nil, nil, []*logger.MockLogger{globalLogger, globalLogger, natsLogger})
		test.Actor(ctx)

		time.Sleep(time.Millisecond * 100)

		if buf.String() != test.Result {
			t.Errorf("Unexpected output: `%s` vs `%s`", buf.String(), test.Result)
		}
	}
}

func CreateMockedContext(ctrl *gomock.Controller, natsFilter, natsExclude, gFilter, gExclude *string, loggers []*logger.MockLogger) (ctx *ContextWithCancel) {
	frog := "frog"
	debug := "debug"
	empty := ""
	url := nats.DefaultURL

	mockedNatsConfig := NewMockIConfig(ctrl)
	mockedNatsConfig.EXPECT().GetStringDefault("subject", "default-logger-subject").Return(&frog)
	mockedNatsConfig.EXPECT().GetStringDefault("url", nats.DefaultURL).Return(&url)
	mockedNatsConfig.EXPECT().GetString("filter").Return(natsFilter).AnyTimes()
	mockedNatsConfig.EXPECT().GetString("exclude").Return(natsExclude).AnyTimes()
	mockedNatsConfig.EXPECT().GetStringDefault("logLevel", "debug").Return(&debug).AnyTimes()
	mockedNatsConfig.EXPECT().GetBoolDefault("ansicodes", false).Return(false).AnyTimes()
	mockedNatsConfig.EXPECT().GetStringDefault("natsSeed", "").Return(&empty).AnyTimes()

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
