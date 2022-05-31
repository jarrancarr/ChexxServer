package tutor

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jarrancarr/ChexxServer/utils"
)

type Step struct {
	Title      string     `json:"title"`
	Width      int8       `json:"width"`
	PosX       int8       `json:"posX"`
	PosY       int8       `json:"posY"`
	Background string     `json:"background"`
	Text       []string   `json:"text"`
	White      []string   `json:"white"`
	Black      []string   `json:"black"`
	Hilite     [][]string `json:"hilite"`
	Log        []string   `json:"log"`
	Answer     [][]string `json:"answer"`
	Order      []string   `json:"order"`
	Hint       []string   `json:"hint"`
}

type Lesson struct {
	Step []Step `json:"step"`
}

var Lessons map[string]Lesson = map[string]Lesson{
	"Intro": {
		Step: []Step{
			{Title: "Welcome to Chexx", Width: 70, PosX: 50, PosY: 50, Background: "#00202020", Log: []string{},
				Black: []string{"Qf31", "Kf33", "Ea31", "Ia33", "Pf4", "Nf53", "Ba4", "Aa53", "Sb4"},
				White: []string{"Qd33", "Kd31", "Ec33", "Ic31", "Pe4", "Nd53", "Bd4", "Ac53", "Sc4"},
				Text: []string{
					"Like Chess, Chexx is strategic game of Army against Army.  White Army vs Black Army.",
					"Although Chexx is played on a hexagonal board; here are the types of pieces at your command.  Some are familiar, others are...well, we'll get to that."}},
			{Title: "Welcome to Chexx", Width: 70, PosX: 50, PosY: 50, Background: "#f0f02050",
				Black: []string{"Qf31", "Kf33", "Ea31", "Ia33", "Pf4", "Nf53", "Ba4", "Aa53", "Sb4"},
				White: []string{"Qd33", "Kd31", "Ec33", "Ic31", "Pe4", "Nd53", "Bd4", "Ac53", "Sc4"},
				Text:  []string{"This might seem a bit intimidating, with 9 types of pieces, but please be patient.  Even without any Chess experience, you will be able to learn step by step how each piece contributes to the fight and begin a most enjoyable journey on a battlefield of wits."}},
			{Title: "Welcome to Chexx", Width: 70, PosX: 50, PosY: 50, Background: "#00f0f050",
				Black: []string{"Qf31", "Kf33", "Ea31", "Ia33", "Pf4", "Nf53", "Ba4", "Aa53", "Sb4"},
				White: []string{"Qd33", "Kd31", "Ec33", "Ic31", "Pe4", "Nd53", "Bd4", "Ac53", "Sc4"},
				Text:  []string{"If you do have any Chess experience, you will find very little conceptually that is terribly different.  Only after a few games you should be comfortable with the hexagonal layout.  Things that are actually different is still quite derivitive from existing Chess concepts, so you should have no trouble learning them."}},
			{Width: 40, PosX: 50, PosY: 40,
				White: []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"},
				Black: []string{},
				Text:  []string{"This is White's army... 29 strong."}},
			{Width: 40, PosX: 50, PosY: 40,
				White:  []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"},
				Black:  []string{},
				Text:   []string{"You have 8 pawns..."},
				Hilite: [][]string{{"d55 d44 d33 d21 c22 c31 c41 Pc51", "stroke", "#ff4"}}},
			{Width: 40, PosX: 50, PosY: 40,
				White:  []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"},
				Black:  []string{},
				Text:   []string{"... 5 spearmen,"},
				Hilite: [][]string{{"d43 d32 d2 c32 c42", "stroke", "#ff4"}}},
			{Width: 40, PosX: 50, PosY: 40,
				White:  []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"},
				Black:  []string{},
				Text:   []string{"3 Knights (red), 3 Bishops (blue), 3 Rooks (yellow), and 3 Archers (green)..."},
				Hilite: [][]string{{"d53 d51 c33", "stroke", "#f00"}, {"c53 c55 d52", "stroke", "#00f"}, {"d54 d5 c52", "stroke", "#ff0"}, {"d42 d3 c43", "stroke", "#0f0"}}},
			{Width: 40, PosX: 50, PosY: 40,
				White:  []string{"Rd54", "Rd5", "Rc52", "Nd53", "Nd51", "Nc33", "Bc53", "Bc55", "Bd52", "Qd41", "Kc44", "Id31", "Ed4", "Pd55", "Pd44", "Pd33", "Pd21", "Pc22", "Pc31", "Pc41", "Pc51", "Sd43", "Sd32", "Sd2", "Sc32", "Sc42", "Ad42", "Ad3", "Ac43"},
				Black:  []string{},
				Text:   []string{"...and the Royal family:", "The Prince on his horse, the Princess with her bow,", "and of course, his magisty, the King and her royal highness, the Queen."},
				Hilite: [][]string{{"d41 c44 d31 d4", "stroke", "#a0f"}}},
			{Width: 40, PosX: 50, PosY: 60,
				White: []string{},
				Black: []string{"Ra5", "Rf52", "Ra54", "Nf53", "Nf55", "Na31", "Ba53", "Ba51", "Bf54", "Qf44", "Ka41", "If33", "Ea4", "Pf51", "Pf41", "Pf31", "Pf22", "Pa21", "Pa33", "Pa44", "Pa55", "Sf42", "Sf32", "Sa2", "Sa32", "Sa43", "Af43", "Aa3", "Aa42"},
				Text:  []string{"Black has the identical army."}},
		},
	},

	"Interface": {
		Step: []Step{
			{Title: "Interface", Width: 40, PosX: 70, PosY: 30, Background: "#f0f02050", Order: []string{"menu tutor", "draw circle 91,8,8,#F00"}, Log: []string{},
				Black: []string{}, White: []string{},
				Text: []string{"All help and tutorials are accessed through the <?> menu in the upper right.", "When selected, all the individual tutorials are shown in the black circular area."}},
			{Title: "Main Menu", Width: 40, PosX: 70, PosY: 60, Background: "#f0f02050",
				Black: []string{}, White: []string{}, Order: []string{"menu main", "draw circle 91,91,8,#F00"},
				Text: []string{"On the lower right, it the main menu.", "Here is where you will find access to logging in, matches and quite a bit more."}},
			{Title: "Clocks", Width: 40, PosX: 30, PosY: 50, Background: "#f0f02050",
				Black: []string{}, White: []string{}, Order: []string{"menu ", "draw circle 9,9,8,#F00", "draw circle 9,91,8,#F00"},
				Text: []string{"To the left of the board are the clocks.",
					"During a match, you will see how much time you have left.  Also, you move isn't confirmed until you tap your clock.  Until then, your move will not be final and your clock will run down.",
					"There are two hands on the clock.  The larger black one is your move clock.  It will circle once around and stop.  Then, the red hand begins to count down.  This is your game clock.  Once it reaches the top, its game over... you lost."}},
			{Title: "Ledger", Width: 55, PosX: 50, PosY: 50, Background: "#00606090",
				Black: []string{}, White: []string{},
				Log: []string{"Pc31~b22", "Nf53~f21", "Id31~d", "Sa32~a22", "IdxSa2", "Qf44xIa2", "Ed4~c11", "Aa3~f"},
				Text: []string{"The thick black ring around the board has a few purposes.",
					"You will have notices that menu options are shown here.",
					"During a match, though, it records each move with tabs that circle the board.  Hovinging over one highlights that move on the board.  Clicking on a tab will give detailed information about the move."}},
		},
	},
	"Rules": {
		Step: []Step{
			{Title: "Rules", Width: 70, PosX: 50, PosY: 50, Background: "#20202050", Log: []string{},
				Black: []string{}, White: []string{},
				Text: []string{"A Chexx match begins with White making the first move and then each side makes a legal moves in turn until the game and ends in checkmate, stalemate, draw, timeout, surrender.",
					"The ultimate goal of each player is to achieve checkmate.  Checkmate is when you are attacking the opponent's King and there is no move that can free him of an attack.",
					"Lest you can achieve checkmate, stalemate might be your next option.  If you can manage to position your King in such a way as to be safe, but any move you make will check your King, meaning he is attacked."}},
		},
	},
	"Board": {
		Step: []Step{
			{Title: "Battlefield", Width: 70, PosX: 50, PosY: 50, Background: "#00000010", Order: []string{"menu tutor"},
				Black: []string{}, White: []string{}, Log: []string{},
				Text: []string{"The Chexx battlefield is set on a hexagonal grid.  Black spaces, white spaces and neutral spaces."}},
			{Title: "Compass Rose", Width: 70, PosX: 50, PosY: 65, Background: "#00000010", Black: []string{}, White: []string{}, Hilite: [][]string{{"*", "stroke", "#ff0"}},
				Text: []string{"The dead center of the is marked with a gold medallion... the 'compass rose'.", "Notated by '*'."}},

			{Title: "Columns", Width: 25, PosX: 10, PosY: 50, Background: "#999999a0", Black: []string{}, White: []string{},
				Hilite: [][]string{{"a5 a4 a3 a2 a1 a * d d1 d2 d3 d4 d5", "stroke", "#ff0"}, {"a52 a42 a32 a22 b1 b11 c1 c21 c32 c43 c54", "stroke", "#0f0"}, {"f54 f43 f32 f21 f1 e11 e1 d22 d32 d42 d52", "stroke", "#00f"}},
				Text: []string{
					"If you remember from chess, rows run horizontal and columns run vertical.  On the Hexagonal grid, it is bit different... there are still vertical columns.",
					"Highlighted in yellow, is the center vertical column.",
					"In green is the a52~c54 column, and in blue, the f54~d52 column"}},
			{Title: "", Width: 20, PosX: 15, PosY: 15, Background: "#888888a0", Black: []string{}, White: []string{},
				Hilite: [][]string{{"a5 a4 a3 a2 a1 a d d1 d2 d3 d4 d5", "stroke", "#ff0"}, {"b5 b4 b3 b2 b1 b e e1 e2 e3 e4 e5", "stroke", "#08f"}, {"c5 c4 c3 c2 c1 c f f1 f2 f3 f4 f5", "stroke", "#0F4"}},
				Text:   []string{"But columns are a bit more arbitrary on this trilinear topology...", "Lets keep the central vertical column in yellow, but will add the other columns that run through the center."}},

			{Title: "", Width: 20, PosX: 15, PosY: 15, Background: "#888888a0", Black: []string{}, White: []string{},
				Hilite: [][]string{{"a5 a4 a3 a2 a1 a d d1 d2 d3 d4 d5", "stroke", "#ff0"}, {"b5 b4 b3 b2 b1 b e e1 e2 e3 e4 e5", "stroke", "#08f"}, {"c5 c4 c3 c2 c1 c f f1 f2 f3 f4 f5", "stroke", "#0F4"}},
				Text:   []string{"Given this 3 axis format, rows and column coordinate system is less than ideal.", "So we introduce a new coordinate system."}},

			{Title: "", Width: 60, PosX: 50, PosY: 90, Background: "#888888a0", Black: []string{}, White: []string{},
				Hilite: [][]string{
					{"a55 a54 a53 a53 a51 a5 a44 a43 a42 a41 a4 a33 a32 a31 a3 a22 a21 a2 a11 a1 a", "stroke", "#f00"},
					{"b55 b54 b53 b53 b51 b5 b44 b43 b42 b41 b4 b33 b32 b31 b3 b22 b21 b2 b11 b1 b", "stroke", "#ff0"},
					{"c55 c54 c53 c53 c51 c5 c44 c43 c42 c41 c4 c33 c32 c31 c3 c22 c21 c2 c11 c1 c", "stroke", "#0f0"},
					{"d55 d54 d53 d53 d51 d5 d44 d43 d42 d41 d4 d33 d32 d31 d3 d22 d21 d2 d11 d1 d", "stroke", "#0ff"},
					{"e55 e54 e53 e53 e51 e5 e44 e43 e42 e41 e4 e33 e32 e31 e3 e22 e21 e2 e11 e1 e", "stroke", "#00f"},
					{"f55 f54 f53 f53 f51 f5 f44 f43 f42 f41 f4 f33 f32 f31 f3 f22 f21 f2 f11 f1 f", "stroke", "#f0f"}},
				Text: []string{"The battlefield is divided into 6 regions... A through F.", "'A' in red, 'B' in yellow, 'C' in green, 'D' in cyan, 'E' in blue and 'F' in purple."}},

			{Title: "", Width: 70, PosX: 50, PosY: 85, Background: "#888888a0", Order: []string{"menu tutor"},
				Hilite: [][]string{{"a b c d e f", "stroke", "#fff"}},
				Black:  []string{}, White: []string{"Nd32"},
				Text: []string{"Surrounding the compass rose <*> at the center of the board: a,b,c,d,e,f make up the center ring...",
					"so a pawn just above <*> is notated: 'Pa'", "A Bishop, to the upper left of <*> is denoted as 'Bf'"}},

			{Title: "", Width: 70, PosX: 50, PosY: 30, Background: "#888888a0", Order: []string{"menu tutor"},
				Hilite: [][]string{{"d32", "stroke", "#00f"}},
				Black:  []string{}, White: []string{"Nd32"},
				Text:   []string{"The Knight highlighted in red wants to go to <d>, since he is already selected, move him to now..."},
				Answer: [][]string{{"", "That is not how a knight even moves.", "22,60,4,#F66,#000"}, {"d", "Thats the correct position.", "15,60,4,#6F6,#FF0"}, {"e c22 c33 d51 d55 e4 e31 e21", "The Knight can move there, but that is not where he want to be.", "25,70,4,#aa6,#000"}},
				Hint:   []string{"Its next to the compass rose.||22,60,4,#F66,#000"}},

			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#00000010",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},

			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#00000010",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Pawn": {
		Step: []Step{
			{Title: "Pawn", Width: 70, PosX: 50, PosY: 50, Background: "#20a0a070", Order: []string{"menu tutor"},
				Black: []string{"Pa31", "Pf33"}, White: []string{"Pc33", "Pd31"},
				Text: []string{"The lowly, humble pawn.  They are you front line; first to fight and to die for their King.  Alone, they are weak and vulnerable, but united, their formitability dictates the course of battle.",
					"The pawn is the most fundamental aspect of Chexx.  They will be the focus of you strategy since what tactics you employ by the other pieces will be derived from the positioning of your pawn structure."}},
			{Title: "Pawn", Width: 70, PosX: 50, PosY: 30, Background: "#20a0a070",
				Hilite: [][]string{{"d", "stroke", "#0f0"}, {"d11 c11", "stroke", "#f00"}, {"d1", "stroke", "#00f"}},
				Black:  []string{}, White: []string{"Pc22", "Pd1", "Pd21"},
				Text: []string{"If Chess is familiar to you, this will be as you might expect.  A pawn move forward one space at a time.",
					"Although their attacks are only to the forward flanks.  Below, the center pawn selected, shows in green where it can move to if the space is unoccupied.  In red are the spaces to where this pawn can attack an enemy piece.  Only when an enemy piece is present may this pawn take and accupy this space."}},
			{Title: "Pawn", Width: 70, PosX: 50, PosY: 30, Background: "#20a0a070",
				Hilite: [][]string{{"c22 d21", "stroke", "#00f"}},
				Black:  []string{}, White: []string{"Pc22", "Pd1", "Pd21"},
				Text: []string{"Since the pawns are c22 and d21 'attack' the d1 space, occupied by one of their own, that d1 pawn is said to be protected by the c22 and d21 pawns.  What this really means is that if the d1 pawn were taken, there would be potential consequesnce."}},

			{Title: "Pawn Formations", Width: 70, PosX: 50, PosY: 30, Background: "#20a0a070",
				Hilite: [][]string{{"c22 c32 c42 d3 d31", "stroke", "#00f"}},
				Black:  []string{}, White: []string{"Pc22", "Pc32", "Pc42", "Pd3", "Pd31"},
				Text: []string{"Pawns that protect each other in such a formation are said to be a 'pawn chain'.",
					"Such a configuration is powerful.  Used wisely, these formation can cost many times what these individual pawns are worth."}},

			{Title: "", Width: 70, PosX: 50, PosY: 20, Background: "#20a0a070",
				Black: []string{"Pa1", "Pf11", "Pa11", "Pa22", "Pa33"}, White: []string{"Pc22", "Pc32", "Pc42", "Pd2", "Pd31"},
				Text: []string{"You will play White... select the lead pawn, and move him forward one space."},
				Answer: [][]string{{"", "Incorrect.", "40,55,6,#a00,#88F"},
					{"a1 f11 a11 a22 a33", "These pawns don't belong to your army.", "15,45,5,#6F6,#FF0"},
					{"c22 c32 c42 d2 d31", "", "", "~"}, {"c22~c1", "Not that far.", "35,40,4,#6F6,#FF0"},
					{"c22~c11", "Ok, good start.", "25,60,4,#6F6,#FF0", "-"}}},

			{Title: "", Width: 70, PosX: 50, PosY: 20, Background: "#20a0a070",
				Black: []string{"Pa", "Pf11", "Pa11", "Pa22", "Pa33"}, White: []string{"Pc11", "Pc32", "Pc42", "Pd2", "Pd31"},
				Text: []string{"Nice.... the pawn advances forward 1 space,  but now, that lead pawn is all by himself.  Quite unfortunate for a pawn.",
					"Our strategy is to disrupt blacks ranks with a right flank attack.  So Pc32 will take the lead.  Now on a pawn's starting position, he can optionally charge forward two spaces.  From then on, he marches on with his signular focus, one space at a time."},
				Answer: [][]string{{"", "Incorrect.", "40,55,6,#a00,#88F"},
					{"a f11 a11 a22 a33", "These pawns don't belong to your army.", "10,51,5,#000,#FF0"},
					{"c11 c32 c42 d2 d31", "", "", "~"},
					{"c32~c21", "<c32> is meant for greatness!  Be bold an move that extra space.", "10,55,3,#faf,#050"},
					{"c32~c1", "Brilliant!", "70,60,6,#6F6,#FF0", "-"}},
				Hint: []string{"Look at which pawns that could take the lead by moving 2 spaces forward.||20,35,2,#000,#fff"}},

			{Title: "", Width: 70, PosX: 50, PosY: 30, Background: "#20a0a070",
				Black: []string{"Pa", "Pf11", "Pa11", "Pa22", "Pb2"}, White: []string{"Pc11", "Pc1", "Pc42", "Pd2", "Pd31"},
				Text: []string{"That lead pawn will need some back-up.  Give him some support from the left."},
				Answer: [][]string{{"", "Incorrect.", "40,55,6,#a00,#88F"},
					{"a f11 a11 a22 b2", "These pawns don't belong to your army.", "10,50,5,#606,#FF0"},
					{"c11 c1 c42 d2 d31", "", "", "~"},
					{"c42~c31", "<c42> HEY! Get into the fight!", "25,55,4,#6F6,#FF0"},
					{"c42~c2", "Good.", "70,55,5,#6F6,#808", "-"}}},

			{Title: "", Width: 70, PosX: 50, PosY: 15, Background: "#20a0a070",
				Black: []string{"Pf11", "Pa", "Pb", "Pb1", "Pb2"}, White: []string{"Pd11", "Pd", "Pc", "Pb11", "Pb22"},
				Text: []string{"Lets fast foward a few moves.  Time for some bloodshead.",
					"It would have been nice if black would take first at <b11> but he understands that after PbxPb11, and the counter Pb22xPb11, his forces are divides in the center and ours are in tact.",
					"Will will therefore take the initiative.  Take <b> with <b11>."},
				Answer: [][]string{{"", "Incorrect.", "40,55,6,#a00,#88F"},
					{"f11 a b b1 b2", "These pawns don't belong to your army.", "10,50,5,#606,#FF0"},
					{"d11 d c b11 b22", "", "", "~"},
					{"b11~b", "First blood.", "70,55,5,#F00,#055", "-"}}},

			{Title: "Pawn", Width: 70, PosX: 50, PosY: 20, Background: "#20a0a070",
				Black: []string{"Pf11", "Pb", "Pb1", "Pb2"}, White: []string{"Pd11", "Pd", "Pc", "Pb22"},
				Text: []string{"Black takes our pawn.  We will sings songs about his bravery."}},

			{Title: "Pawn", Width: 70, PosX: 50, PosY: 30, Background: "#20a0a070",
				Hilite: [][]string{{"", "stroke", "#00f"}},
				Black:  []string{}, White: []string{""},
				Text: []string{"A pawn, especially, must think as a team.  Valliant a pawn might be, standing alone he is simply a casualty waiting happen.  He needs his fellow soldiers' support.",
					"A single pawn could mean the difference between victory and defeat."}},
		},
	},
	"Spearman": {
		Step: []Step{
			{Title: "Spearman", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Knight": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Bishop": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Rook": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Archer": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Queen": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"King": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Prince": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Princess": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Special": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Promotion": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Forks": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Skewers": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Pins": {
		Step: []Step{
			{Title: "", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Tactics": {
		Step: []Step{
			{Title: "Tactics", Width: 70, PosX: 50, PosY: 50, Background: "#f0202050",
				Black: []string{},
				White: []string{},
				Text:  []string{""}},
		},
	},
	"Unimplemented": {
		Step: []Step{
			{Title: "Unimplemented", Width: 70, PosX: 50, PosY: 50, Background: "#f020f040",
				Black: []string{},
				White: []string{},
				Text:  []string{"Sorry, perhaps check back later for an update."}},
		},
	},
}

type Thing struct {
	Lesson string `json:"lesson"`
}

func Tutorial(w http.ResponseWriter, r *http.Request) {

	// respDump, err := httputil.DumpRequest(r, true)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("RESPONSE:\n%s", string(respDump))
	var s Thing

	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		log.Fatalf("Error happened in JSON Decode. Err: %s", err)
		return
	}

	res, err := json.Marshal(Lessons[s.Lesson])
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	// fmt.Println(string(res))
	utils.CorsHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
