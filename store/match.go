package store

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jarrancarr/ChexxServer/utils"
	"github.com/jinzhu/gorm"
)

type Army struct {
	UserId uint     `json:"userid"`
	Pieces []string `json:"pieces"`
	Time   int      `json:"time"`
}

type Type struct {
	Name      string `json:"name"`
	GameClock uint32 `json:"game"`
	MoveClock uint32 `json:"move"`
	Status    string `json:"status`
	Rating    uint32 `json:"rating"`
}
type Match struct {
	gorm.Model
	Title         string   `json:"name"`
	White         Army     `json:"white" gorm:"-"`
	Black         Army     `json:"black" gorm:"-"`
	WhiteArmy     string   `json:"-"`
	BlackArmy     string   `json:"-"`
	WhitePlayerId uint     `json:"-"`
	BlackPlayerId uint     `json:"-"`
	Log           []string `json:"log" gorm:"-"`
	Logs          string   `json:"-"`
	Blog          string   `json:"blog"`
	Notes         string   `json:"notes"`
	Game          Type     `gorm:"embedded"`
	board         *Board   `json:"-" gorm:"-"`
	LastMove      string   `json:"-" gorm:"-"`
}

type Board struct {
	Attacks      map[int][]int
	Moves        map[int][]int
	Attacked     map[int][2][]int
	Occupant     map[int]string // [wb][KQIERABNPS]
	Pinned       map[int]int    // map[pinned]skewerer so pinned piece can attack
	WhiteInCheck bool
	BlackInCheck bool
	Mate         bool // if no one is in check, then stalemate
	Score        int
	MoveCount    int
}

func (m *Match) Validate() (map[string]interface{}, bool) {

	if m.Title == "" {
		return utils.Message(false, "No Title"), false
	}

	return utils.Message(false, "Requirement passed"), true
}

func (m *Match) Create() map[string]interface{} {

	if resp, ok := m.Validate(); !ok {
		return resp
	}

	m.Logs = strings.Join(m.Log, " ")
	m.BlackArmy = strings.Join(m.Black.Pieces, " ") + "|" + fmt.Sprintf("%d", m.Black.Time)
	m.WhiteArmy = strings.Join(m.White.Pieces, " ") + "|" + fmt.Sprintf("%d", m.White.Time)
	m.BlackPlayerId = m.Black.UserId
	m.WhitePlayerId = m.White.UserId

	GetDB().Create(m)

	if m.ID <= 0 {
		return utils.Message(false, "Failed to create message, connection error.")
	}

	response := utils.Message(true, "Match created")
	response["id"] = m.ID
	return response
}
func (m *Match) Update() map[string]interface{} {

	if resp, ok := m.Validate(); !ok {
		return resp
	}

	m.Logs = strings.Trim(strings.Join(m.Log, " "), " ")
	m.BlackArmy = strings.Join(m.Black.Pieces, " ")
	m.WhiteArmy = strings.Join(m.White.Pieces, " ")
	strings.Trim(strings.Replace(m.WhiteArmy, "  ", " ", -1), " ")
	strings.Trim(strings.Replace(m.BlackArmy, "  ", " ", -1), " ")
	m.WhiteArmy += "|" + fmt.Sprintf("%d", m.White.Time)
	m.BlackArmy += "|" + fmt.Sprintf("%d", m.Black.Time)
	m.BlackPlayerId = m.Black.UserId
	m.WhitePlayerId = m.White.UserId

	GetDB().Save(m)

	if m.ID <= 0 {
		return utils.Message(false, "Failed to create message, connection error.")
	}

	response := utils.Message(true, "Match updated")
	return response
}

func GetMatch(id uint) *Match {

	m := &Match{}
	GetDB().Table("matches").Where("id = ?", id).First(m)

	m.Log = strings.Split(strings.Trim(m.Logs, " "), " ")

	black := strings.Split(m.BlackArmy, "|")
	blackClock, _ := strconv.Atoi(black[1])
	m.Black = Army{UserId: m.BlackPlayerId, Pieces: strings.Split(strings.Trim(black[0], " "), " "), Time: blackClock}

	white := strings.Split(m.WhiteArmy, "|")
	whiteClock, _ := strconv.Atoi(white[1])
	m.White = Army{UserId: m.WhitePlayerId, Pieces: strings.Split(strings.Trim(white[0], " "), " "), Time: whiteClock}

	return m
}

