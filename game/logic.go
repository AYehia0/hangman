package game

// check if the char exists in the word and return all the correct positions if true
func (g *Game) tryChar(char byte) {
	// TODO: TIP: use hashmap
	var found bool
	for i := 0; i < len(g.word.Word); i++ {
		if g.word.Word[i] == char {
			found = true
			g.currentGuess[i] = char
		}
	}
	if !found {
		g.guesses += 1
	}
}
