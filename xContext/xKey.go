package xContext

import "github.com/opentracing/opentracing-go/ext"

type XKey string

func (m XKey) ToExtTagName() ext.StringTagName {
	return ext.StringTagName(m)
}

func (m XKey) String() string {
	return string(m)
}
