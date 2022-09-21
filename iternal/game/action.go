package game

type action int

/*
	Actions

	initialMeld        = "Начальная ставка"
	addPiece           = "Выложить фишку"
	removePiece        = "Убрать фишку"
	replacePiece       = "Заменить фишку"
	addCombination     = "Выложить комбинацию"
	concatCombinations = "Соединить комбинации"
	splitCombination   = "Разделить комбинацию"
	pass               = "Пропустить ход"

*/

const (
	InitialMeld action = iota
	AddPiece
	RemovePiece
	ReplacePiece
	AddCombination
	ConcatCombinations
	SplitCombination
	Pass
)

var initialActions [2]action = [2]action{InitialMeld, Pass}

var mainActions [7]action = [7]action{
	AddPiece, RemovePiece, ReplacePiece,
	AddCombination, ConcatCombinations, SplitCombination, Pass,
}
