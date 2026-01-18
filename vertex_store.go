package dag

import "sync"

// vertexStore centralizes vertex storage and per-vertex locks.
// DAG only manages graph topology and delegates vertex lifecycle to this store.
type vertexStore struct {
	mu     sync.RWMutex
	byHash map[interface{}]string
	byID   map[string]Vertex
	hashFn func(interface{}) interface{}
	locks  *dMutex
}

func newVertexStore(hashFn func(interface{}) interface{}) *vertexStore {
	if hashFn == nil {
		hashFn = defaultVertexHashFunc
	}
	return &vertexStore{
		byHash: make(map[interface{}]string),
		byID:   make(map[string]Vertex),
		hashFn: hashFn,
		locks:  newDMutex(),
	}
}

func (s *vertexStore) setHashFunc(hashFn func(interface{}) interface{}) {
	if hashFn == nil {
		hashFn = defaultVertexHashFunc
	}
	s.mu.Lock()
	s.hashFn = hashFn
	s.mu.Unlock()
}

func (s *vertexStore) hash(value interface{}) interface{} {
	s.mu.RLock()
	hashFn := s.hashFn
	s.mu.RUnlock()
	if hashFn == nil {
		return defaultVertexHashFunc(value)
	}
	return hashFn(value)
}

func (s *vertexStore) add(id string, value interface{}) error {
	if value == nil {
		return VertexNilError{}
	}
	hash := s.hash(value)
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byHash[hash]; exists {
		return VertexDuplicateError{value}
	}
	if _, exists := s.byID[id]; exists {
		return IDDuplicateError{id}
	}
	s.byID[id] = &vertex{id: id, val: value}
	s.byHash[hash] = id
	return nil
}

func (s *vertexStore) get(id string) (Vertex, bool) {
	s.mu.RLock()
	v, ok := s.byID[id]
	s.mu.RUnlock()
	return v, ok
}

func (s *vertexStore) value(id string) (interface{}, bool) {
	v, ok := s.get(id)
	if !ok {
		return nil, false
	}
	return v.Value(), true
}

func (s *vertexStore) hasID(id string) bool {
	s.mu.RLock()
	_, ok := s.byID[id]
	s.mu.RUnlock()
	return ok
}

func (s *vertexStore) hasHash(hash interface{}) bool {
	s.mu.RLock()
	_, ok := s.byHash[hash]
	s.mu.RUnlock()
	return ok
}

func (s *vertexStore) hashByID(id string) (interface{}, bool) {
	v, ok := s.get(id)
	if !ok {
		return nil, false
	}
	return s.hash(v.Value()), true
}

func (s *vertexStore) idByHash(hash interface{}) (string, bool) {
	s.mu.RLock()
	id, ok := s.byHash[hash]
	s.mu.RUnlock()
	return id, ok
}

func (s *vertexStore) delete(id string, hash interface{}) {
	s.mu.Lock()
	delete(s.byID, id)
	delete(s.byHash, hash)
	s.mu.Unlock()
}

func (s *vertexStore) count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.byID)
}

func (s *vertexStore) values() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]interface{}, len(s.byID))
	for id, v := range s.byID {
		out[id] = v.Value()
	}
	return out
}

func (s *vertexStore) eachByHash(fn func(hash interface{}, id string, value interface{})) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for hash, id := range s.byHash {
		v, ok := s.byID[id]
		if !ok {
			continue
		}
		fn(hash, id, v.Value())
	}
}

func (s *vertexStore) lock(hash interface{}) {
	s.locks.lock(hash)
}

func (s *vertexStore) unlock(hash interface{}) {
	s.locks.unlock(hash)
}

/***************************
********** dMutex **********
****************************/

type cMutex struct {
	mutex sync.Mutex
	count int
}

// Structure for dynamic mutexes.
type dMutex struct {
	mutexes     map[interface{}]*cMutex
	globalMutex sync.Mutex
}

// Initialize a new dynamic mutex structure.
func newDMutex() *dMutex {
	return &dMutex{
		mutexes: make(map[interface{}]*cMutex),
	}
}

// Get a lock for instance i
func (d *dMutex) lock(i interface{}) {

	// acquire global lock
	d.globalMutex.Lock()

	// if there is no cMutex for i, create it
	if _, ok := d.mutexes[i]; !ok {
		d.mutexes[i] = new(cMutex)
	}

	// increase the count in order to show, that we are interested in this
	// instance mutex (thus now one deletes it)
	d.mutexes[i].count++

	// remember the mutex for later
	mutex := &d.mutexes[i].mutex

	// as the cMutex is there, we have increased the count, and we know the
	// instance mutex, we can release the global lock
	d.globalMutex.Unlock()

	// and wait on the instance mutex
	(*mutex).Lock()
}

// Release the lock for instance i.
func (d *dMutex) unlock(i interface{}) {

	// acquire global lock
	d.globalMutex.Lock()

	// unlock instance mutex
	d.mutexes[i].mutex.Unlock()

	// decrease the count, as we are no longer interested in this instance
	// mutex
	d.mutexes[i].count--

	// if we are the last one interested in this instance mutex delete the
	// cMutex
	if d.mutexes[i].count == 0 {
		delete(d.mutexes, i)
	}

	// release the global lock
	d.globalMutex.Unlock()
}

