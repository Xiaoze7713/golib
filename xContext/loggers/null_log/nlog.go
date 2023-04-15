package null_log

type LoggerNull struct{}

func (LoggerNull) Debug(specialKV map[string]interface{}, args ...interface{})              {}
func (LoggerNull) Info(specialKV map[string]interface{}, args ...interface{})               {}
func (LoggerNull) Warn(specialKV map[string]interface{}, args ...interface{})               {}
func (LoggerNull) Error(specialKV map[string]interface{}, args ...interface{})              {}
func (LoggerNull) Fatal(specialKV map[string]interface{}, args ...interface{})              {}
func (LoggerNull) Debugf(specialKV map[string]interface{}, fmt string, args ...interface{}) {}
func (LoggerNull) Warnf(specialKV map[string]interface{}, fmt string, args ...interface{})  {}
func (LoggerNull) Errorf(specialKV map[string]interface{}, fmt string, args ...interface{}) {}
func (LoggerNull) Fatalf(specialKV map[string]interface{}, fmt string, args ...interface{}) {}
func (LoggerNull) Infof(specialKV map[string]interface{}, fmt string, args ...interface{})  {}
func (LoggerNull) DebugKV(specialKV map[string]interface{}, KV map[string]interface{})      {}
func (LoggerNull) InfoKV(specialKV map[string]interface{}, KV map[string]interface{})       {}
func (LoggerNull) WarnKV(specialKV map[string]interface{}, KV map[string]interface{})       {}
func (LoggerNull) ErrorKV(specialKV map[string]interface{}, KV map[string]interface{})      {}
func (LoggerNull) FatalKV(specialKV map[string]interface{}, KV map[string]interface{})      {}
