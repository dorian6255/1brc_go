package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strconv"
	"strings"
	"sync"
	// "text/scanner"
)

var wgread sync.WaitGroup

var wgprocess sync.WaitGroup

var wgstore sync.WaitGroup

const batchSize = 1000

const uniqueStationName = 10000

// my plan is to divide the file by x goroutine  -3
// 1 goroutine will be allocated to find the min value of each one
// another one will be allocated to find the max value of each one
// another one will be allocated to find the avg value of each one
// i'll have to try if one goroutine can handle all 3 of it / if it's faster
// each one will be given a part of the file to read
// we store everything in ram, a map[NAME][]VALUES
// once we fully reach the end of the file, we must find the avg of each one
// then we return the result
type valuesOut struct {
	Name string
	Min  float32
	Max  float32
	Avg  float32
}
type storeUnit struct {
	name  string
	value float32
}

type valueIn struct {
	name  string
	value float32
}

func main() {
	minimum := make(map[string]float32, uniqueStationName)
	maximum := make(map[string]float32, uniqueStationName)
	average := make(map[string]float32, uniqueStationName)
	nb := make(map[string]float32, uniqueStationName)
	// maximum := sync.Map{}
	// average := sync.Map{}
	// nb := sync.Map{}
	//
	f, err := os.Create("profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()
	// Start tracing
	traceFile, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer traceFile.Close()

	if err := trace.Start(traceFile); err != nil {
		panic(err)
	}
	defer trace.Stop()
	//start
	// var res = sync.Map{}

	filename := os.Args[1]
	// content, err := os.ReadFile(filename)
	lines := make(chan string, 10000)
	in := make(chan valueIn)

	max := make(chan storeUnit, batchSize)
	min := make(chan storeUnit, batchSize)
	avg := make(chan storeUnit, batchSize)

	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// scanner := bufio.NewScanner(bytes.NewReader(content))
	// scanner.Split(bufio.ScanLines)

	for i := 0; i < 1; i++ {

		go sendScan(scanner, lines)

		wgread.Add(1)
	}

	for i := 0; i < 1; i++ {

		go processLine(lines, in)

		wgprocess.Add(1)
	}

	go handleAvg(avg, average, nb)

	wgstore.Add(1)
	go handleMax(max, maximum)

	wgstore.Add(1)
	go handleMin(min, minimum)

	wgstore.Add(1)

	go handleValues(in, min, avg, max)
	wgstore.Add(1)

	wgread.Wait()
	close(lines)

	wgprocess.Wait()

	close(in)

	wgstore.Wait()

	resmap := make(map[string]valuesOut, uniqueStationName)

	for key, value := range maximum {
		resmap[key] = valuesOut{Name: key, Max: value}
	}
	// maximum.Range(func(key, value any) bool {
	// 	k, _ := key.(string)
	// 	v := value.(float32)
	// 	return true
	// })
	//
	for key, value := range minimum {

		tmp := resmap[key]
		tmp.Min = value
		resmap[key] = tmp
	}
	for key, value := range average {

		tmp := resmap[key]
		tmp.Avg = value

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

	fmt.Print("{")
	res := make([]string, uniqueStationName)
	idx := 0
	for k := range resmap {
		res[idx] = k

		idx++
	}
	sort.Strings(res)
	i := false
	for _, str := range res {
		value := resmap[str]
		if !i {

			fmt.Printf("%v=%.2f/%.2f/%.2f", value.Name, value.Min, value.Avg, value.Max)
			i = true
		} else {

			fmt.Printf(", %v=%.2f/%.2f/%.2f", value.Name, value.Min, value.Avg, value.Max)
		}
	}
	fmt.Print("}")

}
func sendScan(scanner *bufio.Scanner, lines chan string) {

	// read_lines := strings.Split(string(*data), "\n")
	//
	// read_lines = read_lines[:len(read_lines)-1]
	//
	// for _, line := range read_lines {
	// 	lines <- line
	// }
	for scanner.Scan() {

		lines <- scanner.Text()
	}

	wgread.Done()
	return

}

func handleMax(in chan storeUnit, maximum map[string]float32) {

	for elem := range in {
		max, ok := maximum[elem.name]

		// f, _ := max.(float32)

		if !ok || elem.value > max {
			maximum[elem.name] = elem.value

		}
	}

	wgstore.Done()
	return
}
func handleAvg(in chan storeUnit, average, number map[string]float32) {

	for elem := range in {
		nb, ok := number[elem.name]

		if ok {

			prevAvg, _ := average[elem.name]

			value := ((prevAvg * nb) + elem.value) / (nb + 1)
			average[elem.name] = value
		} else {
			number[elem.name] = 1.0

			average[elem.name] = elem.value
		}
	}

	wgstore.Done()
	return
}
func handleMin(in chan storeUnit, minimum map[string]float32) {

	for elem := range in {
		// min, ok := minimum.LoadOrStore(elem.name, elem.value)

		min, ok := minimum[elem.name]

		if !ok || elem.value < min {
			minimum[elem.name] = elem.value

		}
	}

	wgstore.Done()
	return
}

func handleValues(in chan valueIn, minimum, average, maximum chan storeUnit) {

	for elem := range in {
		minimum <- storeUnit{name: elem.name, value: elem.value}
		maximum <- storeUnit{name: elem.name, value: elem.value}
		average <- storeUnit{name: elem.name, value: elem.value}
	}
	close(minimum)
	close(maximum)
	close(average)
	wgstore.Done()
}

func processLine(lines chan string, in chan valueIn) {

	var group sync.WaitGroup
	batch := make([]string, batchSize)
	i := 0
	for newLine := range lines {

		batch[i] = newLine

		i++
		if i == batchSize {
			group.Add(1)
			go func() {
				for y := 0; y < i; y++ {

					splitLine := strings.Split(batch[y], ";")
					name := splitLine[0]
					value, _ := strconv.ParseFloat(splitLine[1], 32)
					in <- valueIn{name: name, value: float32(value)}

				}
				group.Done()
				return
			}()
			i = 0
		}

	}

	group.Wait()
	wgprocess.Done()
	return
}
