package mock_logger

import "vm/pkg/constants"

type StubLogger struct{}

func (l *StubLogger) Debug(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
}
func (l *StubLogger) Info(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
}
func (l *StubLogger) Warn(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
}
func (l *StubLogger) Error(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
}
func (l *StubLogger) Fatal(cat constants.Category, sub constants.SubCategory, msg string, extra map[constants.ExtraKey]interface{}) {
}

func (l *StubLogger) Debugf(template string, args ...interface{}) {}
func (l *StubLogger) Infof(template string, args ...interface{})  {}
func (l *StubLogger) Warnf(template string, args ...interface{})  {}
func (l *StubLogger) Errorf(template string, args ...interface{}) {}
func (l *StubLogger) Fatalf(template string, args ...interface{}) {}
func (l *StubLogger) Init() {
	// No-op or simple console logger for test visibility
}