func (m *Match) Move(move string) {
	m.LastMove = move
	pos := strings.Split(move, "~")
	repl := -1
	for p := range m.White.Pieces {
		if m.White.Pieces[p][1:] == pos[0] {
			repl = p
		}
	}
	if repl > -1 {
		if m.White.Pieces[repl][1:] == pos[0] {
			m.White.Pieces[repl] = m.White.Pieces[repl][:1] + pos[1]
		}
		// now remove black piece
		rem := -1
		for p := range m.Black.Pieces {
			if m.Black.Pieces[p][1:] == pos[1] {
				rem = p
			}
		}
		if rem > -1 {
			m.Log = append(m.Log, m.White.Pieces[repl][:1]+pos[0]+"x"+m.Black.Pieces[rem])
			m.Black.Pieces = append(m.Black.Pieces[:rem], m.Black.Pieces[rem+1:]...)
		} else {
			m.Log = append(m.Log, m.White.Pieces[repl][:1]+pos[0]+"~"+pos[1])
		}
	} else {
		for p := range m.Black.Pieces {
			if m.Black.Pieces[p][1:] == pos[0] {
				repl = p
			}
		}
		if repl > -1 {
			if m.Black.Pieces[repl][1:] == pos[0] {
				m.Black.Pieces[repl] = m.Black.Pieces[repl][:1] + pos[1]
			}
			// now remove white piece
			rem := -1
			for p := range m.White.Pieces {
				if m.White.Pieces[p][1:] == pos[1] {
					rem = p
				}
			}
			if rem > -1 {
				m.Log = append(m.Log, m.Black.Pieces[repl][:1]+pos[0]+"x"+m.White.Pieces[rem])
				m.White.Pieces = append(m.White.Pieces[:rem], m.White.Pieces[rem+1:]...)
			} else {
				m.Log = append(m.Log, m.Black.Pieces[repl][:1]+pos[0]+"~"+pos[1])
			}
		}
	}
	// TODO: still have to do special moves

	// TODO: checkmate?
	// TODO: stalemate?
	// TODO: timeout?
	//return m.Update()
}

func (m *Match) clone(origin *Match) {
	m.White.Pieces = make([]string, len(origin.White.Pieces))
	m.Black.Pieces = make([]string, len(origin.Black.Pieces))
	for p := range origin.White.Pieces {
		m.White.Pieces[p] = origin.White.Pieces[p]
	}
	for p := range origin.Black.Pieces {
		m.Black.Pieces[p] = origin.Black.Pieces[p]
	}
}

