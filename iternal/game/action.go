package game

type action string

const (
	initialMeld        action = "Начальная ставка"
	addPiece           action = "Выложить фишку"
	removePiece        action = "Убрать фишку"
	replacePiece       action = "Заменить фишку"
	addCombination     action = "Выложить комбинацию"
	concatCombinations action = "Соединить комбинации"
	splitCombination   action = "Разделить комбинацию"
	pass               action = "Пропустить ход"
)

var initialActions [2]action = [2]action{initialMeld, pass}

var mainActions [7]action = [7]action{
	addPiece, removePiece, replacePiece,
	addCombination, concatCombinations, splitCombination, pass,
}
