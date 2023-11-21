## slicle扩容学习

* 使用append追加元素的时候超过容量就会对切片进行扩容，扩容代码如下 runtime/slice.go(基于go 1.20)
```go
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
    newcap := oldCap
	doublecap := newcap + newcap
	if newLen > doublecap {
		newcap = newLen
	} else {
		const threshold = 256
		if oldCap < threshold {
			newcap = doublecap
		} else {
			for 0 < newcap && newcap < newLen {
				newcap += (newcap + 3*threshold) / 4
			}
			if newcap <= 0 {
				newcap = newLen
			}
		}
	}
	}
```
* 如果newLen(oldCap+num)大于当前容量的2倍,就使用当前newLen
* 如果旧切片容量小于256，就会将容量翻倍
* 如果旧切片容量大于等于256，就会每次增加25%的容量+192(当newcap越来越大时，实际上是趋向于1.25倍)，直至新容量大于newLen
#### 这时候还不是最终的容量，需要根据切片中的元素所占字节大小(1、2或8的倍数)进行对齐内存，默认情况下是将目标容量和元素大小相乘得到占用的内存，对齐内存代码在runtime/msize.go：
```go
func roundupsize(size uintptr) uintptr {
	if size < _MaxSmallSize {
		if size <= smallSizeMax-8 {
			return uintptr(class_to_size[size_to_class8[divRoundUp(size, smallSizeDiv)]])
		} else {
			return uintptr(class_to_size[size_to_class128[divRoundUp(size-smallSizeMax, largeSizeDiv)]])
		}
	}
	if size+_PageSize < size {
		return size
	}
	return alignUp(size, _PageSize)
}
```
* roundupsize是对待申请的内存进行向上取整，使用class_to_size数组中的整数可以提高贴在分配效率并减少碎片
### 如果扩容后的新切片不会赋值回原变量，就需要注意下会发生复制从而影响性能
