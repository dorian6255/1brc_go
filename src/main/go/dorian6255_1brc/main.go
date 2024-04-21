package main

import (
	"dorian6255_1brc/internals/dataType"
	"dorian6255_1brc/internals/handling"
	"dorian6255_1brc/internals/reading"
	"os"
)

// optimisations idea :
// string build to return the output in 1 time
// either pass float value as pointer or create my own type of float to reduce size because value are between -99.00 && 99.99 in worst case scenario to reduce copy time
// remove Name from values out
// store []byte instead of string ?
//  pass struct as reference

func Handle(filename string) {
	minimum := make(map[string]float32, dataType.UniqueStationName)
	maximum := make(map[string]float32, dataType.UniqueStationName)
	average := make(map[string]float32, dataType.UniqueStationName)
	nb := make(map[string]float32, dataType.UniqueStationName)
	// maximum := sync.Map{}
	// average := sync.Map{}
	// nb := sync.Map{}
	//
	//start
	// var res = sync.Map{}

	content, _ := os.ReadFile(filename)
	lines := make(chan [dataType.BufferSize]byte, dataType.BatchSize)
	in := make(chan dataType.ValueIn, dataType.BatchSize)

	max := make(chan dataType.StoreUnit, dataType.BatchSize)
	min := make(chan dataType.StoreUnit, dataType.BatchSize)
	avg := make(chan dataType.StoreUnit, dataType.BatchSize)

	go reading.SendScanBuf(content, lines)

	dataType.Wgread.Add(1)

	for i := 0; i < 1; i++ {

		go reading.ProcessLine(lines, in)

		dataType.Wgprocess.Add(1)
	}

	go handling.HandleAvg(avg, average, nb)

	dataType.Wgstore.Add(1)
	go handling.HandleMax(max, maximum)

	dataType.Wgstore.Add(1)
	go handling.HandleMin(min, minimum)

	dataType.Wgstore.Add(1)

	go handling.HandleValues(in, min, avg, max)
	dataType.Wgstore.Add(1)

	dataType.Wgread.Wait()
	close(lines)

	dataType.Wgprocess.Wait()

	close(in)

	dataType.Wgstore.Wait()

	resmap := make(map[string]dataType.ValuesOut, dataType.UniqueStationName)

	for key, value := range maximum {
		resmap[key] = dataType.ValuesOut{Name: key, Max: &value}
	}
	// maximum.Range(func(key, value any) bool {
	// 	k, _ := key.(string)
	// 	v := value.(float32)
	// 	return true
	// })
	//
	for key, value := range minimum {

		tmp := resmap[key]
		tmp.Min = &value
		resmap[key] = tmp
	}
	for key, value := range average {

		tmp := resmap[key]
		tmp.Avg = &value

		resmap[key] = tmp
	}
	// minimum.Range(func(key, value any) bool {
	//
	// 	k, _ := key.(string)
	// 	v := value.(float32)
	// 	tmp := resmap[k]
	// 	tmp.Min = v
	// 	resmap[k] = tmp
	//
	// 	return true
	// })
	// average.Range(func(key, value any) bool {
	// 	k, _ := key.(string)
	// 	v := value.(float32)
	// 	tmp := resmap[k]
	// 	tmp.Avg = v
	// 	resmap[k] = tmp
	// 	return true
	// })

	// fmt.Print("{")
	// res := make([]string, len(resmap))
	// idx := 0
	// for k := range resmap {
	// 	res[idx] = k
	//
	// 	idx++
	// }
	// sort.Strings(res)
	// i := false
	// for _, str := range res {
	// 	value := resmap[str]
	// 	if !i {
	//
	// 		fmt.Printf("%v=%.2f/%.2f/%.2f", value.Name, *value.Min, *value.Avg, *value.Max)
	// 		i = true
	// 	} else {
	//
	// 		fmt.Printf(", %v=%.2f/%.2f/%.2f", value.Name, *value.Min, *value.Avg, *value.Max)
	// 	}
	// }
	// fmt.Print("}")

}
func main() {

	filename := os.Args[1]
	Handle(filename)
}
