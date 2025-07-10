package random

import (
	"math/rand"
	"sort"
	"strings"
	"time"

	"log/slog"
)

type charFreq struct {
	char byte
	freq int
}

func NewRandomString(UrlLink string, storage interface {
	AliasExists(alias string) (bool, error)
}, log *slog.Logger) string {
	return GenerateAlias(UrlLink, storage, log)
}

func parceTheLink(UrlLink string) []byte {
	frequency := make(map[byte]int)

	lowerLink := strings.ToLower(UrlLink)
	for i := 0; i < len(lowerLink); i++ {
		char := lowerLink[i]
		if char >= 'a' && char <= 'z' {
			frequency[char]++
		}
	}

	var charFreqs []charFreq
	for char, freq := range frequency {
		charFreqs = append(charFreqs, charFreq{char: char, freq: freq})
	}

	sort.Slice(charFreqs, func(i, j int) bool {
		if charFreqs[i].freq == charFreqs[j].freq {
			return charFreqs[i].char < charFreqs[j].char
		}
		return charFreqs[i].freq > charFreqs[j].freq
	})

	result := make([]byte, 0, 6)
	for i := 0; i < len(charFreqs) && i < 6; i++ {
		result = append(result, charFreqs[i].char)
	}

	return result
}

func GenerateAlias(UrlLink string, storage interface {
	AliasExists(alias string) (bool, error)
}, log *slog.Logger) string {
	const op = "random.GenerateAlias"

	log.Info("начинаем генерацию алиаса",
		slog.String("url", UrlLink),
		slog.String("op", op))

	frequentChars := parceTheLink(UrlLink)

	allChars := make([]byte, 0, 26)
	for char := byte('a'); char <= 'z'; char++ {
		allChars = append(allChars, char)
	}

	availableChars := append(frequentChars, allChars...)

	charMap := make(map[byte]bool)
	uniqueChars := make([]byte, 0)
	for _, char := range availableChars {
		if !charMap[char] {
			charMap[char] = true
			uniqueChars = append(uniqueChars, char)
		}
	}

	log.Info("доступные символы для генерации",
		slog.Int("количество", len(uniqueChars)),
		slog.String("символы", string(uniqueChars)))

	alias := generateRandomAlias(uniqueChars, 6)

	exists, err := storage.AliasExists(alias)
	if err != nil {
		log.Error("ошибка при проверке существования алиаса",
			slog.String("alias", alias),
			slog.String("error", err.Error()),
			slog.String("op", op))
		return alias
	}

	if !exists {
		log.Info("алиас сгенерирован успешно",
			slog.String("alias", alias),
			slog.String("op", op))
		return alias
	}

	log.Info("алиас уже существует, изменяем одну букву",
		slog.String("alias", alias),
		slog.String("op", op))

	modifiedAlias := modifyRandomChar(alias, uniqueChars, storage, log)

	return modifiedAlias
}

func generateRandomAlias(chars []byte, length int) string {
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func modifyRandomChar(alias string, availableChars []byte, storage interface {
	AliasExists(alias string) (bool, error)
}, log *slog.Logger) string {
	const op = "random.modifyRandomChar"

	aliasBytes := []byte(alias)
	maxAttempts := 10

	for attempt := 0; attempt < maxAttempts; attempt++ {
		pos := rand.Intn(len(aliasBytes))
		originalChar := aliasBytes[pos]

		var newChar byte
		for {
			newChar = availableChars[rand.Intn(len(availableChars))]
			if newChar != originalChar {
				break
			}
		}

		aliasBytes[pos] = newChar
		newAlias := string(aliasBytes)

		log.Info("попытка изменения алиаса",
			slog.Int("попытка", attempt+1),
			slog.String("старый_алиас", alias),
			slog.String("новый_алиас", newAlias),
			slog.String("измененная_позиция", string(originalChar)+"->"+string(newChar)),
			slog.String("op", op))

		exists, err := storage.AliasExists(newAlias)
		if err != nil {
			log.Error("ошибка при проверке измененного алиаса",
				slog.String("alias", newAlias),
				slog.String("error", err.Error()),
				slog.String("op", op))
			continue
		}

		if !exists {
			log.Info("успешно изменен алиас",
				slog.String("новый_алиас", newAlias),
				slog.String("op", op))
			return newAlias
		}

		aliasBytes[pos] = originalChar
	}

	log.Warn("не удалось найти уникальный алиас после всех попыток",
		slog.String("исходный_алиас", alias),
		slog.Int("попыток", maxAttempts),
		slog.String("op", op))

	suffix := generateRandomAlias(availableChars, 2)
	finalAlias := alias + suffix

	log.Info("возвращаем алиас с суффиксом",
		slog.String("финальный_алиас", finalAlias),
		slog.String("op", op))

	return finalAlias
}

func Randomizing(charFreqs []charFreq) string {
	rand.Seed(time.Now().UnixNano())

	return ""
}
