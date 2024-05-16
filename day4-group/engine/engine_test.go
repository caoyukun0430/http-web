package engine

import "testing"

func TestNestedGroup(t *testing.T) {
	engine := New()
	v1 := engine.Group("/v1")
	v2 := v1.Group("/v2")
	v3 := v2.Group("/v3")
	if v2.prefix != "/v1/v2" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
	if v3.prefix != "/v1/v2/v3" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
	// init engine.groups have size 1
	if len(engine.groups) != 4 {
		t.Fatal("engine group size is not correct,should be 4")
	}
}
