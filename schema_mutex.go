package JSONParse

import (
)

type SchemaMutex struct {
	mutexes		map[string]bool
}

func NewSchemaMutex() *SchemaMutex {
	mutex := new(SchemaMutex)
	mutex.mutexes = make(map[string]bool)

	return mutex
}

func (sm *SchemaMutex) Find(mutex_id string) (bool, bool) {
	val, found := sm.mutexes[mutex_id]
	if found {
		return val, found
	} else {
		return false, false
	}
}

func (sm *SchemaMutex) Add(mutex_id string) {
	sm.mutexes[mutex_id] = true
}

func (sm *SchemaMutex) Set(mutex_id string, value bool) {
	sm.mutexes[mutex_id] = value
}

func (sm *SchemaMutex) Remove(mutex_id string) {
	delete(sm.mutexes, mutex_id)
}
