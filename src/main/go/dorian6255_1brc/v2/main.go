package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

const endlineSymbole byte = '\n'

const splitLineSymbol byte = ';'
const maxSizeName int = 35
const maxDifferentName int = 1000

func main() {
	filename := os.Args[1]
	data := loadFile(filename)
	res := process(data)
	showResult(res)

}

type outputType struct {
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

// receive runtime.numcp() map and merge the result of each map into a res map
func mergeResult(data ...map[string]outputType) map[string]outputType {
	//TODO: use make to see perrformance evolution
	//and try to use the first map
	var res = map[string]outputType{}

	for _, m := range data {

		for k, v := range m {
			value, ok := res[k]
			if ok {
				switch {
				//TODO: make a function that take two outputType, and return one with min,max,avg, updated
				case value.Min > v.Min:
					value.Min = v.Min
				case value.Max < v.Max:
					value.Max = v.Max

				}
				//NOTE: for avg and nb

				value.Avg = ((value.Avg * value.Nb) + (v.Avg * v.Nb)) / (value.Nb + v.Nb)
				value.Nb += v.Nb

				res[k] = value
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
func process(data []byte) map[string]outputType {

	res := make(map[string]outputType, maxDifferentName)

	//TODO: process line by line
	//for the size of the data
	// if we find an endline
	// we get the name, value of the line
	// we process the result (min, max, avg)
	// we go to the next one
	return res
}

// INFO:: for now it take ~34 ns per operation, maybe we can improve some stuf : reduce size int?, byte size ? paralell ?
func interpretLine(line []byte) ([]byte, int) {
	lineIdx := 0
	//handle name
	for ; lineIdx < len(line) && line[lineIdx] != splitLineSymbol; lineIdx++ {

	}

	return line[:lineIdx-1], interpretValue(line[lineIdx+1:])
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
func showResult(res map[string]outputType) {

}
