package main

import "log"
import "testing"

type sinkWriter struct {
	data []byte
}

func (sw *sinkWriter) Write(p []byte) (n int, err error) {
	sw.data = append(sw.data, p...)
	return len(p), nil
}

func testPWOutput(t *testing.T, prefix string, writes []string, expected_output string) {
	t.Logf("Checking prefix=%#v writes=%#v", prefix, writes)
	sink := &sinkWriter{}
	pw := NewPrefixingWriter(sink, prefix)
	for _, datum := range writes {
		dbytes := []byte(datum)
		n, err := pw.Write(dbytes)
		if err != nil {
			t.Error("Error writing to PW", err)
		}
		if n != len(dbytes) {
			t.Errorf("Wrote %d bytes, returned %d", len(dbytes), n)
		}
	}
	pw.Close()
	output := string(sink.data)
	if output != expected_output {
		t.Errorf("Expected %#v, got %#v", expected_output, output)
	}
}

func TestSingleLines(t *testing.T) {
	testPWOutput(t, "PREFIX",
		[]string{ "foo\n", "bar\n", "baz\n" },
		"PREFIX foo\nPREFIX bar\nPREFIX baz\n")
}

func TestNoFinalNewline(t *testing.T) {
	testPWOutput(t, "PREFIX",
		[]string{ "foo\n", "bar\n", "baz" },
		"PREFIX foo\nPREFIX bar\nPREFIX baz\n")
}

func TestBlocksAcrossLines(t *testing.T) {
	testPWOutput(t, "PREFIX",
		[]string{ "foo\nba", "r\nbaz" },
		"PREFIX foo\nPREFIX bar\nPREFIX baz\n")
}
