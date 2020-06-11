package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func writeStackInfo(key string, value string) {

	in := fmt.Sprintf("%s"+"="+"%s\n", key, value)

	f, err := os.OpenFile("stack.info", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(in)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func getStackInfo(value string) string {

	a := value + "="

	content, err := ioutil.ReadFile("stack.info")
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	pos := strings.LastIndex(text, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(text) {
		return ""
	}
	return text[adjustedPos:len(text)]
}

func pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func downloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
