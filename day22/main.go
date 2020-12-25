package main

import "fmt"

var allGames = make(map[string]bool)

func main() {
	deck0 := []int{43, 21, 2, 20, 36, 31, 32, 37, 38, 26, 48, 47, 17, 16, 42, 12, 45, 19, 23, 14, 50, 44, 29, 34, 1}
	deck1 := []int{40, 24, 49, 10, 22, 35, 28, 46, 7, 41, 15, 5, 39, 33, 11, 8, 3, 18, 4, 13, 6, 25, 30, 27, 9}

	// deck0 := []int{9, 2, 6, 3, 1}
	// deck1 := []int{5, 8, 4, 7, 10}

	deck0, deck1 = playGame2(deck0, deck1, 0)

	winner := deck0
	if len(deck1) > 0 {
		winner = deck1
	}

	sum := 0
	for i := 1; i <= len(winner); i++ {
		sum += winner[len(winner)-i] * i
	}

	fmt.Printf("Lösung: %d\n", sum)
}

func playGame(deck0, deck1 []int) ([]int, []int) {
	for len(deck0) > 0 && len(deck1) > 0 {
		if deck0[0] > deck1[0] {
			deck0 = append(deck0, deck0[0], deck1[0])
			copy(deck0, deck0[1:])
			copy(deck1, deck1[1:])
			deck0 = deck0[:len(deck0)-1]
			deck1 = deck1[:len(deck1)-1]
		} else {
			deck1 = append(deck1, deck1[0], deck0[0])
			copy(deck0, deck0[1:])
			copy(deck1, deck1[1:])
			deck0 = deck0[:len(deck0)-1]
			deck1 = deck1[:len(deck1)-1]
		}
	}

	return deck0, deck1
}

func playGame2(deck0, deck1 []int, depth int) ([]int, []int) {
	gameHash := fmt.Sprintf("%v|%v", deck0, deck1)
	if result, ok := allGames[gameHash]; ok {
		fmt.Println("Vorheriges Ergebnis gefunden!")

		if result {
			return []int{1}, []int{}
		}
		return []int{}, []int{1}
	}

	previousRounds := make(map[string]struct{})

	for len(deck0) > 0 && len(deck1) > 0 {
		roundHash := fmt.Sprintf("%v|%v", deck0, deck1)
		if _, ok := previousRounds[roundHash]; ok {
			fmt.Println("Schleife gefunden!")
			allGames[gameHash] = true
			return []int{1}, []int{}
		}

		zeroHasWon := false

		if len(deck0) > deck0[0] && len(deck1) > deck1[0] {
			copy0 := make([]int, deck0[0])
			copy1 := make([]int, deck1[0])
			copy(copy0, deck0[1:])
			copy(copy1, deck1[1:])

			// play a sub game
			copy0, _ = playGame2(copy0, copy1, depth+1)
			zeroHasWon = len(copy0) != 0
		} else {
			zeroHasWon = deck0[0] > deck1[0]
		}

		if zeroHasWon {

			deck0 = append(deck0, deck0[0], deck1[0])
			copy(deck0, deck0[1:])
			copy(deck1, deck1[1:])
			deck0 = deck0[:len(deck0)-1]
			deck1 = deck1[:len(deck1)-1]

		} else {

			deck1 = append(deck1, deck1[0], deck0[0])
			copy(deck0, deck0[1:])
			copy(deck1, deck1[1:])
			deck0 = deck0[:len(deck0)-1]
			deck1 = deck1[:len(deck1)-1]

		}

		previousRounds[roundHash] = struct{}{}
	}

	allGames[gameHash] = len(deck0) != 0
	return deck0, deck1
}
