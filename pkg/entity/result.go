package entity

import (
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const MaxCount = 10000

var (
	once   sync.Once
	mu     sync.Mutex
	unique = make(map[int]struct{})
)

type Result struct {
	Number string `json:"number"`
	Hash   []byte `json:"hash"`
}

func (r *Result) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *Result) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}

func NewResult(number int) *Result {
	mu.Lock()
	_, ok := unique[number]
	if ok {
		if len(unique) == MaxCount {
			return nil
		}
		mu.Unlock()
		once.Do(func() {
			rand.Seed(time.Now().UnixNano())
		})
		number = number - (number % MaxCount) + rand.Intn(MaxCount)
		return NewResult(number)
	}
	unique[number] = struct{}{}
	mu.Unlock()
	s := strconv.Itoa(number)
	h := sha3.New256()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return &Result{Number: s, Hash: bs}
}
