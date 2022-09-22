package game

// Action code
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
	// Initial meld
	InitialMeld action = iota
	// Add piece to the existing combination
	AddPiece
	// Remove piece from the existing combination
	RemovePiece
	// Replace a piece from the hands with a piece from the combination
	ReplacePiece
	// Add new combination to the game field
	AddCombination
	// Concat two existing combinations
	ConcatCombinations
	// Split the existing combination
	SplitCombination
	// Pass
	Pass
)

// Actions available for initial meld stage
var initialActions [2]action = [2]action{InitialMeld, Pass}

// Actions available for main game stage
var mainActions [7]action = [7]action{
	AddPiece, RemovePiece, ReplacePiece,
	AddCombination, ConcatCombinations, SplitCombination, Pass,
}
