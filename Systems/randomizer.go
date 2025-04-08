package systems

import (
	"crypto/md5"
	"encoding/binary"
	"io"
	"math/rand"
	"time"
)

var randomizer_instance *Randomizer

type Randomizer struct {
	randomGenerator *rand.Rand
}

func GetRandomizer() *Randomizer {
	if randomizer_instance == nil {
		source := rand.NewSource(time.Now().UnixMicro())
		randomizer_instance = &Randomizer{
			randomGenerator: rand.New(source),
		}
	}

	return randomizer_instance
}

func (r *Randomizer) Seed(input string) {
	seed := r.convertStringToSeed(input)
	r.randomGenerator.Seed(seed)
}

func (r *Randomizer) convertStringToSeed(input string) int64 {
	hash := md5.New()
	io.WriteString(hash, input)
	seed := binary.BigEndian.Uint64(hash.Sum(nil))
	return int64(seed)
}

func (r *Randomizer) GetNumber(max int32) int32 {
	return r.randomGenerator.Int31n(max)
}

func (r *Randomizer) GetPercentage() float32 {
	return r.randomGenerator.Float32()
}
