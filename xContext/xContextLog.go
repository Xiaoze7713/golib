package xContext

const (
	LogFieldTraceID = "trace_id"
	LogFieldSpanID  = "span_id"
)

func (x *XContext) InfoKV(kvm KVMType) {
	x.LogFields(kvm.ToAnyArgs()...)
	x.LoggerIF.Info(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, x.KVsFormats(kvm))
}

func (x *XContext) Info(args ...interface{}) {
	x.LoggerIF.Info(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, append([]interface{}{}, x.ArgsFormats(args...))...)
}

func (x *XContext) Debug(args ...interface{}) {
	x.LoggerIF.Debug(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, append([]interface{}{}, x.ArgsFormats(args...))...)
}

func (x *XContext) Warn(args ...interface{}) {
	x.LoggerIF.Warn(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, append([]interface{}{}, x.ArgsFormats(args...))...)
}

func (x *XContext) Error(args ...interface{}) {
	x.LoggerIF.Error(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, append([]interface{}{}, x.ArgsFormats(args...))...)
}

func (x *XContext) Fatal(args ...interface{}) {
	x.LoggerIF.Fatal(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, append([]interface{}{}, x.ArgsFormats(args...))...)
}

func (x *XContext) Infof(format string, args ...interface{}) {
	//format = x.XContextPrefix() + format
	x.LoggerIF.Infof(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, format, args...)
}

func (x *XContext) Debugf(format string, args ...interface{}) {
	//format = x.XContextPrefix() + format
	x.LoggerIF.Debugf(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, format, args...)
}

func (x *XContext) Warnf(format string, args ...interface{}) {
	//format = x.XContextPrefix() + format
	x.LoggerIF.Warnf(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, format, args...)
}

func (x *XContext) Errorf(format string, args ...interface{}) {
	//format = x.XContextPrefix() + format
	x.LoggerIF.Errorf(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, format, args...)
}

func (x *XContext) Fatalf(format string, args ...interface{}) {
	//format = x.XContextPrefix() + format
	x.LoggerIF.Fatalf(map[string]interface{}{LogFieldTraceID: x.TraceID().String(), LogFieldSpanID: x.SpanID().String()}, format, args...)
}
