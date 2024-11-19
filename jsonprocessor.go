package dynamictags

import "encoding/json"

func ProcessJson(str string) {
	var res interface{}
	json.Unmarshal([]byte(str), res)
}
