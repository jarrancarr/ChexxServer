package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func CorsHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	CorsHeader(w)
	json.NewEncoder(w).Encode(data)
}

var Score = map[string]int{"P": 1000, "S": 1000, "N": 4000, "B": 3000, "A": 4000, "R": 6000, "I": 9000, "E": 9000, "Q": 9000, "K": 100000}
var Hex = map[int]string{612: "*", 610: "a", 711: "b", 713: "c", 614: "d", 513: "e", 511: "f"}
var XY = map[string]int{}

func Init() {

	for i := 1; i < 6; i++ {
		Hex[610-i*2] = fmt.Sprintf("a%d", i)
		Hex[711+99*i] = fmt.Sprintf("b%d", i)
		Hex[713+101*i] = fmt.Sprintf("c%d", i)
		Hex[614+i*2] = fmt.Sprintf("d%d", i)
		Hex[513-99*i] = fmt.Sprintf("e%d", i)
		Hex[511-101*i] = fmt.Sprintf("f%d", i)
		for j := 1; j < i+1; j++ {
			Hex[610+101*j-2*i] = fmt.Sprintf("a%d%d", i, j)
			Hex[711+99*i+2*j] = fmt.Sprintf("b%d%d", i, j)
			Hex[713+101*i-99*j] = fmt.Sprintf("c%d%d", i, j)
			Hex[614-101*j+2*i] = fmt.Sprintf("d%d%d", i, j)
			Hex[513-99*i-2*j] = fmt.Sprintf("e%d%d", i, j)
			Hex[511-101*i+99*j] = fmt.Sprintf("f%d%d", i, j)
		}
	}
	for xy, Hex := range Hex {
		XY[Hex] = xy
	}
}

func OnBoard(xy int) bool {
	y := xy % 100
	x := (xy - y) / 100
	return !(x < 0 || x > 12 || x+y < 6 || x+y > 31 || x-y > 6 || y-x > 19)
}

func InStartPos(hex string, isWhite bool) bool {
	if isWhite {
		return strings.Contains("d55 d44 d33 d21 c22 c31 c41 c51 d43 d32 d2 c32 c42 ", hex+" ")
	}
	return strings.Contains("f51 f41 f31 f22 a21 a33 a44 a55 f42 f32 a2 a32 a43 ", hex+" ")
}

func Add2MapArr(theMapArr map[int][]int, xy, add int) {
	if !OnBoard(add) {
		return
	}
	theMapArr[xy] = append(theMapArr[xy], add)
}
