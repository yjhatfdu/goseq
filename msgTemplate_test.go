package goseq

import (
	"testing"
)

func TestRenderMsgTemplate(t *testing.T) {
	ret := RenderMsgTemplate("hello{name}{test}{name}", map[string]string{"name": "world"})
	if ret != "helloworld{test}world" {
		t.Fail()
	}
}

type test struct {
	f1 string
	f2 int
}

func TestRenderMsgTemplateStruct(t *testing.T) {
	ret := RenderMsgTemplate("hello{name}{test}{name}{@struct}", map[string]string{"name": "world", "struct": "{\"foo\":1,\"bar\":\"ok\"}"})
	if ret != "helloworld{test}world{foo: 1, bar: ok}" {
		t.Fail()
	}
}
