package reqdecoder

import (
	"compress/flate"
	"compress/gzip"
	"io"
)

func (d *requestDecoder) addDecoder(contentEncoding string, decoder decoderFunc) {
	d.decoders[contentEncoding] = decoder
}

func (d *requestDecoder) fillDecoders() {
	d.addDecoder("gzip", decoderGzip)
	d.addDecoder("deflate", decoderDeflate)
}

func decoderGzip(reader io.ReadCloser) io.ReadCloser {
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil
	}
	return gz
}

func decoderDeflate(reader io.ReadCloser) io.ReadCloser {
	// Deflate doesn't return error for some reason
	df := flate.NewReader(reader)
	return df
}
