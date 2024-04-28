package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const endlineSymbole byte = '\n'

const splitLineSymbol byte = ';'
const maxSizeName int = 35
const maxSizeLine int = 50
const maxDifferentName int = 1000

func main() {
	filename := os.Args[1]
	data := loadFile(filename)
	content := splitContent(data)

	var mapSplit = make([]map[string]readingType, len(content))
	for i := 0; i < len(content); i++ {

		//TODO: use go function
		mapSplit[i] = processContent(content[i])
	}
	res := mergeResult(mapSplit...)
	showResult(res)

}

type readingType struct {
	// name string
	Min int
	Max int
	Avg int
	Nb  int
}

// INFO: i did not include test because of how simple the function is, and the lack of opporunity i can see to improve the efficiency
func loadFile(filename string) []byte {

	file, err := os.ReadFile(filename)
	if err != nil {
		errorString := fmt.Sprintf("invalide filename %s", filename)
		panic(errorString)
	}
	return file

}

// PERF: regarding the size of the function, i do not see how i can improve this one
// Used to split the file at the right place and not in the middle on an entry
// given a content, we return the idx of the latest endlinesymbole that we can find
func findNextEndline(content []byte) int {

	i := len(content) - 1
	for ; i > 0 && content[i] != endlineSymbole; i-- {

	}
	return i
}

// INFO: the res := make(..... improve the performance by 2 accordings to bench !!
// Given a file as content, we want this function to return X arrays of byte
// X is gonna be the number of thread on the cpu
// the goal here is to process the file multiple part at a time
func splitContent(content []byte) [][]byte {
	res := make([][]byte, runtime.NumCPU())
	splitSize := len(content) / runtime.NumCPU()

	idxStart := 0
	idxEnd := findNextEndline(content[0:splitSize])
	for i := 0; i < runtime.NumCPU(); i++ {

		res[i] = content[idxStart:idxEnd]
		idxStart = idxEnd + 1
		newEnd := idxEnd + splitSize
		if newEnd > len(content) {
			newEnd = len(content)
		}
		idxEnd = idxEnd + findNextEndline(content[idxEnd:newEnd])

	}

	return res
}

// INFO: tested indirectly via mergeResult function tests
// TODO: Implement Bench && try with min and max function
func mergeTwoOuputType(f1, f2 readingType) readingType {

	switch {
	//TODO: make a function that take two outputType, and return one with min,max,avg, updated
	case f1.Min > f2.Min:
		f1.Min = f2.Min
	case f1.Max < f2.Max:
		f1.Max = f2.Max

	}

	f1.Avg = ((f1.Avg * f1.Nb) + (f2.Avg * f2.Nb)) / (f1.Nb + f2.Nb)
	f1.Nb += f2.Nb

	return f1

}

// receive runtime.numcp() map and merge the result of each map into a res map
func mergeResult(data ...map[string]readingType) map[string]readingType {
	//TODO: use make to see perrformance evolution
	//and try to use the first map
	var res = map[string]readingType{}

	for _, m := range data {

		for k, v := range m {
			value, ok := res[k]
			if ok {

				res[k] = mergeTwoOuputType(value, v)
			} else {
				res[k] = v
			}

		}
	}

	return res
}

// TODO:
// process take a part of the file in params
// parse each line
// fill a map
// return the map
// it will run in a go routine
func processContent(data []byte) map[string]readingType {

	res := make(map[string]readingType, maxDifferentName)

	var lineBuffer = [maxSizeLine]byte{}
	bufferIdx := 0
	for dataIdx := 0; dataIdx < len(data); dataIdx++ {
		lineBuffer[bufferIdx] = data[dataIdx]

		if lineBuffer[bufferIdx] == endlineSymbole {

			name, value := interpretLine(lineBuffer[:])
			nameS := string(name)
			var valueOutputType = readingType{value, value, value, 1}
			v, ok := res[nameS]
			if ok {
				res[nameS] = mergeTwoOuputType(v, valueOutputType)

			} else {
				res[nameS] = readingType{value, value, value, 1}
			}

			bufferIdx = 0
		} else {
			bufferIdx++
		}

	}

	return res
}

// INFO:: for now it take ~34 ns per operation, maybe we can improve some stuf : reduce size int?, byte size ? paralell ?
func interpretLine(line []byte) ([]byte, int) {
	lineIdx := 0
	//handle name
	for ; lineIdx < len(line) && line[lineIdx] != splitLineSymbol; lineIdx++ {

	}

	return line[:lineIdx], interpretValue(line[lineIdx+1:])
}

// PERF: for now, it takes only ~0,00004ns per operation (unless the bench is wrong)
func interpretValue(line []byte) int {
	var res int
	var neg bool = false

	// value format is sign ??.??
	var value [4]byte

	//we store the sign here
	neg = line[0] == '-'
	var lineIdx = 0
	var valueIdx = 0
	if neg {
		lineIdx = 1 // we skip the "-"
		for ; lineIdx < len(line) && line[lineIdx] != endlineSymbole; lineIdx++ {
			if line[lineIdx] != '.' {

				value[valueIdx] = line[lineIdx]
				valueIdx++
			}
		}

		tmp, _ := strconv.ParseInt(string(value[:valueIdx]), int(endlineSymbole), 0)
		res = int(tmp)
		return -res
	}
	lineIdx = 0
	for ; lineIdx < len(line) && line[lineIdx] != 10; lineIdx++ {
		if line[lineIdx] != '.' {

			value[valueIdx] = line[lineIdx]
			valueIdx++
		}
	}

	tmp, _ := strconv.ParseInt(string(value[:valueIdx]), int(endlineSymbole), 0)
	res = int(tmp)
	return res
}

func showResult(res map[string]readingType) {
	var tmpStringArray = make([]string, len(res))
	i := 0
	for k := range res {

		tmpStringArray[i] = k
		i++
	}
	sort.Strings(tmpStringArray)
	var sb strings.Builder
	sb.WriteString("{")
	for i := 0; i < len(tmpStringArray); i++ {

		tmp := res[tmpStringArray[i]]
		min := float32(tmp.Min) / 10
		max := float32(tmp.Max) / 10
		avg := float32(tmp.Avg) / 10
		sb.WriteString(fmt.Sprintf("%v=%v/%v/%v, ", tmpStringArray[i], min, avg, max))

	}
	fmt.Println(sb.String())

}
