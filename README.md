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
With chi:
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