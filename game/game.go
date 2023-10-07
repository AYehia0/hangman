package game

import "github.com/charmbracelet/ssh"

/*
- API to get the words : https://www.wordgamedb.com/api/v1/words/random
- Interface -> CheckWord, UpdateUI
*/

/*

 Hangman
 Word: c r _ _ d
  _______
 |       0
 |      /|\
 |       |
 |      /
 |
---

Keyboard:
| c | r | a | i | d |   |   |   |
|   |   |   |   |   |   |   |   |
|   |   |   |   |   |   |   |   |

*/
// creed : : maxGuesses : 6, guesses : 4
// keyboard : {c, r, d, a} --> coloured
// c : 0
// r : 1
// e : [2, 3]
// d : 4

// identifiy the keyboard location for exmaple(a: 0, selected: true) as selected or not.
type selected int

type Game struct {
	answer       string
	won          bool
	session      ssh.Session
	guesses      int
	currentGuess string
	maxGuesses   int
	keyboard     map[string]selected
}
