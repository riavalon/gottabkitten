package gottabkitten

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
)

type wordPair struct {
	Original   string `yaml:"original"`
	Impurroved string `yaml:"impurroved"`
}

func (wp wordPair) checkMatch(word string) bool {
	re := regexp.MustCompile(`\W`)
	sanitizedWord := string(re.ReplaceAll([]byte(word), []byte("")))
	return strings.ToLower(sanitizedWord) == wp.Original
}

type commonWords struct {
	CommonWordPairs []wordPair `yaml:"commonWordPairs"`
	LetterGroups    []wordPair `yaml:"letterGroups"`
}

// ReadSourceFile for checking passed in files
func ReadSourceFile(fileSrc string) []byte {
	file, err := os.OpenFile(fileSrc, os.O_RDONLY, 0755)
	if err != nil {
		panic(fmt.Errorf("I could not open the source file with the provided path. Got this error:\n%v", err))
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("There was an issue when I tried to read the file. Got this error:\n%v", err))
	}

	return data
}

func parseYAML() commonWords {
	var baseWords commonWords
	file, err := os.OpenFile("gottabkitten.yaml", os.O_RDONLY, 0755)
	if err != nil {
		panic(fmt.Errorf("I couldn't open the YAML file to get it's content. Got this error:\n%v", err))
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("I couldn't read the YAML file contents. Got this error:\n%v", err))
	}

	err = yaml.Unmarshal(data, &baseWords)
	if err != nil {
		panic(fmt.Errorf("I couldn't unmarshal the YAML content. Got this error:\n%v", err))
	}

	return baseWords
}

func setCapitalization(origWord, impurrovedWord string) string {
	firstLetter := origWord[0]
	matched, err := regexp.Match(`[A-Z]`, []byte{firstLetter})
	if err != nil {
		panic(fmt.Errorf("There was an issue with regular expression.. error is:\n%v", err))
	}

	if matched {
		return strings.Title(impurrovedWord)
	}
	return impurrovedWord
}

// Impurrove takes text content and makes it punny
func Impurrove(content string) string {
	commonWords := parseYAML()
	words := strings.Split(content, " ")
	var newWords []string
	for _, word := range words {
		newWord := word
		for _, wp := range commonWords.CommonWordPairs {
			if wp.checkMatch(word) {
				newWord = setCapitalization(word, wp.Impurroved)
				break
			}
		}
		if word == newWord {
			// replace any of the matched letterGroups with their impurroved version
			for _, lg := range commonWords.LetterGroups {
				re := regexp.MustCompile(lg.Original)
				if re.Match([]byte(strings.ToLower(word))) {
					replacedVersion := string(re.ReplaceAll([]byte(strings.ToLower(word)), []byte(lg.Impurroved)))
					newWord = setCapitalization(word, replacedVersion)
				}
			}
		}
		newWords = append(newWords, newWord)
	}

	converted := strings.Join(newWords, " ")

	fmt.Println(converted)

	return converted
}
