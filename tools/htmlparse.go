package tools

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func GetIconFromMyApp(packageName string) (string, error) {
	doc, err := goquery.NewDocument(fmt.Sprintf("https://sj.qq.com/appdetail/%s", packageName))
	if err != nil {
		return "", err
	}
	se := doc.Find("img.jsx-2463046046")
	val, ex := se.Attr("src")
	if ex {
		return val, nil
	}
	return "", fmt.Errorf("Can not get Icon.[%s]", packageName)
}
