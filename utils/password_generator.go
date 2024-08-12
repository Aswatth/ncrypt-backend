package utils

import (
	"math/rand"
	"ncrypt/utils/logger"
	"strconv"
	"strings"
)

func GeneratePassword(digits bool, uppercase bool, special_char bool, length int) string {
	logger.Log.Println("Generating password")
	generated_password := ""

	special_char_list := []rune{'!', '@', '#', '$', '%', '^', '&', '*'}

	character_matrix := [][]rune{
		{'a', 'b', 'c', 'd', 'e', 'f'},
		{'g', 'h', 'i', 'd', 'e', 'f'},
		{'m', 'n', 'o', 'p', 'q', 'r'},
		{'s', 't', 'u', 'v', 'w', 'x'},
		{'y', 'z', '0', '1', '2', '3'},
		{'4', '5', '6', '7', '8', '9'},
	}

	logger.Log.Println("Setting initial characters as per constraints")
	if digits {
		generated_password += strconv.Itoa(rand.Intn(10))
	}

	if uppercase {
		generated_password += string(rune(rand.Intn(26) + 'A'))
	}

	if special_char {
		generated_password += string(special_char_list[rand.Intn(len(special_char_list))])
	}

	logger.Log.Println("Randomly generating characters to match required password length")
	for len(generated_password) <= length {
		i := rand.Intn(len(character_matrix))
		j := rand.Intn(len(character_matrix[i]))

		should_be_special := rand.Intn(2)

		//Random character is a digit
		if character_matrix[i][j] >= '0' && character_matrix[i][j] <= '9' {
			//Would map to special character above respective digit in keyboard
			if should_be_special == 1 && special_char && character_matrix[i][j] >= '1' && character_matrix[i][j] <= '8' {
				generated_password += string(special_char_list[(character_matrix[i][j]-'0')-1])
			} else if digits {
				generated_password += string(character_matrix[i][j])
			}
		} else if character_matrix[i][j] >= 'a' && character_matrix[i][j] <= 'z' { //Random character is an alphabet
			if should_be_special == 1 && uppercase {
				generated_password += strings.ToUpper(string(character_matrix[i][j]))
			} else {
				generated_password += string(character_matrix[i][j])
			}
		}
	}

	logger.Log.Println("Shuffling characters")
	return shuffule(generated_password, rand.Intn(len(generated_password)))
}

func shuffule(generated_password string, shuffle_count int) string {
	generated_password_rune := []rune(generated_password)

	//Shuffle first three characters which was generated as per constraints
	random_index := rand.Intn(len(generated_password_rune))
	generated_password_rune[0], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[0]
	generated_password_rune[0], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[0]
	generated_password_rune[0], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[0]

	random_index = rand.Intn(len(generated_password_rune))
	generated_password_rune[1], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[1]
	generated_password_rune[1], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[1]
	generated_password_rune[1], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[1]

	random_index = rand.Intn(len(generated_password_rune))
	generated_password_rune[2], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[2]
	generated_password_rune[2], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[2]
	generated_password_rune[2], generated_password_rune[random_index] = generated_password_rune[random_index], generated_password_rune[2]

	//Random shuffle
	for range shuffle_count {
		random_index_1 := rand.Intn(len(generated_password_rune))
		random_index_2 := rand.Intn(len(generated_password_rune))

		generated_password_rune[random_index_1], generated_password_rune[random_index_2] = generated_password_rune[random_index_2], generated_password_rune[random_index_1]
	}

	return string(generated_password_rune)
}
