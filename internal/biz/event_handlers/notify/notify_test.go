package notify

import (
	"go-cs/pkg/pprint"
	"regexp"
	"testing"
)

func Test_cleanHtmlTag(t *testing.T) {
	s := "胡<p 123>(hu)</p> <br /> 12"
	want := "胡(hu) \n 12"

	s2 := cleanHtmlTag(s)

	t.Log(s)
	t.Log(s2 == want)
}

func Test(t *testing.T) {
	re := regexp.MustCompile(`(?i)<span[^>]*>(.*?)</span>`)

	s := "123<span class='c'>abc</span>456"

	allStringIndex := re.FindAllStringIndex(s, -1)
	allStringIndex1 := re.FindAllStringSubmatchIndex(s, -1)

	t.Log(allStringIndex)
	t.Log(allStringIndex1)

	t.Log(s[allStringIndex1[0][0]:allStringIndex1[0][1]])
	t.Log(s[allStringIndex1[0][2]:allStringIndex1[0][3]])
}

func Test_parse(t *testing.T) {
	s := `123<span class="minor-color">abc</span>456`
	rich := parseToImRich(s)
	pprint.Println(rich)
}
