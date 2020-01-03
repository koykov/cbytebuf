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
