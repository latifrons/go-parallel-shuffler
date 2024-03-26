package parallelshuffler

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func shuffleArray(arr []uint64) {
	// Use the Fisher-Yates shuffle algorithm
	for i := len(arr) - 1; i >= 1; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func TestEnqueue(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	shuffler := NewShuffler(ShufflerConfig{
		MaxSize:          100,
		ResultBufferSize: 10,
	})

	go func() {
		i := 1
		batch := 100
		for {
			array := make([]uint64, batch)
			for j := 0; j < batch; j++ {
				array[j] = uint64(i)
				i++
			}
			shuffleArray(array)
			for j := 0; j < batch; j++ {
				shuffler.EnqueueTask(array[j], array[j])
			}
			time.Sleep(30 * time.Millisecond)
		}
	}()

	go func() {
		c := shuffler.ResultChan()
		for {
			select {
			case task := <-c:
				fmt.Println("Dequeued", task.Id)
			}
		}
	}()

	select {}
}
