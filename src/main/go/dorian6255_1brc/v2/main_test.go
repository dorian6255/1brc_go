package main

import (
	"fmt"
	"runtime"
	"testing"
)

func Test_findNextEndline(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "test simple case", args: args{[]byte{48, 49, endlineSymbole, 58, 47, 53}}, want: 2},
		{name: "test is latest", args: args{[]byte{48, 49, endlineSymbole}}, want: 2},
		{name: "test is first", args: args{[]byte{endlineSymbole, 48, 49}}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// showResult(tt.args.res)
			got := findNextEndline(tt.args.content)
			if got != tt.want {
				t.Errorf("Expected to find endline at %v, found at %v", tt.want, got)

			}
		})
	}

}

// func BenchmarkInterpretLine(b *testing.B) {
func BenchmarkSplitContent1M(b *testing.B) {

	type args struct {
		content []byte
	}
	type testParams struct {
		name string
		args args
	}
	testFiles := "tests/measurements1000000.txt"

	content := loadFile(testFiles)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		splitContent(content)

	}

}
func Test_splitContent(t *testing.T) {

	type args struct {
		content []byte
	}
	type testParams struct {
		name string
		args args
	}
	var testFolder string = "tests/"
	testFiles := [5]string{testFolder + "test24.txt", testFolder + "test48.txt", testFolder + "test96.txt", testFolder + "test100.txt", testFolder + "measurements1000000.txt"}

	var tests []testParams
	var content [][]byte
	for _, file := range testFiles {

		content = append(content, loadFile(file))
		tests = append(tests, testParams{name: file, args: args{loadFile(file)}})
	}

	//check if split in 12
	//check if last if \n
	const zero byte = 48
	const nine byte = 57
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := splitContent(tt.args.content)
			if len(got) != runtime.NumCPU() {
				t.Errorf("split not using all cpu, got: %v", len(got))
			}

			for _, split := range got {

				if split[len(split)-1] < zero || split[len(split)-1] > nine {
					t.Errorf("Last symbole is not a number as it should %s, %s", split[len(split)-1:], split[:])
				}
				// fmt.Println(split[len(split)-1:])

				// fmt.Println(string(zero))
			}

		})
	}
}
func BenchmarkMergeResult(b *testing.B) {
	type args struct {
		data []map[string]outputType
	}
	tests := []struct {
		name string
		args args
		want map[string]outputType
	}{
		{name: "testWithMoreComplexCase1", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{-120, 50, 150, 15}}, map[string]outputType{"toto1": outputType{0, 200, 25, 7}}}}, want: map[string]outputType{"toto1": outputType{-120, 200, 110, 22}}},
		{name: "testWithMoreComplexCase2", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{-20, 200, 10, 20}}, map[string]outputType{"toto1": outputType{50, 80, 60, 9}}}}, want: map[string]outputType{"toto1": outputType{-20, 200, 110, 25}}},
	}

	b.ResetTimer()
	for idx, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				mergeResult(tests[idx].args.data...)

			}
		})
	}

}

func Test_mergeMinMaxAvg(t *testing.T) {
	type args struct {
		data []map[string]outputType
	}
	tests := []struct {
		name string
		args args
		want map[string]outputType
	}{
		{name: "testOneMap", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{1, 1, 1, 1}}}}, want: map[string]outputType{"toto1": outputType{1, 1, 1, 1}}},
		{name: "testTwoMapAvg", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{1, 1, 10, 1}}, map[string]outputType{"toto1": outputType{1, 1, 20, 1}}}}, want: map[string]outputType{"toto1": outputType{1, 1, 15, 2}}},
		{name: "testTwoMapMin", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{2, 1, 10, 1}}, map[string]outputType{"toto1": outputType{5, 1, 20, 1}}}}, want: map[string]outputType{"toto1": outputType{2, 1, 15, 2}}},
		{name: "testTwoMapMax", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{1, 2, 10, 1}}, map[string]outputType{"toto1": outputType{1, 5, 20, 1}}}}, want: map[string]outputType{"toto1": outputType{1, 5, 15, 2}}},
		{name: "testWithMoreComplexCase1", args: args{[]map[string]outputType{map[string]outputType{"toto1": outputType{-120, 50, 150, 15}}, map[string]outputType{"toto1": outputType{0, 200, 25, 7}}}}, want: map[string]outputType{"toto1": outputType{-120, 200, 110, 22}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeResult(tt.args.data...)
			tmpGot := got["toto1"]
			tmpWant := tt.want["toto1"]
			if tmpGot.Min != tmpWant.Min {
				t.Errorf("Min Value incorrect got %v, want %v", got["toto1"].Min, tt.want["toto1"].Min)
			}
			if tmpGot.Max != tmpWant.Max {
				t.Errorf("Max Value incorrect got %v, want %v", got["toto1"].Max, tt.want["toto1"].Max)
			}
			if tmpGot.Avg != tmpWant.Avg {
				t.Errorf("Avg Value incorrect got %v, want %v", got["toto1"].Avg, tt.want["toto1"].Avg)
			}
			if tmpGot.Nb != tmpWant.Nb {
				t.Errorf("NB Value incorrect got %v, want %v", got["toto1"].Nb, tt.want["toto1"].Nb)
			}
		})
	}

}
func Test_mergeResult(t *testing.T) {

	type args struct {
		data []map[string]outputType
	}
	type test = struct {
		name string
		args args
		want map[string]outputType
	}
	var tests []test

	//NOTE: want creation
	// for the test, i just use client{{i}}= i,i,i as value that i need to find
	want := map[string]outputType{}

	for i := 0; i < runtime.NumCPU(); i++ {

		want[fmt.Sprintf("client%v", i)] = outputType{i, i, i, 1}
	}

	//NOTE: args creation
	// I test with 2,3,4,5 ... 12 maps for the sake of it
	for i := 2; i < runtime.NumCPU(); i++ {

		var tmp = test{}
		tmp.name = fmt.Sprintf("test with %v", i)
		var tmpArgs = args{}
		for nbMap := 0; nbMap < i; nbMap++ {

			mapSplit := runtime.NumCPU() / i

			for x := 0; x < runtime.NumCPU(); x += mapSplit {

				var tmpMap = map[string]outputType{}

				for j := 0; j < mapSplit; j++ {
					if j < mapSplit {
						tmpMap[fmt.Sprintf("client%v", j+x)] = outputType{j + x, j + x, j + x, 1}
					}

				}
				tmpArgs.data = append(tmpArgs.data, tmpMap)
			}
		}
		tmp.args = tmpArgs

		tmp.want = want
		tests = append(tests, tmp)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := mergeResult(tt.args.data...)
			if len(got) != len(tt.want) {
				t.Errorf("incorrect number of entry receive, got %v, want %v", got, tt.want)
			}

		})
	}
}

