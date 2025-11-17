package jaeger_trace

import (
	"net/url"
	"testing"
)

func TestEscape(t *testing.T) {
	print(url.QueryEscape("adapt_in9gwarca6@484389a0b9ed796_in9gwarca6@53df7ad2afe8301/api/traces"))
}
