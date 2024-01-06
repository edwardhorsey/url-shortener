package base32

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
var base = len(letters) // 36

func Encode(num int) string {
	word := []rune{}

	for num > 0 {
		letter := letters[num%base]
		word = append(word, letter)
		num = num / base
	}

	return string(word)
}

func Decode(str string) int {
	result := 0
	word := []rune(str)

	for i := len(word); i > 0; i-- {
		letter := word[i-1]

		for j := range letters {
			if letter == letters[j] {
				result = result*base + j

				break
			}
		}
	}

	return result
}
