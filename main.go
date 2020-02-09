package main

import (
	"flag"
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

func main() {
	sourceFile := flag.String("input", "", "Specify an input file with text to read from. If not supplied, expects sourceText")
	sourceText := flag.String("text", "", "Text to be used. If not specified, sourceFile is expected. SourceFile takes precendence.")
	flag.Parse()

	if len(*sourceFile) == 0 && len(*sourceText) == 0 {
		fmt.Println("You must specify either a source file or source text")
		os.Exit(1)
	}

	var textContent string
	if len(*sourceFile) > 0 {
		textContent = string(readSourceFile(*sourceFile))
	} else {
		textContent = *sourceText
	}

	makeTextPunny(textContent)
}

func readSourceFile(fileSrc string) []byte {
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
	file, err := os.OpenFile("gottaBkitten.yaml", os.O_RDONLY, 0755)
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

func makeTextPunny(content string) {
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
}
