package utils

import (
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	readability "github.com/go-shiori/go-readability"
)

var blockedExt = []string{
	"jpg", "jpeg", "png", "gif", "pdf", "zip", "rar",
}

func IsAllowedExtension(url string) bool {
	ext := strings.ToLower(path.Ext(url))
	if len(ext) > 0 {
		ext = ext[1:] // ตัดจุดนำหน้า
	}
	for _, b := range blockedExt {
		if ext == b {
			return false
		}
	}
	return true
}

// IsHTMLURL ใช้ HEAD ตรวจสอบ Content-Type ต้องมี "text/html"
func IsHTMLURL(url string) (bool, error) {
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	return strings.Contains(ct, "text/html"), nil
}

func GetTextFromURL(url string) (string, error) {
	// // Parse URL เป็นบทความ พร้อม timeout
	// ok, err := IsHTMLURL(url)
	// if err != nil {
	// 	return "", err
	// }
	// if !ok {
	// 	return "", errors.New("skipped: content is not HTML")
	// }

	// // ตรวจสอบว่า URL มี extension ที่ไม่อนุญาตหรือไม่
	// if !IsAllowedExtension(url) {
	// 	return "", errors.New("skipped: URL has blocked extension")
	// }

	article, err := readability.FromURL(url, 60*time.Second)
	if err != nil {
		return "", err
	}
	// คืนแต่ข้อความล้วนๆ ไม่รวมโค้ด HTML
	return article.TextContent, nil
}

var wsRegexp = regexp.MustCompile(`\s+`)

// NormalizeSpace ตัด leading/trailing whitespace
// และแทนที่ทุกชุด whitespace (space, \n, \t ฯลฯ) ด้วย single space
func NormalizeSpace(s string) string {
	s = strings.TrimSpace(s)
	return wsRegexp.ReplaceAllString(s, " ")
}

var sentenceSplitter = regexp.MustCompile(`(?m)([^\.!?]+[\.!?])`)

// SummarizeText returns up to maxSentences sentences from s.
func SummarizeText(s string, maxSentences int) string {
	matches := sentenceSplitter.FindAllString(s, maxSentences)
	if len(matches) == 0 {
		return s
	}
	return strings.Join(matches, " ")
}
