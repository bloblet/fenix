package databases

import (
	"bytes"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func MakeULID() ulid.ULID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

func MakeFlatULID(t time.Time) ulid.ULID {
	b := bytes.Buffer{}
	b.Write(make([]byte, 16))
	entropy := ulid.Monotonic(&b, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

func TimestampFromID(id ulid.ULID) time.Time {
	return ulid.Time(id.Time())
}
