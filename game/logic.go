package game

// check if the char exists in the word and return all the correct positions if true
func (g *Game) tryChar(char string) {
	// TODO: TIP: use hashmap
	var found bool
	for i := 0; i < len(g.word.Word); i++ {
		if string(g.word.Word[i]) == char {
			found = true
			g.currentGuess[i] = char
		}
	}
	if !found {
		g.guesses += 1
	}
}

// show blank chars
func (g *Game) showWord() string {
	var word string

	for i := 0; i < len(g.currentGuess); i++ {
		if g.currentGuess[i] == "" {
			word += "_"
		} else {
			word += g.currentGuess[i]
		}
	}

	return word
}
