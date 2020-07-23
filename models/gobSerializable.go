package models

import "bytes"

type GobSerializable interface {
	ToGob() bytes.Buffer
	FromGob(bytes.Buffer) interface{}
}