func (m *Match) AI(width, depth int, finished chan bool) *Match {
	// fmt.Printf("*-*-*AI    w:%d h:%d\n", width, depth)
	if width < 1 {
		width = 1
	}
	m.Analyse()
	lowest := 0
	bestNMoves := make([]Match, 0)
	isWhite := len(m.Log)%2 == 0
	for perp, assaults := range m.board.Moves {
		for attack := range assaults {
			if isWhite {
				if m.board.Occupant[perp][0] == 'w' {
					testMove := utils.Hex[perp] + "~" + utils.Hex[assaults[attack]]
					testMatch := Match{White: Army{Pieces: []string{}}, Black: Army{Pieces: []string{}}}
					testMatch.clone(m)
					testMatch.Move(testMove)
					testMatch.Analyse()
					// fmt.Printf("*-*-*:::{%v} {%v} move:%s    score:%d\n", testMatch.White.Pieces, testMatch.Black.Pieces, testMatch.LastMove, testMatch.board.Score)
					if !testMatch.board.WhiteInCheck {
						if len(bestNMoves) < width {
							bestNMoves = append(bestNMoves, testMatch)
						} else if testMatch.board.Score > bestNMoves[lowest].board.Score {
							bestNMoves[lowest] = testMatch
						}
						for l := range bestNMoves {
							if bestNMoves[l].board.Score < bestNMoves[lowest].board.Score {
								lowest = l
							}
						}
					}
				}
			} else {
				if m.board.Occupant[perp][0] == 'b' {
					testMove := utils.Hex[perp] + "~" + utils.Hex[assaults[attack]]
					testMatch := Match{White: Army{Pieces: []string{}}, Black: Army{Pieces: []string{}}}
					testMatch.clone(m)
					testMatch.Move(testMove)
					testMatch.Analyse()
					// fmt.Printf("*-*-*:::{%v} {%v} move:%s    score:%d\n", testMatch.White.Pieces, testMatch.Black.Pieces, testMatch.LastMove, testMatch.board.Score)
					if !testMatch.board.BlackInCheck {
						if len(bestNMoves) < width {
							bestNMoves = append(bestNMoves, testMatch)
						} else if testMatch.board.Score < bestNMoves[lowest].board.Score {
							bestNMoves[lowest] = testMatch
						}
						for l := range bestNMoves {
							if bestNMoves[l].board.Score > bestNMoves[lowest].board.Score {
								lowest = l
							}
						}
					}
				}
			}
		}
	}
	if len(bestNMoves) == 0 {
		// fmt.Printf("*-*-**-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*\n no move found for position: \n%v\n%v\n*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*\n", m.White.Pieces, m.Black.Pieces)
		if finished != nil {
			finished <- true // let parent know i am done
		}
		return nil
	}
	highest := 0
	done := make(chan bool)
	if depth > 1 {
		for l := range bestNMoves {

			bestNMoves[l].Log = append(bestNMoves[l].Log, "xxx")
			go bestNMoves[l].AI(width-1, depth-1, done)
		}
		for rec := 0; rec < len(bestNMoves); rec += 1 {
			// fmt.Printf("*-*-*waiting on children\n")
			<-done
			// fmt.Printf("*-*-*child %d done\n", rec)
		}
	}
	for l := range bestNMoves {
		// fmt.Printf("*-*-*%v, %v, %s....%d\n", bestNMoves[l].White.Pieces, bestNMoves[l].Black.Pieces, bestNMoves[l].LastMove, bestNMoves[l].board.Score)
		if isWhite {
			if bestNMoves[l].board.Score > bestNMoves[highest].board.Score {
				highest = l
			}
		} else {
			if bestNMoves[l].board.Score < bestNMoves[highest].board.Score {
				highest = l
			}
		}
	}
	m.board.Score = bestNMoves[highest].board.Score
	if finished != nil {
		finished <- true // let parent know i am done
	}

	if len(bestNMoves) == 0 {
		return nil // no valid move
	}
	return &bestNMoves[highest]
}

func (match *Match) slide(xy int, direction []int) {
	dirs := []int{}
	for d := range direction {
		dirs = append(dirs, direction[d], -direction[d])
	}
	for d := range dirs {
		keepGoing := true
		for pos := xy + dirs[d]; utils.OnBoard(pos) && keepGoing; pos += dirs[d] {
			victim, iam := match.board.Occupant[pos]
			if !iam {
				utils.Add2MapArr(match.board.Attacks, xy, pos)
			} else {
				keepGoing = false
				utils.Add2MapArr(match.board.Attacks, xy, pos)
				if victim[:1] != match.board.Occupant[xy][:1] { // not same team
					// now see if pos piece is pinned
					for cont := pos + dirs[d]; utils.OnBoard(cont); cont = cont + dirs[d] {
						isKing, ami := match.board.Occupant[cont]
						if ami {
							if isKing[0] != match.board.Occupant[xy][0] && isKing[1] == 'K' {
								// this is the enemy king, so piece is pinned
								match.board.Pinned[pos] = xy
							} else { // this is a piece, but not the enemy king so stop.
								cont = 0
							}
						}
					}
				}
			}
		}
	}
}

func (match *Match) jump(xy int, direction []int) {
	// fmt.Printf("*-*-*jump %d\n", xy)
	for d := range direction {
		utils.Add2MapArr(match.board.Attacks, xy, xy+direction[d])
		utils.Add2MapArr(match.board.Attacks, xy, xy-direction[d])
	}
}

