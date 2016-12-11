package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/midlang/mid/testdata/generated/go_beans/demo"
)

func main() {
	info := demo.Info{
		User: demo.User{
			Id:         1,
			Name:       "user1",
			OtherNames: []string{"用户1", "ユーザー1"},
			Code:       [6]byte{1, 2, 3, 4, 5, 6},
		},
		Desc: "test info",
		Xxx: map[int64][]map[int][5]bool{
			10: []map[int][5]bool{
				map[int][5]bool{
					100: [5]bool{true, true, false, true, true},
					200: [5]bool{true, true, false, false, true},
				},
				map[int][5]bool{
					300: [5]bool{true, true, false, false, false},
				},
				map[int][5]bool{},
			},
			20: []map[int][5]bool{
				map[int][5]bool{
					400: [5]bool{true, false, false, false, false},
				},
			},
			30: []map[int][5]bool{},
		},
		A: 1,
		B: 2,
		C: 3,
		D: 4,
		E: 5,
		F: 6,
		G: 7,
		H: 8,
		I: 9,
		J: 10,
		K: true,
		L: 12,
	}

	var buf bytes.Buffer
	if err := info.Encode(&buf); err != nil {
		fmt.Printf("Encode error: %v\n", err)
		return
	}
	info2 := demo.Info{}
	if err := info2.Decode(&buf); err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return
	}
	data, _ := json.MarshalIndent(info2, "", "    ")
	fmt.Printf("decoded data: %v\n", string(data))
}
