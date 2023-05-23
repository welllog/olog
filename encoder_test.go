package olog

import (
	"fmt"
	"testing"
)

func TestEncoder_EncodeJsonValue(t *testing.T) {
	e := ValueEncoder{&Buffer{}}
	e.EncodeJsonValue(map[string]string{"name": "lisi", "age": "18"})
	fmt.Println(string(e.Bytes()))
}
