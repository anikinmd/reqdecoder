package reqdecoder

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_decoderDeflate(t *testing.T) {
	testBody := strings.Repeat("test", 10)

	var buf bytes.Buffer
	df, _ := flate.NewWriter(&buf, 5)
	_, _ = df.Write([]byte(testBody))
	_ = df.Close()
	compressed := buf.String()

	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantStr string
	}{
		{
			name:    "Success decode",
			input:   compressed,
			wantNil: false,
			wantStr: testBody,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := io.NopCloser(strings.NewReader(tt.input))
			res := decoderDeflate(rc)
			if tt.wantNil {
				assert.Nil(t, res)
			} else {
				assert.NotNil(t, res)
				buf, err := io.ReadAll(res)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStr, string(buf))
			}
		})
	}
}

func Test_decoderGzip(t *testing.T) {
	testBody := strings.Repeat("test", 10)

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, _ = gz.Write([]byte(testBody))
	_ = gz.Close()
	compressed := buf.String()

	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantStr string
	}{
		{
			name:    "Success decode",
			input:   compressed,
			wantNil: false,
			wantStr: testBody,
		},
		{
			name:    "Not compressed body",
			input:   testBody,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := io.NopCloser(strings.NewReader(tt.input))
			res := decoderGzip(rc)
			if tt.wantNil {
				assert.Nil(t, res)
			} else {
				assert.NotNil(t, res)
				buf, err := io.ReadAll(res)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStr, string(buf))
			}
		})
	}
}
