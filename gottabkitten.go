package gottabkitten

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
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

// FIXME(riavalon) temporarily embedding the YAML data inside package
// until I have time to set this up in a not terrible way.
var data = commonWords{
	CommonWordPairs: []wordPair{
		wordPair{Original: "pardon", Impurroved: "purrdon"},
		wordPair{Original: "forget", Impurroved: "furget"},
		wordPair{Original: "forgot", Impurroved: "forgot"},
		wordPair{Original: "attitude", Impurroved: "cattitude"},
		wordPair{Original: "forever", Impurroved: "furever"},
		wordPair{Original: "appauling", Impurroved: "apawling"},
		wordPair{Original: "inferior", Impurroved: "infurior"},
		wordPair{Original: "metaphorically", Impurroved: "metafurically"},
		wordPair{Original: "forward", Impurroved: "furward"},
		wordPair{Original: "perceive", Impurroved: "purrceive"},
		wordPair{Original: "perfect", Impurroved: "purrfect"},
		wordPair{Original: "formiliar", Impurroved: "furmiliar"},
		wordPair{Original: "catastrophe", Impurroved: "catastrophe"},
		wordPair{Original: "feeling", Impurroved: "feline"},
		wordPair{Original: "kidding", Impurroved: "kitten"},
		wordPair{Original: "pause", Impurroved: "pause"},
		wordPair{Original: "literally", Impurroved: "litter-ally"},
		wordPair{Original: "now", Impurroved: "meow"},
		wordPair{Original: "unfortunate", Impurroved: "unfurtunate"},
		wordPair{Original: "possibility", Impurroved: "pawsibility"},
		wordPair{Original: "perhaps", Impurroved: "purrhaps"},
		wordPair{Original: "fortunate", Impurroved: "furtunate"},
		wordPair{Original: "formidable", Impurroved: "furmidable"},
		wordPair{Original: "awful", Impurroved: "clawful"},
	},
	LetterGroups: []wordPair{
		wordPair{Original: "for", Impurroved: "fur"},
		wordPair{Original: "pos", Impurroved: "paws"},
		wordPair{Original: "pr", Impurroved: "purr"},
		wordPair{Original: "pol", Impurroved: "pawl"},
		wordPair{Original: "por", Impurroved: "purr"},
		wordPair{Original: "per", Impurroved: "purr"},
		wordPair{Original: "aw", Impurroved: "claw"},
		wordPair{Original: "liter", Impurroved: "litter"},
		wordPair{Original: "par", Impurroved: "paw"},
		wordPair{Original: "pal", Impurroved: "paw"},
		wordPair{Original: "pur", Impurroved: "purr"},
	},
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

// func parseYAML() commonWords {
// 	var baseWords commonWords
// 	file, err := os.OpenFile("gottabkitten.yaml", os.O_RDONLY, 0755)
// 	if err != nil {
// 		panic(fmt.Errorf("I couldn't open the YAML file to get it's content. Got this error:\n%v", err))
// 	}
// 	defer file.Close()

// 	data, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		panic(fmt.Errorf("I couldn't read the YAML file contents. Got this error:\n%v", err))
// 	}

// 	err = yaml.Unmarshal(data, &baseWords)
// 	if err != nil {
// 		panic(fmt.Errorf("I couldn't unmarshal the YAML content. Got this error:\n%v", err))
// 	}

// 	return baseWords
// }

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
	words := strings.Split(content, " ")
	var newWords []string
	for _, word := range words {
		newWord := word
		for _, wp := range data.CommonWordPairs {
			if wp.checkMatch(word) {
				newWord = setCapitalization(word, wp.Impurroved)
				break
			}
		}
		if word == newWord {
			// replace any of the matched letterGroups with their impurroved version
			for _, lg := range data.LetterGroups {
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
