package entity

import (
	"encoding/binary"
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
	unique = make(map[uint64]struct{})
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

func NewResult(number uint64) *Result {
	mu.Lock()
	_, ok := unique[number]
	if ok {
		mu.Unlock()
		if len(unique) == MaxCount {
			return nil
		}

		once.Do(func() {
			rand.Seed(time.Now().UnixNano())
		})

		number = number - (number % MaxCount) + uint64(rand.Intn(MaxCount))
		return NewResult(number)
	}
	unique[number] = struct{}{}
	mu.Unlock()

	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, number)
	h := sha3.New256()
	h.Write(bytes)
	bs := h.Sum(nil)
	return &Result{Number: strconv.FormatUint(number, 10), Hash: bs}
}
