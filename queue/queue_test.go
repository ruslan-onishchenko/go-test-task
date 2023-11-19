package queue

import (
	"gotesttask/types"
	"testing"
)

func TestMakeImpl(t *testing.T) {
    actual := MakeImpl(10,15)
    if actual.Storage == nil{
        t.Errorf("MakeImpl returns nil")
    }
}

func TestImplPut(t *testing.T) {
    impl := MakeImpl(10,15)
    if impl.Put("a", types.MsgBody{Message: "aaa"}) != nil{
        t.Errorf("MakeImpl returns nil")
    }
	
}