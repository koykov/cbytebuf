# CbyteBuf

Alloc-free replacement for [bytes.Buffer](https://golang.org/pkg/bytes/#Buffer) based on [cbyte](https://github.com/koykov/cbyte).

## Usage

```go
package main

import (
	"fmt"
	"github.com/koykov/cbytebuf"
)

func main() {
	buf := cbytebuf.NewCByteBuf()
	defer buf.Release()
	buf.WriteString("foo ")
	buf.WriteString("bar ")
	// ...
	buf.WriteString("end.")
	fmt.Println(buf.String()) // "foo bar ... end."
}
```

No escapes to heap:

```bash
$ go build -gcflags '-m' example.go 
# command-line-arguments
example/example.go:9:9: inlining call to cbytebuf.NewCByteBuf
example/example.go:9:9: main &cbytebuf.b·2 does not escape
```

See [test file](https://github.com/koykov/cbytebuf/blob/master/cbytebuf_test.go) for more examples and benchmarks for number pf allocations.

## How it works

This package was inspired by article [Allocation efficiency in high-performance Go services](https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/).
Please read it before continue.

If you will use a lot of bytes.Buffer (or any analogues) you may notice that GC pressure will increase during the time
even if you use sync.Pool. This occurs since all slices in the pools (or any storage) checks by GC during mark phase.

The main approach of CbyteBuf is to avoid using any references and pointers inside it and, consecutive, avoid escapes to heap.
In fact the instance of CbyteBuf contains only [SliceHeader](https://golang.org/pkg/reflect/#SliceHeader) and temporary int variable - one uintptr and three integers in result.
As result any new instance of CBB allocates in stack instead of heap.
In fact allocations in heap occurs, but they produces by [cbyte](https://github.com/koykov/cbyte) and GC doesn't know nothing about them.

We've experienced increasing in more than 2 times the intervals between GC cycles, that is very good for our project. Also we noticed decreasing of GC CPU usage in ~3 times.

## Benchmarks

```
BenchmarkCByteBuf_Write-8                 500000      3430 ns/op       0 B/op       0 allocs/op
BenchmarkCByteBuf_WriteLong-8               3000    439168 ns/op       0 B/op       0 allocs/op
BenchmarkByteSlice_Append-8              1000000      1534 ns/op    2040 B/op       8 allocs/op
BenchmarkByteSlice_AppendLong-8             2000    631983 ns/op 4646289 B/op      25 allocs/op
BenchmarkByteBufferNative_Write-8         500000      2657 ns/op    2416 B/op       5 allocs/op
BenchmarkByteBufferNative_WriteLong-8       5000    308436 ns/op 1646724 B/op      10 allocs/op
BenchmarkByteBufferValyala_Write-8       1000000      1553 ns/op    2040 B/op       8 allocs/op
BenchmarkByteBufferValyala_WriteLong-8      2000    666237 ns/op 4646282 B/op      25 allocs/op
```

As you can see, CbyteBuf is slowest than any byte buffer or byte slice when writing short pieces of data, but has good speed for long writes.
Interesting that long writes is more faster that using append().

Anyway it's acceptable cost since it produces zero allocations even if you doesn't use pools. But I recommend to use it together with pool since it reduces amount of CGO calls in [cbyte](https://github.com/koykov/cbyte).
