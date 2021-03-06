package crawler

import (
	"digger/models"
	"digger/utils"
	"github.com/antchfx/htmlquery"
	"github.com/hetianyi/gox/logger"
	"regexp"
	"strings"
)

var (
	charsetRegex = regexp.MustCompile(".*charset=(.*)")
)

// 非utf8页面编码处理
func handleEncoding(cxt *models.Context) {
	doc, err := parseXpathDocument(cxt, false)
	if err != nil {
		logger.Debug("handleEncoding Error: cannot parse xpath document")
		return
	}
	list, _ := htmlquery.QueryAll(doc, "//meta[@charset]")

	if len(list) > 0 {
		node := list[0]
		charset := strings.TrimSpace(htmlquery.InnerText(node))
		if charset != "" {
			reParseXpathDocWithCharset(cxt, charset)
		}
		return
	}

	list, _ = htmlquery.QueryAll(doc, "//meta[contains(@content, \"charset=\")]")
	if len(list) > 0 {
		node := list[0]
		val := ""
		for _, a := range node.Attr {
			if strings.ToLower(a.Key) == "content" {
				val = strings.TrimSpace(a.Val)
			}
		}
		charset := charsetRegex.ReplaceAllString(val, "$1")
		if charset != "" {
			reParseXpathDocWithCharset(cxt, charset)
		}
		return
	}
}

func reParseXpathDocWithCharset(cxt *models.Context, charset string) {
	cxt.ResponseData = utils.Trans2UTF8(charset, cxt.ResponseData)
	parseXpathDocument(cxt, true)
}
