package utils

import (
	"encoding/json"
	"fmt"
)

// JsonDump 打印
func JsonDump(v any) {
	marshal, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(marshal))
}
