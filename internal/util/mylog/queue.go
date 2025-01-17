// 定义一个协程安全的队列结构，用于全局控制日志
package mylog

import "sync"

// 日志结构
type logItem struct {
	string
}

// 存放日志的队列
type logQueue struct {
	items []*logItem
}

var mu sync.RWMutex

func (li *logItem) String() string {
	return li.string
}

// 获取队列的长度
func (lq *logQueue) Len() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(lq.items)
}

// 从队头中取出一个元素
func (lq *logQueue) pollFirst() *logItem {
	if lq.Len() == 0 {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	item := lq.items[0]
	lq.items = lq.items[1:]
	return item
}

// 往队列的末尾添加一个元素
func (lq *logQueue) offerLast(item *logItem) {
	mu.Lock()
	defer mu.Unlock()
	lq.items = append(lq.items, item)
}

// 往队列的末尾批量添加元素
func (lq *logQueue) offerLastAll(items ...*logItem) {
	mu.Lock()
	defer mu.Unlock()
	lq.items = append(lq.items, items...)
}

// 返回队列最后一个元素, 不弹出
func (lq *logQueue) peekLast() *logItem {
	mu.RLock()
	defer mu.RUnlock()

	if lq.Len() == 0 {
		return nil
	}

	return lq.items[lq.Len()-1]
}
