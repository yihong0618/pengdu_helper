package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/bmaupin/go-epub"
)

var BASE_LRC_URL = "https://music.163.com/api/song/lyric/?id="
var BASE_MUSIC_URL = "https://y.music.163.com/m/song?id="

var (
	telegramID    int64
	telegramToken string
	neString      string
	fileString    string
	withTime      bool
)

func init() {
	flag.StringVar(&neString, "nestring", "", "")
	flag.StringVar(&fileString, "filestring", "", "")
	flag.Int64Var(&telegramID, "tgid", 0, "telegram room id")
	flag.StringVar(&telegramToken, "tgtoken", "", "token from telegram")
	flag.BoolVar(&withTime, "withtime", false, "if with time line or not")
}

func parseTitle(toParseString string) string {
	m := regexp.MustCompile(`《(.*)》`)
	if m.MatchString(toParseString) {
		return m.FindStringSubmatch(toParseString)[1]
	}
	return ""
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

// for generate srt epub or ass epub for movie or drama also return the name
func generate_ass(file_name string) ([]string, string) {
	file_string_list := strings.Split(file_name, ".")
	if len(file_string_list) != 2 {
		fmt.Println("file name error")
		os.Exit(1)
	}
	file_type := file_string_list[1]
	reader, err := os.Open(file_name)

	if err != nil {
		fmt.Println(err)
	}
	defer func() { _ = reader.Close() }()
	subtitles := astisub.NewSubtitles()
	switch file_type {
	case "ass", "ssa":
		subtitles, err = astisub.ReadFromSSA(reader)
		if err != nil {
			fmt.Println(err)
		}
	case "srt":
		subtitles, err = astisub.ReadFromSRT(reader)
		if err != nil {
			fmt.Println(err)
		}
	}
	// not support just raise
	string_list := []string{}
	for _, item := range subtitles.Items {
		time_string := item.StartAt.String()
		for _, line := range item.Lines {
			for _, lineItem := range line.Items {
				if len(lineItem.Text) > 0 {
					if withTime {
						string_list = append(string_list, time_string)
					}
					string_list = append(string_list, lineItem.Text)
				}
			}
		}
	}
	return string_list, file_string_list[0]
}

// is dir
func isDirectory(path_name string) (bool, error) {
	fileInfo, err := os.Stat(path_name)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}

func make_body(assContentList []string) string {
	body := ""
	for i := 0; i < len(assContentList); i++ {
		if len(assContentList[i]) > 0 {
			body = body + ("<p>" + assContentList[i] + "</p>")
		}
	}
	return body
}

func main() {
	flag.Parse()
	// if we can not pares title, we will ignore it
	title := ""
	if neString != "" {
		title = parseTitle(neString)
		musicId := parseId(neString)
		content := getLRCByID(musicId)
		m := regexp.MustCompile(`(\[(.*?)\])`)
		content = m.ReplaceAllString(content, "")
		lrcContentList := strings.Split(content, "\n")
		e := epub.NewEpub(title)

		e.SetAuthor("hongyi_bot")
		body := "<h1>" + title + "</h1>"
		link := BASE_MUSIC_URL + musicId
		body += "<a href=\"" + link + "\">" + link + "</a>"
		for i := 0; i < len(lrcContentList); i++ {
			if len(lrcContentList[i]) > 0 {
				body = body + ("<p>" + lrcContentList[i] + "</p>")
			}
		}
		// Add the section
		e.AddSection(body, "", "", "")
		err := e.Write(musicId + ".epub")
		if err != nil {
			log.Fatal(err)
		}
	}
	if fileString != "" {
		info, err := os.Stat(fileString)
		if err != nil {
			fmt.Println(fileString + "is not a file or dir")
			return
		}
		// TODO support it in the future
		if info.IsDir() {
			fmt.Println("not support for dir for now")
			return
		}
		assBodyList, assName := generate_ass(fileString)
		assBody := make_body(assBodyList)
		e := epub.NewEpub(fileString)
		e.SetAuthor("hongyi_bot")
		e.AddSection(assBody, "", "", "")
		assName = strings.Replace(assName, " ", "_", -1)
		err = e.Write(assName + ".epub")
		if err != nil {
			log.Fatal(err)
		}
	}
}
