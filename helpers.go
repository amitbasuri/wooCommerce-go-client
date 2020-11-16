package wooCommerce

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func nextPage(header http.Header, re *regexp.Regexp) string {
	headerLink := header.Get(string(HeaderKeyLink))
	if headerLink == "" {
		return ""
	}

	headerLinks := strings.Split(headerLink, ",")
	var nextLink string
	for _, link := range headerLinks {
		l := re.FindStringSubmatch(link)
		if len(l) > 1 {
			nextLink = l[1]
			break
		}
	}

	if nextLink == "" {
		return ""
	}

	nextUrl, err := url.Parse(nextLink)
	if err != nil {
		return ""
	}

	return nextUrl.Query().Get(string(QueryParamPage))
}
