package parallelshuffler

import (
	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"sync"
	"time"
)

type ShufflerConfig struct {
	MaxSize          uint64
	ResultBufferSize uint64
	AllowSkip        bool
	SkipTimeout      time.Duration
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
	LastSequence  uint64
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
	return s
}

func (s *Shuffler) ResultChan() <-chan TaskContext {
	return s.resultChan
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
		//fmt.Println("Top", v.Id, s.LastSequence)
		if v.Id == s.LastSequence+1 {
			s.resultChan <- *v
			s.priorityQueue.Dequeue()
			//fmt.Println("SortedTask", v.Id)
			s.LastSequence++
		} else {
			break
		}
	}
}
