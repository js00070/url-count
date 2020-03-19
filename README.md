# url-count
从100GB的文本文件中找到出现次数最多的前100个url，内存限制1GB

## 思路
首先是需要求出每个url在文件中出现的次数，相当于是算一个group by聚合。有两种方式可以做，一种是基于hash的聚合，一种是基于sort的聚合。

基于hash的聚合：
1. 先用一个哈希函数h0将所有的具有相同h0的url进行分组(比如分成100组)，写入磁盘
2. 如果一个组的大小过大(比如大于1GB)，则需要用哈希函数h1对这个组再次进行分组
3. 分组完成后，相同的url必然在同一个组中，则针对每个组，利用内存中的hashmap计算聚合函数

基于sort的聚合：
1. 按顺序分块读入url，在内存中对每一块进行排序，然后分块写入磁盘
2. 利用堆将这些内部有序的块合并成总体有序的序列，相同的url会排在一起，迭代计算聚合函数

这两种方式都是总共需要读两次磁盘，写一次磁盘

## 实现

由于感觉sort-aggregate的内存使用相对于hash-aggregate来说比较可控，所以选择了实现基于sort的聚合。

程序入口是util/count.go中的CountTopN函数，里面的流程是这样的：
1. 执行了PartitionSort函数，完成了分块排序
2. 创建MergeSorter，不断地调用MergeSorter.Next获取从小到大字典序的[]byte流，同时迭代计算count
3. 维护一个容量为n的最小堆counterHeap，得到count前n大的记录

util/sorting.go中PartitionSort函数的流程是这样的：
1. 创建两块buffer，一块用于磁盘io，一块用于排序，交换两块buffer的指针来实现任务交接
2. 创建读取磁盘的协程，与排序并行执行

## 性能优化

基本功能实现后，又进行了一些优化：
1. 实现了并行的快速排序, 见parallel_sort.go，实现的思路是将快排中的partition动作抽离出来并行执行
2. 一开始写的时候比较随意，制造了很多的内存碎片，而golang的GC并不能及时清除，于是就利用pprof把所有能够复用的内存都复用了，不再有巨量的内存碎片产生.
3. 上面提到的，分出两块buffer让读取磁盘与计算排序并行执行

## 改进空间

1. 在PartitionSort中，写磁盘也可以与排序并行进行
2. 用统一的内存池实现内存复用
3. 完善异常处理
4. 尝试hash-aggregate

## 测试
```bash
$ cd util
$ go test -bench BenchmarkCountTop100In10GBWithBufferSize300M -benchmem -cpuprofile 10GB_300M_cpu.out -memprofile 10GB_300M_mem.out
```
首次运行时会先随机生成10GB大小的测试文件，然后使用300M大小的两个buffer进行测试。

不算首次运行时的生成过程，在我的机器上测试下来是花了103秒且占用内存始终保持在700多M，而一开始没有进行任何优化时，是需要花大概160多秒而且内存占用极其恐怖(golang GC的问题)