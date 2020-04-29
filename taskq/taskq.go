package taskq

import (
	"errors"
	"sync"
)

type request struct {
	key      string
	method   func() (interface{}, error)
	callback chan interface{}
}

type reply struct {
	key    string
	result interface{}
}

type TaskQueue struct {
	callbackDic map[string][]chan interface{}
	inputQueue  chan *request
	outputQueue chan *reply
	shutdown    chan bool
	wg          sync.WaitGroup
	isDestory   bool
	destoryOnce sync.Once
}

func NewTaskQueeu() *TaskQueue {
	m := &TaskQueue{
		callbackDic: make(map[string][]chan interface{}),
		inputQueue:  make(chan *request, 16),
		outputQueue: make(chan *reply, 4),
		shutdown:    make(chan bool),
		isDestory:   false,
	}
	go m.runloop()
	return m
}

func (m *TaskQueue) runloop() {
	for {
		select {
		case <-m.shutdown:
			{
				return
			}
		case rep := <-m.outputQueue:
			{
				target, ok := m.callbackDic[rep.key]
				if ok {
					for _, callback := range target {
						callback <- rep.result
					}
					delete(m.callbackDic, rep.key)
				}
			}
		case req := <-m.inputQueue:
			{
				target, ok := m.callbackDic[req.key]
				if ok {
					m.callbackDic[req.key] = append(target, req.callback)
				} else {
					target = make([]chan interface{}, 1)
					target[0] = req.callback
					m.callbackDic[req.key] = target

					go func(key string, method func() (interface{}, error)) {
						res, err := Lambda(method, 0)
						if err != nil {
							m.outputQueue <- &reply{key: key, result: err}
						} else {
							m.outputQueue <- &reply{key: key, result: res}
						}
					}(req.key, req.method)
				}
			}
		}
	}
}

func (m *TaskQueue) Destory() {
	m.destoryOnce.Do(func() {
		m.isDestory = true
	})
	m.wg.Wait()
	close(m.shutdown)
	close(m.inputQueue)
	close(m.outputQueue)
}

func (m *TaskQueue) Exec(key string, method func() (interface{}, error)) (interface{}, error) {
	if m.isDestory {
		return nil, errors.New("TaskQueue Destoried")
	}

	m.wg.Add(1)
	defer m.wg.Done()

	callback := make(chan interface{}, 1)
	m.inputQueue <- &request{key: key, method: method, callback: callback}

	res := <-callback
	close(callback)
	switch res.(type) {
	case error:
		return nil, res.(error)
	default:
		return res, nil
	}
}
