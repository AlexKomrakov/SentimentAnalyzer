package main
import (
	"strings"
	"github.com/kljensen/snowball"
	"os"
	"bufio"
	"strconv"
	"fmt"
	"time"
)

var (
	text = `
Чёртовы расисты из 20th Century Fox Television, принципиально не желавшие видеть на главной роли в продолжении сериала «24 часа» белого актёра (что, разумеется, тоже расизм, но узаконенный), сделали ход чёрным конём и в качестве наследника Кифера Сазерленда выбрали Кори Хоукинса — звезду фильма «Голос улиц» (Straight Outta Compton).

Нате, полюбуйтесь на эту будущую грозу террористов из сериала «24: Наследие»:


Персонажа с внешностью студента колледжа (в лучшем случае) будут звать Эрик Картер, и это, конечно, сплошное недоразумение. Если, как утверждает народное творчество, «Джек Бауэр» в переводе с арабского означает «Мне пиздец», то «Эрик Картер» тянет максимум на «Какой чудесный день для теракта».

Несмотря на это, сценаристы и режиссёры рьяно попытаются изобразить Картера не пацаном, который обосрался бы после первого же допроса в исполнении Джека Бауэра, а героем войны, вернувшимся домой с неприятностями на хвосте. По возвращении он обращается за помощью в новое CTU, чтобы спасти свою жизнь (уже неплохое начало для «наследника» Бауэра) и заодно остановить крупномасштабную террористическую атаку на США.

Ещё в «Наследии» ожидается сильный женский персонаж — не менее важный, чем главный герой. В глаза не видевшие «24 часа» представители Fox или журналисты с серьёзными нарушениями слуха утверждают, что это прям новое слово во франчайзе — раньше, дескать, Джек бегал по городу один-одинёшенек и не мог опереться на крепкое девичье плечо. Про Рене Уокер и Кейт Морган (а также кучу других чудесных женских персонажей вроде Мишель Десслер) вышеупомянутые дурачки явно позабыли. Действительно — зачем знать хоть что-то о своём сериале? Политкорректность соблюдена — это куда важнее.

Радует одно — Кифер Сазерленд к данному проекту вроде бы не имеет никакого отношения. И мы, пожалуй, побережём свою психику. Спасибо большое, оставьте это добро себе.
`
	separator = " "
	cutset = ",.\n\r\t "
	dictionary = readFile("dictionary/combined_stemmed.csv")
)

type SentimentWord struct {
	Word string
	Tone int
}

func (s *SentimentWord) FromCSV(csv string) {
	splitted := strings.SplitN(csv, ",", 2)
	s.Word = splitted[0]
	tone, err := strconv.Atoi(splitted[1])
	check(err)
	s.Tone = tone
}

func (s *SentimentWord) CSV() string {
	return s.Word + "," + strconv.Itoa(s.Tone)
}

func (s *SentimentWord) GetTone() {
	index := SliceIndex(len(dictionary), func(i int) bool { return dictionary[i].Word == s.Word })
	if index != -1 {
		s.Tone = dictionary[index].Tone
	}
}

func main() {

	start := time.Now()

	text_lower := strings.ToLower(text)
	tokens := strings.Split(text_lower, separator)
	for key, token := range tokens {
		tokens[key] = strings.Trim(token, cutset)
	}

	text_words := make([]SentimentWord, 0)
	result_tone := 0
	for key, token := range tokens {
		stemmed, err := snowball.Stem(token, "russian", true)
		check(err)
		tokens[key] = stemmed
		word := SentimentWord{Word: stemmed}
		word.GetTone()
		result_tone += word.Tone
		text_words = append(text_words, word)
	}

	fmt.Println(result_tone)
	fmt.Println(time.Since(start))

	//combined := readFile("dictionary/combined_stemmed.csv")
	//stemmed := FindDuplicates(combined)
	//writeFile("dictionary/combined_stemmed.csv", stemmed)
}

func readFile(filename string) []SentimentWord {
	file, err := os.Open(filename)
	check(err)

	defer file.Close()

	words := make([]SentimentWord, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		word := SentimentWord{}
		word.FromCSV(line)
		stemmed, err := snowball.Stem(word.Word, "russian", true)
		check(err)
		word.Word = stemmed
		words = append(words, word)
	}

	err = scanner.Err();
	check(err)

	return words
}

func writeFile(filename string, words []SentimentWord) {
	f, err := os.Create(filename)
	check(err)

	defer f.Close()

	w := bufio.NewWriter(f)
	for _, word := range words {
		_, err := w.WriteString(word.CSV() + "\n")
		check(err)
	}
	w.Flush()
}

func FindDuplicates(words []SentimentWord) (words_filtered []SentimentWord) {
	for _, word := range words {
		index := SliceIndex(len(words_filtered), func(i int) bool { return words_filtered[i].Word == word.Word })
		if index == -1 {
			words_filtered = append(words_filtered, word)
		} else {
			fmt.Println(word, words_filtered[index])
		}
	}
	fmt.Println(len(words) - len(words_filtered))

	return
}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}