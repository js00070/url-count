package util

// QuickSort quicksort
func QuickSort(src []int, lessEq func(a, b int) bool) {
	if len(src) <= 8 { // 长度比较小时采用选择排序
		var temp int
		for i := 0; i < len(src); i++ {
			temp = i
			for j := i + 1; j < len(src); j++ {
				if lessEq(src[j], src[temp]) {
					temp = j
				}
			}
			src[i], src[temp] = src[temp], src[i]
		}
		return
	}
	i := 0
	j := len(src) - 1
	tmp := src[0]
	for i < j {
		for lessEq(tmp, src[j]) && i < j {
			j--
		}
		src[i] = src[j]
		for lessEq(src[i], tmp) && i < j {
			i++
		}
		src[j] = src[i]
	}
	src[i] = tmp
	QuickSort(src[0:i], lessEq)
	QuickSort(src[i+1:], lessEq)
}

type partitionResult struct {
	left  []int
	right []int
}

type partitionWorker struct {
	inputCh  chan []int
	outputCh chan partitionResult
	lessEq   func(a, b int) bool
}

func (worker *partitionWorker) init(length int, lessEq func(i, j int) bool) {
	worker.inputCh = make(chan []int, length/256+1)
	worker.outputCh = make(chan partitionResult, length/256+1)
	worker.lessEq = lessEq
}

func (worker *partitionWorker) start() {
	go worker.partitioning()
}

func (worker *partitionWorker) finish() {
	close(worker.inputCh)
}

func (worker *partitionWorker) partitioning() {
	for {
		src, ok := <-worker.inputCh
		if !ok {
			return
		}
		if len(src) <= 256 {
			QuickSort(src, worker.lessEq)
			worker.outputCh <- partitionResult{nil, nil}
			continue
		}
		i := 0
		j := len(src) - 1
		tmp := src[0]
		for i < j {
			for worker.lessEq(tmp, src[j]) && i < j {
				j--
			}
			src[i] = src[j]
			for worker.lessEq(src[i], tmp) && i < j {
				i++
			}
			src[j] = src[i]
		}
		src[i] = tmp
		worker.outputCh <- partitionResult{src[0:i], src[i+1:]}
	}
}

// ParallelSort parallel quicksort
func ParallelSort(src []int, lessEq func(a, b int) bool, workerNum int) {
	if workerNum == 1 {
		QuickSort(src, lessEq)
		return
	}
	workers := make([]partitionWorker, workerNum)
	for i := 0; i < len(workers); i++ {
		workers[i].init(len(src), lessEq)
		workers[i].start()
	}
	remaining := 1
	workers[0].inputCh <- src
	var res partitionResult
	for remaining != 0 {
		for i := 0; i < len(workers); i++ {
			select {
			case res = <-workers[i].outputCh:
				remaining--
				if res.left != nil {
					workers[i].inputCh <- res.left
					workers[(i+1)%workerNum].inputCh <- res.right
					remaining = remaining + 2
				}
			default:
			}
		}
	}
	for i := 0; i < len(workers); i++ {
		workers[i].finish()
	}
}
