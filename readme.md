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

See [test file](https://github.com/koykov/cbytebuf/blob/master/cbytebuf_test.go) for more examples and benchmarks for number of allocations.

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
BenchmarkCByteBuf_Write-8          	  345268	      3602 ns/op	       0 B/op	       0 allocs/op
BenchmarkCByteBuf_WriteLong-8      	    2622	    439151 ns/op	       0 B/op	       0 allocs/op
BenchmarkCByteBuf_AppendBytes-8    	 1373740	       870 ns/op	     896 B/op	       1 allocs/op
BenchmarkCByteBuf_AppendString-8   	 1342476	       869 ns/op	     896 B/op	       1 allocs/op
BenchmarkLBPool-8                  	11660017	       101 ns/op	       0 B/op	       0 allocs/op
BenchmarkPool-8                    	 5501479	       205 ns/op	       0 B/op	       0 allocs/op
```

Also you can see more comparison benchmarks in [versus](https://github.com/koykov/versus/tree/master/cbytebuf) project:
```
BenchmarkByteArray_Append-8             	  767320	      1449 ns/op	    2040 B/op	       8 allocs/op
BenchmarkByteArray_AppendLong-8         	    1557	    754013 ns/op	 4646288 B/op	      25 allocs/op
BenchmarkByteBufferNative_Write-8       	  517546	      2376 ns/op	    2416 B/op	       5 allocs/op
BenchmarkByteBufferNative_WriteLong-8   	    3441	    346512 ns/op	 1646722 B/op	      10 allocs/op
BenchmarkByteBufferPool_Write-8         	  904567	      1335 ns/op	       0 B/op	       0 allocs/op
BenchmarkByteBufferPool_WriteLong-8     	    1555	    754847 ns/op	 4667398 B/op	      29 allocs/op
BenchmarkCByteBuf_Write-8               	  380574	      3171 ns/op	       0 B/op	       0 allocs/op
BenchmarkCByteBuf_WriteLong-8           	    2631	    454779 ns/op	       0 B/op	       0 allocs/op
```

As you can see, CbyteBuf is slowest than any byte buffer or byte slice when writing short pieces of data, but has good speed for long writes.
Interesting that long writes is more faster that using append().

Anyway it's acceptable cost since it produces zero allocations even if you doesn't use pools. But I recommend to use it together with pool since it reduces amount of CGO calls in [cbyte](https://github.com/koykov/cbyte).
