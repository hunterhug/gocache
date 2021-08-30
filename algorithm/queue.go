package algorithm

import (
	"sync"
)

// use queue implement iterator
type linkQueue struct {
	root *linkNode  // 链表起点
	size int        // 队列的元素数量
	lock sync.Mutex // 为了并发安全使用的锁
}

// link node
type linkNode struct {
	next  *linkNode
	value bsTreeNode
}

// HasNext has next, queue size > 0
func (queue *linkQueue) HasNext() bool {
	if queue.size > 0 {
		return true
	}

	return false
}

func (queue *linkQueue) Next() (key string, value interface{}) {
	// 不断出队列
	element := queue.remove()

	// panic here
	if element == nil {
		panic("Next() empty")
	}

	// 左子树非空，入队列
	if element.leftOf() != nil {
		queue.add(element.leftOf())
	}

	// 右子树非空，入队列
	if element.rightOf() != nil {
		queue.add(element.rightOf())
	}

	//if element == nil {
	//	panic("Next() 2empty")
	//}
	return element.values()
}

// 入队
func (queue *linkQueue) add(v bsTreeNode) {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	// 如果栈顶为空，那么增加节点
	if queue.root == nil {
		queue.root = new(linkNode)
		queue.root.value = v
	} else {
		// 否则新元素插入链表的末尾
		// 新节点
		newNode := new(linkNode)
		newNode.value = v

		// 一直遍历到链表尾部
		nowNode := queue.root
		for nowNode.next != nil {
			nowNode = nowNode.next
		}

		// 新节点放在链表尾部
		nowNode.next = newNode
	}

	// 队中元素数量+1
	queue.size++
}

// 出队
func (queue *linkQueue) remove() bsTreeNode {
	queue.lock.Lock()
	defer queue.lock.Unlock()

	// 队中元素已空
	if queue.size == 0 {
		//panic("over limit")
		return nil
	}

	// 顶部元素要出队
	topNode := queue.root
	v := topNode.value

	// 将顶部元素的后继链接链上
	queue.root = topNode.next

	// 队中元素数量-1
	queue.size--

	return v
}