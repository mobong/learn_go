#### map的结构体
```go
type hmap struct {
	count     int               // 当前哈希表中的元素数量
	flags     uint8             // 写标志
	B         uint8             // buckets数组的长度的对数
	noverflow uint16            // overflow的bucket近似数
	hash0     uint32            // 哈希表种子
	buckets    unsafe.Pointer   // 指向bmap数组，大小为2^B
	oldbuckets unsafe.Pointer   // 扩容时保存之前buckets,大小是当前buckets的一半
	nevacuate  uintptr          // 指示扩容进度,小于此地址的buckets表示完成迁移
	extra *mapextra             // 溢出桶数据
}

type bmap struct {
    // tophash [bucketCnt]uint8  键的哈希高8位,实际结构不止这个字段
	topbits [8]uint8
	keys    [8]keytype
	values  [8]valuetype
	pad     uintptr
	overflow uintptr
}
```
* map的扩容条件有两个：
    1.装载因子(元素数量/桶数量)超过阈值(6.5)
    2.溢出桶过多：当B<15时，溢出桶数量超过2^B；当B>=15时，溢出桶数量超过2^15
* 针对条件1，元素过多了，map的扩容是翻倍扩容；条件2是溢出桶过多，map的扩容是sameSizeGrow(等量扩容),可以把溢出桶里面的数据处理的更紧密点
* map的扩容需要将原有的key/value重新搬迁到新的内存地址，大量的key/value一次性搬迁会非常影响性能，因此采用了"渐进式"的方式进行扩容,且每次最多只会搬迁2个buckets