func Test_processContent(t *testing.T) {
	const filename string = "tests/test100.txt"
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		nbValue int
	}{
		{name: "test100", args: args{data: loadFile(filename)}, nbValue: 86},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processContent(tt.args.data)
			if len(got) != tt.nbValue { //NOTE: nb line - doublons
				t.Errorf("Incorrect number of entry got : %v, want %v", len(got), tt.nbValue)
			}
			value, ok := got["Kingston"]
			fmt.Println(got)

			if !ok {
				t.Errorf("Kingston missing from res")
			}
			if value.Nb != 2 {
				t.Errorf("Incorrect number of value for Kingston got %v", value.Nb)

			}
			if value.Max != 343 {

				t.Errorf("Incorrect max of value for Kingston got %v", value.Max)
			}
			if value.Min != 262 {

				t.Errorf("Incorrect Min of value for Kingston got %v", value.Min)
			}
			if value.Avg != 302 {

				t.Errorf("Incorrect Min of value for Kingston got %v", value.Min)
			}

		})
	}
}

func BenchmarkInterpretLine(b *testing.B) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 int
	}{
		{name: "test1Case", args: args{data: []byte{84, 111, 97, 109, 97, 115, 105, 110, 97, 59, 50, 50, 46, 48, 10}}, want: []byte{84, 111, 97, 109, 97, 115, 105, 110, 97}, want1: 220},
		{name: "test2Case", args: args{data: []byte{84, 101, 108, 32, 65, 118, 105, 118, 59, 50, 54, 46, 55, 10}}, want: []byte{84, 101, 108, 32, 65, 118, 105, 118}, want1: 267},
	}
	for idx, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {

				interpretLine(tests[idx].args.data)

			}
		})
	}

}
func Test_interpretLine(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 int
	}{
		{name: "test1Case", args: args{data: []byte{84, 111, 97, 109, 97, 115, 105, 110, 97, 59, 50, 50, 46, 48, 10}}, want: []byte{84, 111, 97, 109, 97, 115, 105, 110, 97}, want1: 220},
		{name: "test2Case", args: args{data: []byte{84, 101, 108, 32, 65, 118, 105, 118, 59, 50, 54, 46, 55, 10}}, want: []byte{84, 101, 108, 32, 65, 118, 105, 118}, want1: 267},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := interpretLine(tt.args.data)
			if string(got) != string(tt.want) {
				t.Errorf("interpretLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("interpretLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
func BenchmarkInterpretValue(b *testing.B) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name  string
		args  args
		want1 int
	}{
		//22.0 26.7 -26.7 1.5 -2.4 0.0
		{name: "TestPositiveValue", args: args{data: []byte{50, 50, 46, 48, 10}}, want1: 220},
		{name: "testPositiveValue", args: args{data: []byte{50, 54, 46, 55, 10}}, want1: 267},
		{name: "testNegativeValue", args: args{data: []byte{45, 50, 54, 46, 55, 10}}, want1: -267},
		{name: "testSingleDigit", args: args{data: []byte{49, 46, 53, 10}}, want1: 15},
		{name: "testNegSingleDigit", args: args{data: []byte{45, 50, 46, 52, 10}}, want1: -24},
		{name: "testZero", args: args{data: []byte{48, 46, 48, 10}}, want1: 0},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(t *testing.B) {
			for i := 0; i < b.N; i++ {

				interpretValue(tt.args.data)
			}

		})
	}

}

func Test_interpretValue(t *testing.T) {

	type args struct {
		data []byte
	}
	tests := []struct {
		name  string
		args  args
		want1 int
	}{
		//22.0 26.7 -26.7 1.5 -2.4 0.0
		{name: "TestPositiveValue", args: args{data: []byte{50, 50, 46, 48, 10}}, want1: 220},
		{name: "testPositiveValue", args: args{data: []byte{50, 54, 46, 55, 10}}, want1: 267},
		{name: "testNegativeValue", args: args{data: []byte{45, 50, 54, 46, 55, 10}}, want1: -267},
		{name: "testSingleDigit", args: args{data: []byte{49, 46, 53, 10}}, want1: 15},
		{name: "testNegSingleDigit", args: args{data: []byte{45, 50, 46, 52, 10}}, want1: -24},
		{name: "testZero", args: args{data: []byte{48, 46, 48, 10}}, want1: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := interpretValue(tt.args.data)
			if got != tt.want1 {
				t.Errorf("interpretLine() got = %v, want %v", got, tt.want1)
			}
		})
	}
}

func Test_showResult(t *testing.T) {
	type args struct {
		res map[string]outputType
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			showResult(tt.args.res)
		})
	}
}
