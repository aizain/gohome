package compression

import (
	"fmt"
	"testing"
)

func Test_zstdEncode(t *testing.T) {

	tests := []struct {
		name string
		str  string
	}{
		{
			str: "{\\\"total\\\":1,\\\"dataIds\\\":[\\\"abcdefgadsahuioca\\\"]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := zstdEncode([]byte(tt.str))
			fmt.Println(got)
		})
	}

	g := zstdEncode([]byte("你很骄傲123你很骄傲对对对"))
	got := make([]int, 0, len(g))
	for _, c := range g {
		got = append(got, 256-int(c))
	}
	//g1 := string(got)
	fmt.Println(got)
	//fmt.Println(g1)

}

//func Test_zstdEncode(t *testing.T) {
//	type args struct {
//		str []byte
//	}
//	tests := []struct {
//		name string
//		args args
//		want []byte
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := zstdEncode(tt.args.str); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("zstdEncode() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
