package splash

import "sync"

func workerFunc(jobChan <-chan job, errorChan chan<- error, killChan <-chan interface{}, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for {
		select {
		case j := <-jobChan:
			if err := j.exec(); err != nil {
				errorChan <- err
			}
		case <-killChan:
			return
		}
	}
}

func newWorker(jobChan <-chan job, errorChan chan<- error, killChan <-chan interface{}, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)

	go workerFunc(jobChan, errorChan, killChan, waitGroup)
}
