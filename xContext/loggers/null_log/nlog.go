package null_log

type LoggerNull struct{}

func (LoggerNull) Debug(args ...interface{})              {}
func (LoggerNull) Info(args ...interface{})               {}
func (LoggerNull) Warn(args ...interface{})               {}
func (LoggerNull) Error(args ...interface{})              {}
func (LoggerNull) Fatal(args ...interface{})              {}
func (LoggerNull) Debugf(fmt string, args ...interface{}) {}
func (LoggerNull) Warnf(fmt string, args ...interface{})  {}
func (LoggerNull) Errorf(fmt string, args ...interface{}) {}
func (LoggerNull) Fatalf(fmt string, args ...interface{}) {}
func (LoggerNull) Infof(fmt string, args ...interface{})  {}
