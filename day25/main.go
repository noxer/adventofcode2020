package main

import "fmt"

func main() {
	doorPublic := 16616892
	cardPublic := 14505727
	//cardPublic := 5764801
	//doorPublic := 17807724

	doorIterations := 0
	cardIterations := 0

	value := 1
	for i := 1; ; i++ {
		value *= 7
		value %= 20201227

		if value == doorPublic && doorIterations == 0 {
			doorIterations = i
			if cardIterations != 0 {
				break
			}
		}
		if value == cardPublic && cardIterations == 0 {
			cardIterations = i
			if doorIterations != 0 {
				break
			}
		}
	}

	fmt.Printf("Tür: %d\nKarte: %d\n", doorIterations, cardIterations)

	value = 1
	for i := 0; i < cardIterations; i++ {
		value *= doorPublic
		value %= 20201227
	}

	fmt.Printf("Schlüssel: %d\n", value)
}
