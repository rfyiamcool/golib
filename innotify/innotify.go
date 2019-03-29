package innotify

import (
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidHandler      = errors.New("ErrInvalidHandler")
	ErrMessageChanOverflow = errors.New("ErrMessageChanOverflow")
	ErrWaitInnotifyTimeout = errors.New("ErrWaitInnotifyTimeout")
)

// null logger
var defualtLogger = func(level, s string) {}

type loggerType func(level, s string)

func SetLogger(logger loggerType) {
	defualtLogger = logger
}

type NCKey uint64

func NCUint64Key(i uint64) NCKey {
	return NCKey(i)
}

type InnotifyHandler struct {
	nc    *InnotifyCenter
	count int32 // Message counter befor remove listening, negative number indicate unlimit

	Key NCKey
	Ch  chan interface{}
}

func newInnotifyHandler(nc *InnotifyCenter, key NCKey, count int) *InnotifyHandler {
	return &InnotifyHandler{
		nc:    nc,
		count: int32(count),
		Key:   key,
		Ch:    make(chan interface{}, 1),
	}
}

type InnotifyCenter struct {
	cache map[NCKey][]chan interface{}
	lock  *sync.RWMutex
}

func NewInnotifyCenter() *InnotifyCenter {
	n := &InnotifyCenter{
		cache: make(map[NCKey][]chan interface{}),
		lock:  &sync.RWMutex{},
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)

			n.lock.Lock()
			knum := 0
			vnum := 0
			for _, v := range n.cache {
				knum++
				vnum += len(v)
			}
			n.lock.Unlock()
		}
	}()

	return n
}

var instance = NewInnotifyCenter()

func DefaultCenter() *InnotifyCenter {
	return instance
}

func (nc *InnotifyCenter) Remove(nh *InnotifyHandler) (rm bool) {
	nc.lock.Lock()
	if _, exist := nc.cache[nh.Key]; exist == true {
		delete(nc.cache, nh.Key)
		rm = true
	}
	nc.lock.Unlock()

	return
}

func (nc *InnotifyCenter) Register(key NCKey) *InnotifyHandler {
	return nc.register(key, -1)
}

// Remove handler after receive one Innotify
func (nc *InnotifyCenter) RegisterOnce(key NCKey) *InnotifyHandler {
	return nc.register(key, 1)
}

func (nc *InnotifyCenter) register(key NCKey, count int) *InnotifyHandler {
	nc.lock.Lock()

	nh := newInnotifyHandler(nc, key, count)

	if nhs, exist := nc.cache[key]; exist == false {
		nhs = []chan interface{}{nh.Ch}
		nc.cache[key] = nhs
	} else {
		nc.cache[key] = append(nhs, nh.Ch)
	}

	nc.lock.Unlock()

	return nh
}

func (nc *InnotifyCenter) Notify(key NCKey, mesg interface{}) (bool, error) {
	nc.lock.RLock()

	chs, exist := nc.cache[key]
	if exist == false {
		nc.lock.RUnlock()
		return false, nil
	}

	var err error

	for _, ch := range chs {
		select {
		case ch <- mesg:
		default:
			err = ErrMessageChanOverflow
		}
	}

	nc.lock.RUnlock()

	return true, err
}

func (nh *InnotifyHandler) Remove() bool {
	return nh.nc.Remove(nh)
}

func (nh *InnotifyHandler) Wait(timeout time.Duration) (mesg interface{}, err error) {
	if nh.count == 0 {
		err = ErrInvalidHandler
		return
	}

	timer := globalTimerPool.Get(timeout)
	defer globalTimerPool.Put(timer)

	select {
	case mesg = <-nh.Ch:
	case <-timer.C:
		err = ErrWaitInnotifyTimeout
	}

	if nc := atomic.AddInt32(&nh.count, -1); nc == 0 {
		// Remove
		nh.Remove()
	}

	return
}

func MakeId() uint64 {
	rand.Seed(time.Now().UnixNano())
	return uint64(rand.Int63())
}
