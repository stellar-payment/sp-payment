package namegen

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateUUIDFilename(ext string) string {
	return fmt.Sprintf("%s%s", uuid.NewString(), ext)
}

func IntToLetters(number int32) (letters string) {
	number--
	if firstLetter := number / 26; firstLetter > 0 {
		letters += IntToLetters(firstLetter)
		letters += string('A' + number%26)
	} else {
		letters += string('A' + number)
	}

	return
}