func (match *Match) rookMoves(xy int)   { match.slide(xy, []int{2, 101, 99}) }
func (match *Match) bishopMoves(xy int) { match.slide(xy, []int{103, 97, 200}) }
func (match *Match) knightMoves(xy int) { match.jump(xy, []int{105, 95, 301, 299, 204, 196}) }
func (match *Match) archerMoves(xy int) { match.jump(xy, []int{107, 93, 305, 295, 402, 398}) }
func (match *Match) kingMoves(xy int)   { match.jump(xy, []int{2, 99, 101}) }
func (match *Match) pawnMoves(xy, dirs int, inStart bool) {
	utils.Add2MapArr(match.board.Attacks, xy, xy+2*dirs)
	_, block := match.board.Occupant[xy+2*dirs]
	if inStart && !block {
		utils.Add2MapArr(match.board.Attacks, xy, xy+4*dirs)
	}
}
func (match *Match) pawnAttacks(xy, dirs int) {
	utils.Add2MapArr(match.board.Attacks, xy, xy+dirs-100)
	utils.Add2MapArr(match.board.Attacks, xy, xy+dirs+100)
	utils.Add2MapArr(match.board.Attacks, xy, xy) // switch arms
}
func (match *Match) TestAttacks(hex string) {
	match.board = &Board{Occupant: make(map[int]string, len(match.White.Pieces)+len(match.Black.Pieces)), Attacks: make(map[int][]int), Moves: make(map[int][]int), Attacked: make(map[int][2][]int), Pinned: make(map[int]int, 0)}

	for m := range match.White.Pieces {
		p := match.White.Pieces[m]
		match.board.Occupant[utils.XY[p[1:]]] = "w" + p[:1]
	}
	for m := range match.Black.Pieces {
		p := match.Black.Pieces[m]
		match.board.Occupant[utils.XY[p[1:]]] = "b" + p[:1]
	}
	match.attacks(hex)
}

func (match *Match) Show(hex string) {
	fmt.Printf("%v\n", match.board.Attacks[utils.XY[hex]])
}

func (match *Match) attacks(hex string) {
	xy := utils.XY[hex]
	dirs := 1
	isWhite := match.board.Occupant[xy][0] == 'w'
	if isWhite {
		dirs = -1
	}
	switch match.board.Occupant[xy][1] {
	case 'R':
		match.rookMoves(xy)
	case 'N':
		match.knightMoves(xy)
	case 'B':
		match.bishopMoves(xy)
	case 'A':
		match.archerMoves(xy)
	case 'K':
		match.kingMoves(xy)
	case 'Q':
		match.rookMoves(xy)
		match.bishopMoves(xy)
	case 'I':
		match.rookMoves(xy)
		match.knightMoves(xy)
	case 'E':
		match.bishopMoves(xy)
		match.archerMoves(xy)
	case 'P':
		match.pawnMoves(xy, dirs, utils.InStartPos(hex, isWhite))
		match.pawnAttacks(xy, dirs)
	case 'S':
		match.pawnMoves(xy, dirs, utils.InStartPos(hex, isWhite))
		match.pawnAttacks(xy, dirs*3)
	}
}

