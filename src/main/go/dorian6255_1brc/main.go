package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"strconv"
	"strings"
	"sync"
	// "text/scanner"
)

var wgread sync.WaitGroup

var wgprocess sync.WaitGroup

var wgstore sync.WaitGroup

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
	Min float32
	Max float32
	Avg float32
}

type valueIn struct {
	name  string
	value float32
}

func main() {
	minimum := sync.Map{}
	maximum := sync.Map{}
	average := sync.Map{}

	nb := sync.Map{}
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
	content, err := os.ReadFile(filename)
	lines := make(chan string, 1000)
	in := make(chan valueIn, 1000)
	if err != nil {
		panic(err)

	}
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Split(bufio.ScanLines)

	for i := 0; i < 1; i++ {
		go storeElement(in, &minimum, &maximum, &average, &nb)
		wgstore.Add(1)
	}

	go sendScan(*scanner, lines)

	wgread.Add(1)
	for i := 0; i < 1; i++ {

		go processLine(lines, in)

		wgprocess.Add(1)
	}

	wgread.Wait()
	close(lines)

	wgprocess.Wait()

	close(in)

	wgstore.Wait()

	res := make(map[string]valuesOut)
	maximum.Range(func(key, value any) bool {
		k, _ := key.(string)
		v := value.(float32)
		res[k] = valuesOut{Max: v}
		return true
	})
	minimum.Range(func(key, value any) bool {
		k, _ := key.(string)
		v := value.(float32)
		tmp := res[k]
		tmp.Min = v
		res[k] = tmp
		return true
	})
	average.Range(func(key, value any) bool {
		k, _ := key.(string)
		v := value.(float32)
		tmp := res[k]
		tmp.Avg = v
		res[k] = tmp
		return true
	})

	fmt.Println(res)

}
func sendScan(scanner bufio.Scanner, lines chan string) {

	for scanner.Scan() {
		lines <- scanner.Text()
	}

	wgread.Done()

}

func storeElement(in chan valueIn, maximum, minimum, average, number *sync.Map) {

	for elem := range in {

		min, ok := minimum.LoadOrStore(elem.name, elem.value)
		f, _ := min.(float32)

		if ok && elem.value < float32(f) {

			minimum.Store(elem.name, elem.value)
		}

		max, ok := maximum.LoadOrStore(elem.name, elem.value)

		f, _ = max.(float32)
		if ok && elem.value > float32(f) {
			maximum.Store(elem.name, elem.value)

		}
		avg, ok := average.LoadOrStore(elem.name, elem.value)

		f, _ = avg.(float32)
		if ok {
			nb, _ := number.Load(elem.name)
			n := float32(nb.(float64))

			prevAvg, _ := average.Load(elem.name)
			prev := prevAvg.(float32)
			value := ((prev * n) + elem.value) / (n + 1)
			average.Store(elem.name, value)
		} else {
			number.Store(elem.name, 1.0)
		}

	}

	wgstore.Done()

}
func processLine(lines chan string, in chan valueIn) {
	for newLine := range lines {

		splitLine := strings.Split(newLine, ";")
		name := splitLine[0]
		value, _ := strconv.ParseFloat(splitLine[1], 32)
		in <- valueIn{name: name, value: float32(value)}
	}

	wgprocess.Done()
}
