package dataType

import "sync"

var Wgread sync.WaitGroup
var Wgprocess sync.WaitGroup
var Wgstore sync.WaitGroup

const BatchSize = 1000
const NbReader = 10
const BufferSize = 35
const UniqueStationName = 10000

type ValuesOut struct {
	Name string
	Min  *float32
	Max  *float32
	Avg  *float32
}
type StoreUnit struct {
	Name  string
	Value *float32
}

type ValueIn struct {
	Name  string
	Value *float32
}
