package main

import (
	"fmt"
	"os"
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
	name string
	min  int
	max  int
	avg  int
}

// INFO: i did not unclude test because of how simple the function is, and the lack of opporunity i can see to improve the efficiency
func loadFile(filename string) []byte {

	file, err := os.ReadFile(filename)
	if err != nil {
		panic("wrong filename")
	}
	return file

}

// TODO:
func splitContent(content *[]byte) []*[]byte {

	return nil
}

// TODO:
func mergeResult(data ...map[string]outputType) map[string]outputType {

	return nil
}

// TODO:
func process(data []byte) map[string]outputType {

	var tmp [maxSizeName]byte
	for i := 0; i < maxSizeName; i++ {
		tmp[i] = data[i]

	}
	fmt.Println(tmp)
	fmt.Println(string(tmp[:]))
	var b byte = '\n'
	var c byte = ';'
	fmt.Println(b)
	fmt.Println(c)
	return nil
}

// TODO: for now it take ~34 ns per operation, maybe we can improve some stuf : reduce size int?, byte size ? paralell ?
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
