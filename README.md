# storage
A memory storage library for Go

## Benchmarks

### Writes and reads

```bash
go version
go version go1.23.0 linux/amd64

go test -bench=. -benchmem -benchtime=4s ./... -timeout 30m
goos: linux
goarch: amd64
pkg: github.com/vlecluse/storage
cpu: Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
BenchmarkStorage_Set-8         	114962505	        68.83 ns/op	       0 B/op	       0 allocs/op
BenchmarkStorage_SetMedium-8   	124681840	        32.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkStorage_SetBig-8      	117808448	        38.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkStorage_Get-8         	673117422	         6.384 ns/op	       0 B/op	       0 allocs/op
BenchmarkStorage_Delete-8      	766807340	         5.906 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/vlecluse/storage	40.584s
```
