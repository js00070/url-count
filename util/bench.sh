# go test -bench BenchmarkCountTop100WithBufferSize100M -benchmem -cpuprofile 100M_cpu.out -memprofile 100M_mem.out

# go test -bench BenchmarkCountTop100WithBufferSize200M -benchmem -cpuprofile 200M_cpu.out -memprofile 200M_mem.out

# go test -bench BenchmarkCountTop100WithBufferSize500M -benchmem -cpuprofile 500M_cpu.out -memprofile 500M_mem.out

# go test -bench BenchmarkParallelSort -benchmem -cpuprofile parallel_sort_cpu.out -memprofile parallel_sort_mem.out

go test -bench BenchmarkCountTop100In10GBWithBufferSize300M -benchmem -cpuprofile 10GB_300M_cpu.out -memprofile 10GB_300M_mem.out

# go test -bench BenchmarkCountTop100In10GBWithBufferSize500M -benchmem -cpuprofile 10GB_500M_cpu.out -memprofile 10GB_500M_mem.out

# go test -bench BenchmarkCountTop100In10GBWithBufferSize1GB -benchmem -cpuprofile 10GB_1GB_cpu.out -memprofile 10GB_1GB_mem.out

