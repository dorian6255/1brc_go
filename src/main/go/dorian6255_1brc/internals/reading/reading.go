package reading

import (
	"bufio"
	"dorian6255_1brc/internals/dataType"
	"strconv"
	"strings"
	"sync"
)

const endlineSymbol byte = '\n'

func SendScanBuf(data []byte, lines chan [dataType.BufferSize]byte) {

	// read_lines := strings.Split(string(data), "\n")
	// read_lines = read_lines[:len(read_lines)-1]
	var buffer [dataType.BufferSize]byte
	for i := 0; i < len(data); i++ {
		//form line
		j := 0
		for ; i+j < len(data) && data[i+j] != endlineSymbol; j++ {
			buffer[j] = data[i+j]
		}

		i += j

		lines <- buffer
	}

	dataType.Wgread.Done()
	return

}

func ProcessLine(lines chan [dataType.BufferSize]byte, in chan dataType.ValueIn) {

	var group sync.WaitGroup
	batch := make([][dataType.BufferSize]byte, dataType.BatchSize)
	i := 0
	for newLine := range lines {

		batch[i] = newLine

		i++
		if i == dataType.BatchSize {
			batchCopy := batch
			group.Add(1)
			go func() {
				for y := 0; y < i; y++ {

					splitLine := strings.Split(string(batchCopy[y][:]), ";")

					name := splitLine[0]

					v, _ := strconv.ParseFloat(splitLine[1], 32)
					value := float32(v)
					in <- dataType.ValueIn{Name: name, Value: &value}

				}
				group.Done()
				return
			}()
			i = 0

		}

	}

	group.Wait()
	dataType.Wgprocess.Done()
	return
}

// not utilised
func SendScan(scanner *bufio.Scanner, lines chan string) {

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

	dataType.Wgread.Done()
	return

}