func (match *Match) Analyse() {
	match.board = &Board{Occupant: make(map[int]string, len(match.White.Pieces)+len(match.Black.Pieces)), Attacks: make(map[int][]int), Moves: make(map[int][]int), Attacked: make(map[int][2][]int), Pinned: make(map[int]int, 0)}

	// fmt.Printf("*-*-*\n   $$$  W{")
	for m := range match.White.Pieces {
		p := match.White.Pieces[m]
		match.board.Score += utils.Score[p[:1]]
		if p[0] == 'P' || p[0] == 'S' {

		}
		// fmt.Printf("*-*-*%s ", p)
		match.board.Occupant[utils.XY[p[1:]]] = "w" + p[:1]
	}
	// fmt.Printf("*-*-*}  B{")
	for m := range match.Black.Pieces {
		p := match.Black.Pieces[m]
		match.board.Score -= utils.Score[p[:1]]
		// fmt.Printf("*-*-*%s ", p)
		match.board.Occupant[utils.XY[p[1:]]] = "b" + p[:1]
	}
	// fmt.Printf("*-*-*}  ")

	for m := range match.White.Pieces {
		match.attacks(match.White.Pieces[m][1:])
	}
	for m := range match.Black.Pieces {
		match.attacks(match.Black.Pieces[m][1:])
	}
	// fmt.Printf("*-*-*Attacks: %v", match.board.Attacks)
	for xy, attacks := range match.board.Attacks {
		perp := match.board.Occupant[xy]
		// fmt.Printf("*-*-*\n      %s(", perp)
		tally := 0
		for target := range attacks {
			y := attacks[target] % 100
			x := (attacks[target] - y) / 100
			dist := (12-y)*(12-y) + (6-x)*(6-x)
			victim, iam := match.board.Occupant[attacks[target]]
			// fmt.Printf("*-*-*-%s", utils.Hex[attacks[target]])
			// fmt.Printf("*-*-*---%s---%d---%s----%v----%s-----\n", utils.Hex[attacks[target]], dist, perp, iam, victim)
			if perp[0] == 'w' {
				match.board.Score += 200 - dist
				tally += 200 - dist
				if iam {
					if victim[0] == 'b' {
						// fmt.Printf("*-*-**%s", victim[1:])
						match.board.Score += utils.Score[victim[1:]] / 5
					} else {
						// fmt.Printf("*-*-*|%s", victim[1:])
						match.board.Score += utils.Score[victim[1:]] / 15
					}
				}
				match.board.Attacked[attacks[target]] = [2][]int{append(match.board.Attacked[attacks[target]][0], xy), match.board.Attacked[attacks[target]][1]}
			} else {
				match.board.Score += dist - 200
				tally += dist - 200
				if iam {
					if victim[0] == 'w' {
						// fmt.Printf("*-*-**%s", victim[1:])
						match.board.Score -= utils.Score[victim[1:]] / 5
						tally -= utils.Score[victim[1:]] / 5
					} else {
						// fmt.Printf("*-*-*|%s", victim[1:])
						match.board.Score -= utils.Score[victim[1:]] / 15
						tally -= utils.Score[victim[1:]] / 15
					}
				}
				match.board.Attacked[attacks[target]] = [2][]int{match.board.Attacked[attacks[target]][0], append(match.board.Attacked[attacks[target]][1], xy)}
			}
			// legal move checks from xy to target
			if match.legalCheck(xy, attacks[target]) {
				match.board.Moves[xy] = append(match.board.Moves[xy], attacks[target])
				match.board.MoveCount += 1
			}
			// add special moves here
		}
		// fmt.Printf("*-*-*) = %d", tally)
	}

	// TODO: add formation moves here
	for m := range match.White.Pieces {
		if match.White.Pieces[m][0] == 'P' || match.White.Pieces[m][0] == 'S' {
			// formation attacks
		}
	}
	for m := range match.Black.Pieces {
		if match.Black.Pieces[m][0] == 'P' || match.Black.Pieces[m][0] == 'S' {
			// formation attacks
		}
	}

	king := match.getKing(true)
	_, regacide := match.board.Attacked[king]
	if king > 0 && regacide && len(match.board.Attacked[king][1]) > 0 {
		match.board.WhiteInCheck = true
		if len(match.Log)%2 == 0 && !match.canEscape(king) {
			match.board.Mate = true
		}
	}
	king = match.getKing(false)
	_, regacide = match.board.Attacked[king]
	if king > 0 && regacide && len(match.board.Attacked[king][0]) > 0 {
		match.board.BlackInCheck = true
		if len(match.Log)%2 == 1 && !match.canEscape(king) {
			match.board.Mate = true
		}
	}
	// fmt.Printf("*-*-*%s= %d  b+?%v  w+?%v ++?%v\n", match.LastMove, match.board.Score, match.board.BlackInCheck, match.board.WhiteInCheck, match.board.Mate)
}

func (match *Match) getKing(white bool) int {
	for xy, piece := range match.board.Occupant {
		if white {
			if piece[0] == 'w' && piece[1] == 'K' {
				return xy
			}
		} else {
			if piece[0] == 'b' && piece[1] == 'K' {
				return xy
			}
		}
	}
	return 0
}

func (match *Match) canEscape(king int) bool { // is any surrounding space free of attacks and not occupied by friendly?
	index := 0
	if match.board.Occupant[king][0] == 'w' {
		index = 1
	}
	for d := range []int{2, 99, 101, -2, -99, -101} {
		dest, iam := match.board.Occupant[king+d]
		if !iam {
			dest = "xX"
		}
		if utils.OnBoard(king+d) && match.board.Occupant[king][0] != dest[0] && len(match.board.Attacked[king+d][index]) > 0 {
			// is it a rook attack? rook attacking in same direction as escape... still in check
			kingPinned := false
			for atx := range match.board.Attacked[king+d][index] {
				atr := match.board.Occupant[atx]
				if atr != "" && (atr[1] == 'R' || atr[1] == 'Q' || atr[1] == 'I') {
					if (king-atx)%d == 0 {
						kingPinned = true
					}
				}
			}
			if !kingPinned {
				return true
			}
		}
	}
	return false
}

