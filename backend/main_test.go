package main

import (
	"github.com/czQery/llg/backend/tl"
	"testing"
	"time"
)

const timeLayout = "02.01.2006-15:04"

func TestParseTime(t *testing.T) {
	list := []string{"01.02.2024-05:52", "01.02.2024-5:52", "1.2.2024-5:52", "30.12.3030-23:59", "01.01.1970-00:00", "00.00.0000-00:00"}
	for _, item := range list {
		t1, _ := time.Parse(timeLayout, item)
		t2 := tl.ParseTime(item)

		t.Log("COMPARE", t1.Unix(), t2.Unix())
	}

	t.Fail()
}
