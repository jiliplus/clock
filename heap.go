package clock

import (
	"container/heap"
	"time"
)

type taskManager interface {
	hasTaskToRun(now time.Time) bool
	pop() *task
	push(t *task)
	remove(t *task)
}

type task struct {
	deadline time.Time
	// 用于替代 fire，
	runFunc func(t *task) *task
	index   int
}

const removed = -1

func newTask(deadline time.Time, run func(t *task) *task) *task {
	return &task{
		deadline: deadline,
		runFunc:  run,
		index:    removed,
	}
}

func (t *task) run() *task {
	return t.runFunc(t)
}

func (t task) hasStopped() bool {
	return t.index == removed
}

type taskHeap []*task

func newTaskHeap() *taskHeap {
	t := make(taskHeap, 0, 1024)
	return &t
}

// *taskHeap 实现了 taskOrder 接口
func (h *taskHeap) push(t *task) {
	heap.Push(h, t)
}

func (h *taskHeap) pop() (t *task) {
	t, _ = heap.Pop(h).(*task)
	return
}

func (h taskHeap) hasExpiredTask(now time.Time) bool {
	return len(h) != 0 && !now.Before(h[0].deadline)
}

func (h taskHeap) hasTask() bool {
	return len(h) != 0
}

func (h *taskHeap) remove(t *task) {
	if !t.hasStopped() {
		heap.Remove(h, t.index)
	}
}

// *taskHeap 实现了 heap.Interface
func (h taskHeap) Len() int { return len(h) }

func (h taskHeap) Less(i, j int) bool {
	return h[i].deadline.Before(h[j].deadline)
}

func (h taskHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index, h[j].index = i, j
}

func (h *taskHeap) Push(x interface{}) {
	n := len(*h)
	t := x.(*task)
	t.index = n
	*h = append(*h, t)
}

func (h *taskHeap) Pop() interface{} {
	old := *h
	n := len(old)
	t := old[n-1]
	t.index = removed
	*h = old[0 : n-1]
	return t
}
