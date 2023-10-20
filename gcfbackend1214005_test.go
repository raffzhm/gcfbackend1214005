package gcfbackend1214005

import (
	"fmt"
	"testing"
)

func TestGCHandlerFunc(t *testing.T) {
	data := GCHandlerFunc("string", "MONGO_URI", "pointmap")

	fmt.Printf("%+v", data)
}
