package main

import (
	"bytes"
	"github.com/Clever/gearadmin"
	"testing"
)

func TestWriteStatus(t *testing.T) {
	var out bytes.Buffer
	c := NewDebugRaidmanClient(&out)
	status := gearadmin.Status{
		Function:         "fn",
		Total:            10,
		Running:          5,
		AvailableWorkers: 6,
	}
	tags := []string{"tag"}
	writeStatus(c, status, tags)
	expected := `{"Ttl":10,"Time":0,"Tags":["tag"],"Host":"","State":"","Service":"gearmand.status.fn.total","Metric":10,"Description":"","Attributes":null}
{"Ttl":10,"Time":0,"Tags":["tag"],"Host":"","State":"","Service":"gearmand.status.fn.running","Metric":5,"Description":"","Attributes":null}
{"Ttl":10,"Time":0,"Tags":["tag"],"Host":"","State":"","Service":"gearmand.status.fn.available","Metric":6,"Description":"","Attributes":null}
`
	if out.String() != expected {
		t.Fatalf("'%s' != '%s'", out.String(), expected)
	}
}
