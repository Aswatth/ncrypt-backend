package utils

import (
	"math/rand"
	"strconv"
	"strings"
)

func GeneratePassword(digits bool, uppercase bool, special_char bool, length int) string {

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

	if digits {
		generated_password += strconv.Itoa(rand.Intn(10))
	}

	if uppercase {
		generated_password += string(rune(rand.Intn(26) + 'A'))
	}

	if special_char {
		generated_password += string(special_char_list[rand.Intn(len(special_char_list))])
	}

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

	return generated_password
}
