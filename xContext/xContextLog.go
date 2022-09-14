package xContext

import "github.com/go-stack/stack"

func (x *XContext) InfoKV(kvm KVMType) {
	x.LogFields(kvm.ToAnyArgs()...)
	x.LoggerIF.Info(append([]interface{}{x.XContextPrefix()}, x.KVsFormats(kvm))...)
}

func (x *XContext) Info(args ...interface{}) {
	x.LoggerIF.Info(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) Debug(args ...interface{}) {
	x.LoggerIF.Debug(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) Warn(args ...interface{}) {
	x.LoggerIF.Warn(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) Error(args ...interface{}) {
	x.LoggerIF.Error(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) ErrorE(e error) {
	x.SetTag(string(ExecStatus), "error")
	x.LogFields(string(Error), e.Error())
	x.MetricsIF.CounterBy("err", x.TagNames2PLabels(Error)).Add(1)
	x.LoggerIF.Error(append([]interface{}{x.XContextPrefix()}, e.Error(), stack.Caller(7))...)
}

func (x *XContext) Fatal(args ...interface{}) {
	x.LoggerIF.Fatal(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) Infof(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Infof(format, args...)
}

func (x *XContext) Debugf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Debugf(format, args...)
}

func (x *XContext) Warnf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Warnf(format, args...)
}

func (x *XContext) Errorf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Errorf(format, args...)
}

func (x *XContext) Fatalf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Fatalf(format, args...)
}
