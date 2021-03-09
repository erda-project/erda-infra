// Author: recallsong
// Email: songruiguo@qq.com

package i18n

import (
	"sort"
	"strconv"
	"strings"
)

// LanguageCode .
type LanguageCode struct {
	Code    string  `json:"code"`
	Quality float32 `json:"quality"`
}

// RestrictedCode .
func (lc *LanguageCode) RestrictedCode() string {
	idx := strings.Index(lc.Code, "-")
	if idx < 0 {
		return lc.Code
	}
	return lc.Code[:idx]
}

// ElaboratedCode .
func (lc *LanguageCode) ElaboratedCode() string {
	idx := strings.Index(lc.Code, "-")
	if idx < 0 {
		return ""
	}
	return lc.Code[idx+1:]
}

// Codes .
func (lc *LanguageCode) Codes() (string, string) {
	parts := strings.SplitN(lc.Code, "-", 3)
	if len(parts) > 1 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}

// String .
func (lc *LanguageCode) String() string {
	if lc.Quality == 1 {
		return lc.Code
	}
	return lc.Code + ";" + strconv.FormatFloat(float64(lc.Quality), 'f', -1, 32)
}

// LanguageCodes .
type LanguageCodes []*LanguageCode

// Len 返回slice长度
func (ls LanguageCodes) Len() int { return len(ls) }

// Less 比较两个位置上的数据
func (ls LanguageCodes) Less(i, j int) bool {
	return ls[i].Quality > ls[j].Quality
}

// Swap 交换两个位置上的数据
func (ls LanguageCodes) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}

// ParseLanguageCode .
func ParseLanguageCode(text string) (list LanguageCodes, err error) {
	for _, item := range strings.Split(text, ",") {
		parts := strings.SplitN(item, ";", 2)
		lc := &LanguageCode{
			Quality: 1,
		}
		if len(parts) > 1 {
			q := strings.TrimSpace(parts[1])
			if len(q) > 0 {
				kv := strings.Split("=", q)
				if len(kv) == 2 && kv[0] == "q" {
					q, err := strconv.ParseFloat(kv[1], 32)
					if err != nil {
						sort.Sort(list)
						return list, err
					}
					lc.Quality = float32(q)
				}
			}
		}
		lc.Code = strings.TrimSpace(parts[0])
		if len(item) > 0 {
			list = append(list, lc)
		}
	}
	sort.Sort(list)
	return list, nil
}
