package logger

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	"cloud.google.com/go/logging"
)

func TestStdLogger(t *testing.T) {
	var buf bytes.Buffer
	defaultFlags := log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		defer func() {
			log.SetOutput(os.Stderr)
			log.SetFlags(defaultFlags)
		}()
	})

	t.Run("should print a string log with formatting", func(t *testing.T) {
		t.Cleanup(func() {
			buf.Reset()
		})
		newStdLoggerFactory().Logger("test").Debug(context.Background(), logging.Entry{
			Payload: "payload",
		})
		if buf.String() != "[test][DEBUG] \"payload\"\n" {
			t.Errorf("got %s", buf.String())
		}
	})

	t.Run("should print a json log with formatting", func(t *testing.T) {
		t.Cleanup(func() {
			buf.Reset()
		})
		newStdLoggerFactory().Logger("test").Debug(context.Background(), logging.Entry{
			Payload: map[string]any{
				"key": "value",
			},
		})
		if buf.String() != "[test][DEBUG] {\"key\":\"value\"}\n" {
			t.Errorf("got %s", buf.String())
		}
	})
}
