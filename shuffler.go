package parallelshuffler

import (
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"sync"
)

type ShufflerConfig struct {
	MaxSize          uint64
	ResultBufferSize uint64
}

type TaskContext struct {
	Id    uint64
	Value interface{}
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func byPriority(a, b interface{}) int {
	priorityA := a.(*TaskContext).Id
	priorityB := b.(*TaskContext).Id
	switch {
	case priorityA < priorityB:
		return -1
	case priorityA > priorityB:
		return 1
	default:
		return 0
	}
}

type Shuffler struct {
	lastSequence  uint64
	config        ShufflerConfig
	resultChan    chan TaskContext
	priorityQueue *pq.Queue
	mu            sync.Mutex
}

func NewShuffler(config ShufflerConfig) *Shuffler {
	s := &Shuffler{
		config: config,
	}
	s.priorityQueue = pq.NewWith(byPriority)
	s.resultChan = make(chan TaskContext, config.ResultBufferSize)
	s.lastSequence = 0
	return s
}

func (s *Shuffler) ResultChan() <-chan TaskContext {
	return s.resultChan
}

func (s *Shuffler) LastSequence() uint64 {
	return s.lastSequence
}

func (s *Shuffler) SetLastSequence(seq uint64) {
	s.lastSequence = seq
}

func (s *Shuffler) EnqueueTask(id uint64, obj interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.priorityQueue.Enqueue(&TaskContext{Id: id, Value: obj})
	//fmt.Println("EnqueueTask", id)
	// get head
	for {
		top, ok := s.priorityQueue.Peek()
		if !ok {
			break
		}
		v := top.(*TaskContext)
		//fmt.Println("Top", v.Id, s.lastSequence)
		if v.Id == s.lastSequence+1 {
			s.resultChan <- *v
			s.priorityQueue.Dequeue()
			//fmt.Println("SortedTask", v.Id)
			s.lastSequence++
		} else {
			break
		}
	}
}
