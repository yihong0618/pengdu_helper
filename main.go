package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/bmaupin/go-epub"
)

var BASE_LRC_URL = "https://music.163.com/api/song/lyric/?id="
var BASE_MUSIC_URL = "https://y.music.163.com/m/song?id="

var (
	telegramID    int64
	telegramToken string
	toParseString string
)

func init() {
	flag.StringVar(&toParseString, "nestring", "", "")
	flag.Int64Var(&telegramID, "tgid", 0, "telegram room id")
	flag.StringVar(&telegramToken, "tgtoken", "", "token from telegram")
}

func parseTitle(toParseString string) string {
	m := regexp.MustCompile(`《(.*)》`)
	if m.MatchString(toParseString) {
		return m.FindStringSubmatch(toParseString)[1]
	}
	return "Not Found Title"
}

func parseId(toParseString string) string {
	m := regexp.MustCompile((`\?id=(\d+)`))
	if m.MatchString(toParseString) {
		return m.FindStringSubmatch(toParseString)[1]
	}
	return "Not Found ID"
}

func getLRCByID(id string) string {
	lrc_url := BASE_LRC_URL + id + "&lv=-1"
	resp, err := http.Get(lrc_url)
	if err != nil {
		log.Fatal(err)
	}
	var generic map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&generic)
	if err != nil {
		log.Fatal(err)
	}
	lrc := generic["lrc"].(map[string]interface{})
	return lrc["lyric"].(string)
}

func main() {
	flag.Parse()
	title := parseTitle(toParseString)
	musicId := parseId(toParseString)
	content := getLRCByID(musicId)
	m := regexp.MustCompile(`(\[(.*?)\])`)
	content = m.ReplaceAllString(content, "")
	contentList := strings.Split(content, "\n")
	e := epub.NewEpub(title)

	e.SetAuthor("hongyi_bot")
	body := "<h1>" + title + "</h1>"
	link := BASE_MUSIC_URL + musicId
	body += "<a href=\"" + link + "\">" + link + "</a>"
	for i := 0; i < len(contentList); i++ {
		if len(contentList[i]) > 0 {
			body = body + ("<p>" + contentList[i] + "</p>")
		}
	}
	// Add the section
	e.AddSection(body, "", "", "")
	err := e.Write("music_lrc.epub")
	if err != nil {
		log.Fatal(err)
	}
}
