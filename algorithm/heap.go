package algorithm

import (
	"sync"
)

var (
	heapInitCap  = 100
	heapCleanCap = 300
)

// HeapValue 堆中元素
type HeapValue struct {
	// 该值用来比较的
	Value int64
	// 实际的值
	Key string
	// 额外数据
	Extra interface{}
	// 数组索引
	Index int
}

// Heap 最小堆，最小的值永远在树根
type Heap struct {
	// 堆的大小
	size int
	// 使用内部的数组来模拟树
	// 一个节点下标为 i，那么父亲节点的下标为 (i-1)/2
	// 一个节点下标为 i，那么左儿子的下标为 2i+1，右儿子下标为 2i+2
	array []*HeapValue
	lock  sync.Mutex
}

// NewMinHeap 初始化一个最小堆
func NewMinHeap(array []*HeapValue) *Heap {
	if len(array) != 0 {
		return heapFast(array)
	}

	h := new(Heap)
	h.array = make([]*HeapValue, 0, heapInitCap)
	return h
}

// Push 最小堆插入元素
func (h *Heap) Push(x *HeapValue) {
	if h == nil {
		panic("h nil")
	}
	if x == nil {
		panic("x nil")
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	// 往尾巴放一个
	h.array = append(h.array, nil)

	// 堆没有元素时，使元素成为顶点后退出
	if h.size == 0 {
		x.Index = 0
		h.array[0] = x
		h.size++
		return
	}

	// i 是要插入节点的下标
	i := h.size

	// 如果下标存在
	// 将小的值 x 一直上浮
	for i > 0 {
		// parent为该元素父亲节点的下标
		parent := (i - 1) / 2

		// 如果插入的值大于父亲节点，那么可以直接退出循环，因为父亲仍然是最小的
		if x.Value > h.array[parent].Value {
			break
		}

		// 否则将父亲节点与该节点互换，然后向上翻转，将最小的元素一直往上推
		h.array[i] = h.array[parent]
		h.array[i].Index = i
		i = parent
	}

	// 将该值 x 放在不会再翻转的位置
	h.array[i] = x
	h.array[i].Index = i

	// 堆数量加一
	h.size++
}

// Pop 最小堆移除根节点元素，也就是最小的元素
func (h *Heap) Pop() *HeapValue {
	return h.PopIndex(0)
}

// PopIndex 最小堆移除某个元素
func (h *Heap) PopIndex(index int) *HeapValue {
	if h == nil {
		panic("h nil")
	}

	if index >= h.size {
		return nil
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	// 没有元素，返回-1
	if h.size == 0 {
		return nil
	}

	// 取出根节点
	ret := h.array[index]

	// 因为根节点要被删除了，将最后一个节点放到根节点的位置上
	h.size--
	x := h.array[h.size]  // 将最后一个元素的值先拿出来
	h.array[h.size] = ret // 将移除的元素放在最后一个元素的位置上

	// 对根节点进行向下翻转，大的值 x 一直下沉，维持最小堆的特征
	i := index
	for {
		// a，b为下标 i 左右两个子节点的下标
		a := 2*i + 1
		b := 2*i + 2

		// 左儿子下标超出了，表示没有左子树，那么右子树也没有，直接返回
		if a >= h.size {
			break
		}

		// 有右子树，拿到两个子节点中较大小节点的下标
		if b < h.size && h.array[b].Value < h.array[a].Value {
			a = b
		}

		// 父亲节点的值都小于两个儿子较大的那个，不需要向下继续翻转了，返回
		if x.Value < h.array[a].Value {
			break
		}

		// 将较小的儿子与父亲交换，维持这个最小堆的特征
		h.array[i] = h.array[a]
		h.array[i].Index = i

		// 继续往下操作
		i = a
	}

	// 将最后一个元素的值 x 放在不会再翻转的位置
	h.array[i] = x
	h.array[i].Index = i

	// 清除尾巴
	h.array = h.array[:h.size]

	// 垃圾回收
	if cap(h.array) > heapCleanCap && cap(h.array) > 2*h.size {
		before := h.array

		newCap := heapInitCap
		if h.size > heapInitCap {
			newCap = h.size
		}

		h.array = make([]*HeapValue, h.size, newCap)

		if h.size != 0 {
			copy(h.array, before[0:h.size])
		}
	}
	return ret
}

// Min 最小堆获取最小值
func (h *Heap) Min() *HeapValue {
	if h == nil {
		panic("h nil")
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.size == 0 {
		return nil
	}

	return h.array[0]
}

// Get 获取某个元素
func (h *Heap) Get(index int) *HeapValue {
	if h == nil {
		panic("h nil")
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.size == 0 {
		return nil
	}

	if index >= h.size {
		return nil
	}

	return h.array[index]
}

func (h *Heap) Size() int {
	if h == nil {
		panic("h nil")
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	return h.size
}

func heapFast(array []*HeapValue) *Heap {
	// 堆的元素数量
	count := len(array)

	// 最底层的叶子节点下标，该节点位置不定，但是该叶子节点右边的节点都是叶子节点
	start := count/2 + 1

	// 从最底层开始，逐一对节点进行下沉
	for start >= 0 {
		sift(array, start, count)
		start-- // 表示左偏移一个节点，如果该层没有节点了，那么表示到了上一层的最右边
	}

	for k, v := range array {
		v.Index = k
	}

	h := new(Heap)
	h.array = array
	h.size = count
	return h
}

// 下沉操作，需要下沉的元素是 array[start]，参数 count 只要用来判断是否到底堆底，使得下沉结束
func sift(array []*HeapValue, start, count int) {
	// 父亲节点
	root := start

	// 左儿子
	child := root*2 + 1

	// 如果有下一代
	for child < count {
		// 右儿子比左儿子小，那么要翻转的儿子改为右儿子
		if count-child > 1 && array[child].Value > array[child+1].Value {
			child++
		}

		// 父亲节点比儿子大，那么将父亲和儿子位置交换
		if array[root].Value > array[child].Value {
			array[root], array[child] = array[child], array[root]
			// 继续往下沉
			root = child
			child = root*2 + 1
		} else {
			return
		}
	}
}
