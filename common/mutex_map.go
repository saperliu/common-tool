package common

import (
	"fmt"
	"io"
	"sync"
)

type SafeMutexMap struct {
	m map[interface{}] interface{}
	sync.RWMutex
}

func New() *SafeMutexMap {
	return &SafeMutexMap{
		m: make(map[interface{}] interface{}),
	}
}

func (s *SafeMutexMap) Dump(w io.Writer) {
	s.Lock()
	defer s.Unlock()

	keys := make([]interface{}, len(s.m))

	i := 0
	for k, _ := range s.m {
		keys[i] = k
		i = i + 1
	}
	//sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s = %s\n", k, s.m[k])
	}
}

func (s *SafeMutexMap) Get(key interface{}) interface{} {
	s.Lock()
	defer s.Unlock()
	return s.m[key]
}

func (s *SafeMutexMap) Set(key interface{}, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}

func (s *SafeMutexMap) Remove(key interface{}) {
	s.Lock()
	defer s.Unlock()
	delete(s.m,key)
}

func (s *SafeMutexMap) GetValue(key interface{}) (interface{}, bool) {
	s.Lock()
	defer s.Unlock()
	value, ok := s.m[key]
	return value, ok
}
