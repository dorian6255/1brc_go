package main

import (
	"reflect"
	"testing"
)

func Test_loadFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want *[]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadFile(tt.args.filename); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitContent(t *testing.T) {
	type args struct {
		content *[]byte
	}
	tests := []struct {
		name string
		args args
		want []*[]byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitContent(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeResult(t *testing.T) {
	type args struct {
		data []map[string]outputType
	}
	tests := []struct {
		name string
		args args
		want map[string]outputType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeResult(tt.args.data...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_process(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]outputType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := process(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("process() = %v, want %v", got, tt.want)
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
	}
	for i := 0; i < b.N; i++ {

		interpretLine(tests[0].args.data)
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
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("interpretLine() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("interpretLine() got1 = %v, want %v", got1, tt.want1)
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