func (match *Match) legalCheck(xy, att int) bool {
	atr, _ := match.board.Occupant[xy]
	vic, ibe := match.board.Occupant[att]
	// fmt.Printf("*-*-*legalCheck %s:%d~%s:%d\n", atr, xy, vic, att)

	if ibe { // destination hex occupied
		if vic[0] == atr[0] {
			return false // same team
		}
		// we are enemies
		if atr[1] != 'P' && atr[1] != 'S' {
			return !match.Pinned(xy)
		}
		// attacker is a pawn or sprearman
		return !match.Pinned(xy) && ((xy-att)*(xy-att))%4 != 0 // attacks are odd for P and S
	}
	// no one there

	if atr[1] != 'P' && atr[1] != 'S' {
		return !match.Pinned(xy)
	}
	// attacker is a pawn or sprearman
	return !match.Pinned(xy) && ((xy-att)*(xy-att))%4 == 0 // attacks are odd for P and S
}

func (match *Match) Pinned(xy int) bool {
	for pin := range match.board.Pinned {
		if match.board.Pinned[pin] == xy {
			return true
		}
	}
	return false
}

func (match *Match) TextBoard(dat string, hex bool) {
	for y := 0; y < 25; y++ {
		if y%2 == 1 {
			fmt.Printf(" ")
		}
		for x := 0; x < 14; x++ {
			a, b := utils.Hex[x*100+y]
			if b {
				switch dat {
				case "atk":
					c, d := match.board.Attacks[x*100+y]
					if d {
						fmt.Printf("<%2d >", len(c))
					} else {
						fmt.Printf("<   >")
					}
				case "atd":
					c, d := match.board.Occupant[x*100+y]
					if d {
						fmt.Printf("<%d%s%d>", len(match.board.Attacked[x*100+y][0]), string(c[1]), len(match.board.Attacked[x*100+y][1]))
					} else {
						if hex {
							fmt.Printf("<%3s>", a)
						} else {
							if a == "*" {
								fmt.Printf("<%d*%d>", len(match.board.Attacked[x*100+y][0]), len(match.board.Attacked[x*100+y][1]))
								// fmt.Printf("*-*-*<<*>>")
							} else {
								if len(match.board.Attacked[x*100+y][0]) > 0 || len(match.board.Attacked[x*100+y][1]) > 0 {
									fmt.Printf("<%d%s%d>", len(match.board.Attacked[x*100+y][0]), []string{"#", ":", "="}[y%3], len(match.board.Attacked[x*100+y][1]))
								} else {
									fmt.Printf("<%s>", []string{"   ", ":::", "---"}[y%3])
								}
							}
						}
					}

				case "mov":
					fmt.Printf("<%2d >", len(match.board.Moves[x*100+y]))
				default:
				}
			} else {
				if (x+y)%2 == 0 {
					fmt.Printf("     ")
				} else {
					fmt.Printf("   ")
				}
			}
		}
		for x := 0; x < 14; x++ {
			a, b := utils.Hex[x*100+y]
			if b {
				c, d := match.board.Occupant[x*100+y]
				if d {
					if c[0] == 'w' {
						fmt.Printf("< %s >", strings.ToLower(string(c[1])))
					} else {
						fmt.Printf("< %s >", string(c[1]))
					}
				} else {
					if hex {
						fmt.Printf("<%3s>", a)
					} else {
						if a == "*" {
							fmt.Printf("<-*->")
						} else {
							fmt.Printf("<%s>", []string{"---", ":::", "   "}[y%3])
						}
					}
				}
			} else {
				if (x+y)%2 == 0 {
					fmt.Printf("     ")
				} else {
					fmt.Printf("   ")
				}
			}
		}
		fmt.Println()
	}
}

func (match *Match) Examine() {
	fmt.Printf("\nAttacks: %d     Attacked: %d     Moves: %d\n", len(match.board.Attacks), len(match.board.Attacked), len(match.board.Moves))
	fmt.Printf("Score: %d\n", match.board.Score)
	fmt.Printf("Log: %v\n", match.Log)

	match.TextBoard("atd", false)
}
