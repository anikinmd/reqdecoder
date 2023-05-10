[![codecov](https://codecov.io/gh/anikinmd/reqdecoder/branch/main/graph/badge.svg?token=58tQ2DpwRf)](https://codecov.io/gh/anikinmd/reqdecoder)
[![Go Report github.com/anikgithub.com/anikgithub.com/anikinmd/reqdecompinmd/reqdecompinmd/reqdecompCard](https://goreportcard.com/badge/github.com/anikinmd/reqdecoder)](https://goreportcard.com/report/github.com/anikinmd/reqdecoder)
# reqdecoder
Go middleware for decoding/decompressing request body

Features:
* GZip
* Deflate
* Custom decoders

## Usage:
Basic:
```
http.Handle("/", reqdecomp.RequestDecoder(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	...
	w.WriteHeader(http.StatusOK)
})))
```
Using chi:
```
r := chi.NewRouter()
r.Use(reqdecomp.RequestDecoder)
```
Custom decoder:
```
// Decoder should return nil on error
func customDecoder(reader io.ReadCloser) io.ReadCloser {
	...
	return newReader
}
...
reqdecomp.AddDecoder("customContentType", customDecoder)
```