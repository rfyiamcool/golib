package timedmap

import (
	"time"
)

type Section interface {
	Ident() int
	Set(key, value interface{}, expiresAfter time.Duration, cb ...callback)
	GetValue(key interface{}) interface{}
	GetExpires(key interface{}) (time.Time, error)
	SetExpires(key interface{}, d time.Duration) error
	Contains(key interface{}) bool
	Remove(key interface{})
	Refresh(key interface{}, d time.Duration) error
	Flush()
	Size() (i int)
}

type section struct {
	tm  *TimedMap
	sec int
}

func newSection(tm *TimedMap, sec int) *section {
	return &section{
		tm:  tm,
		sec: sec,
	}
}

func (s *section) Ident() int {
	return s.sec
}

func (s *section) Set(key, value interface{}, expiresAfter time.Duration, cb ...callback) {
	s.tm.set(key, s.sec, value, expiresAfter, cb...)
}

func (s *section) GetValue(key interface{}) interface{} {
	v := s.tm.get(key, s.sec)
	if v == nil {
		return nil
	}
	return v.value
}

func (s *section) GetExpires(key interface{}) (time.Time, error) {
	v := s.tm.get(key, s.sec)
	if v == nil {
		return time.Time{}, ErrKeyNotFound
	}
	return v.expires, nil
}

func (s *section) SetExpires(key interface{}, d time.Duration) error {
	return s.tm.setExpire(key, s.sec, d)
}

func (s *section) Contains(key interface{}) bool {
	return s.tm.get(key, s.sec) != nil
}

func (s *section) Remove(key interface{}) {
	s.tm.remove(key, s.sec)
}

func (s *section) Refresh(key interface{}, d time.Duration) error {
	return s.tm.refresh(key, s.sec, d)
}

func (s *section) Flush() {
	for k := range s.tm.container {
		if k.sec == s.sec {
			s.tm.remove(k.key, k.sec)
		}
	}
}

func (s *section) Size() (i int) {
	for k := range s.tm.container {
		if k.sec == s.sec {
			i++
		}
	}
	return
}
