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

var Score = map[string]int{"P": 1000, "S": 1000, "N": 4000, "B": 3000, "A": 4000, "R": 6000, "I": 9000, "E": 9000, "Q": 9000, "K": 100000,
	"PP": 500, "PS": 1000, "PN": 4000, "PB": 3000, "PA": 4000, "PR": 3000, "PI": 4500, "PE": 9000, "PQ": 4500, "PK": 50000,
	"SP": 1000, "SS": 500, "SN": 4000, "SB": 1500, "SA": 4000, "SR": 6000, "SI": 9000, "SE": 4500, "SQ": 9000, "SK": 100000,
	"NP": 1000, "NS": 1000, "NN": 2000, "NB": 3000, "NA": 4000, "NR": 6000, "NI": 4500, "NE": 9000, "NQ": 9000, "NK": 100000,
	"BP": 1000, "BS": 1000, "BN": 4000, "BB": 1500, "BA": 4000, "BR": 6000, "BI": 9000, "BE": 4500, "BQ": 4500, "BK": 100000,
	"AP": 1000, "AS": 1000, "AN": 4000, "AB": 3000, "AA": 2000, "AR": 6000, "AI": 9000, "AE": 4500, "AQ": 9000, "AK": 100000,
	"RP": 1000, "RS": 1000, "RN": 4000, "RB": 3000, "RA": 4000, "RR": 3000, "RI": 4500, "RE": 9000, "RQ": 4500, "RK": 100000,
	"IP": 1000, "IS": 1000, "IN": 4000, "IB": 3000, "IA": 4000, "IR": 6000, "II": 4500, "IE": 9000, "IQ": 9000, "IK": 100000,
	"EP": 1000, "ES": 1000, "EN": 4000, "EB": 3000, "EA": 4000, "ER": 6000, "EI": 9000, "EE": 4500, "EQ": 9000, "EK": 100000,
	"QP": 1000, "QS": 1000, "QN": 4000, "QB": 3000, "QA": 4000, "QR": 6000, "QI": 9000, "QE": 4500, "QQ": 4500, "QK": 100000,
	"KP": 1000, "KS": 1000, "KN": 4000, "KB": 3000, "KA": 4000, "KR": 6000, "KI": 9000, "KE": 4500, "KQ": 9000, "KK": 100000}
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

func InPromotePos(hex string) bool {
	return strings.Contains("f5 f51 f52 f53 f54 f55 a5 a51 a52 a53 a54 a55 b5 c5 c51 c52 c53 c54 c55 d5 d51 d52 d53 d54 d55 e5 ", hex+" ")
}

func Add2MapArr(theMapArr map[int][]int, xy, add int) {
	if !OnBoard(add) {
		return
	}
	theMapArr[xy] = append(theMapArr[xy], add)
}
