package utils

import (
	"context"
	"github.com/andrewelkin/trilib/utils/logger"
	"github.com/golang/mock/gomock"
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
	emptystring := ""

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
		}, "accept test\nexclude test\n_testNS test\n", &emptystring, &emptystring},
		{func(ctx *ContextWithCancel) {
			ctx.Logger.Debugf(acceptString, "test")
			ctx.Logger.Debugf(excludeString, "test")
			ctx.Logger.Debugf(underscoreString, "test")
		}, "accept test\n", &acceptString, &emptystring},
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
		globalLogger.EXPECT().AddOutput(gomock.Any(), gomock.Any(), logger.LogLevel(0), false, true).Do(
			func(filter logger.FilterFunc, output interface{}, minLevel logger.LogLevel, ansi bool, trailCR bool, opts ...interface{}) {
				gfilter = filter
			})
		globalLogger.EXPECT().Infof("*", "Creating global context with default logger level %v and namespace %v", gomock.Any()).After(
			globalLogger.EXPECT().Infof("*", "Adding log file output; filter=%s path=%s level=%s", gomock.Any()).MaxTimes(1),
		)
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

func CreateMockedContext(ctrl *gomock.Controller, filesFilter, filesExclude, globalFilter, globalExclude *string, loggers []*logger.MockLogger) (ctx *ContextWithCancel) {
	mkStr := func(s string) *string {
		r := s
		return &r
	}
	mockedFileWriter := NewMockIConfig(ctrl)
	mockedFileWriter.EXPECT().GetStringDefault("logLevel", "debug").Return(mkStr("frog"))
	mockedFileWriter.EXPECT().GetStringDefault("path", "/tmp/test_logs").Return(mkStr("/tmp/test_logs")).AnyTimes()
	mockedFileWriter.EXPECT().GetStringDefault("filePrefix", "").Return(mkStr("")).AnyTimes()
	mockedFileWriter.EXPECT().GetStringDefault("fileSuffix", ".log").Return(mkStr(".log")).AnyTimes()
	mockedFileWriter.EXPECT().GetBoolDefault("skipRepeating", true).Return(true).AnyTimes()
	mockedFileWriter.EXPECT().GetString("filter").Return(filesFilter).AnyTimes()
	mockedFileWriter.EXPECT().GetString("exclude").Return(filesExclude).AnyTimes()

	mockedOutputsConfig := NewMockIConfig(ctrl)
	mockedOutputsConfig.EXPECT().GetCfg().Return(map[string]interface{}{
		"fileWriter": mockedFileWriter,
	}).AnyTimes()
	mockedOutputsConfig.EXPECT().FromKey("fileWriter").Return(mockedFileWriter).AnyTimes()

	mockedOutputs := NewMockIConfig(ctrl)
	loglevel := "debug"
	mockedOutputs.EXPECT().GetString("loglevel").Return(&loglevel).AnyTimes()
	mockedOutputs.EXPECT().GetString("filter").Return(globalFilter).AnyTimes()
	mockedOutputs.EXPECT().GetString("exclude").Return(globalExclude).AnyTimes()
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
