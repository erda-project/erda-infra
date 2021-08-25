// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package strutil is string util package
package strutil

import (
	"bytes"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Trim trim `s`'s prefix and suffix. If `cutset` not specified, `cutset` = space.
//
// Trim("trim ") => "trim"
//
// Trim(" this  ") => "this"
//
// Trim("athisb", "abs") => "this"
func Trim(s string, cutset ...string) string {
	if len(cutset) == 0 {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, cutset[0])
}

// TrimSuffixes trim `s`'s suffixes.
//
// TrimSuffixes("test.go", ".go") => "test"
//
// TrimSuffixes("test.go", ".md", ".go", ".sh") => "test"
//
// TrimSuffixes("test.go.tmp", ".go", ".tmp") => "test.go"
func TrimSuffixes(s string, suffixes ...string) string {
	originLen := len(s)
	for i := range suffixes {
		trimmed := strings.TrimSuffix(s, suffixes[i])
		if len(trimmed) != originLen {
			return trimmed
		}
	}
	return s
}

// TrimPrefixes trim `s`'s prefixes.
//
// TrimPrefixes("/tmp/file", "/tmp") => "/file"
//
// TrimPrefixes("/tmp/tmp/file", "/tmp", "/tmp/tmp") => "/tmp/file"
func TrimPrefixes(s string, prefixes ...string) string {
	originLen := len(s)
	for i := range prefixes {
		trimmed := strings.TrimPrefix(s, prefixes[i])
		if len(trimmed) != originLen {
			return trimmed
		}
	}
	return s
}

// TrimSlice is the slice version of Trim.
//
// TrimSlice([]string{"trim ", " trim", " trim "}) => []string{"trim", "trim", "trim"}
func TrimSlice(ss []string, cutset ...string) []string {
	r := make([]string, len(ss))
	for i := range ss {
		r[i] = Trim(ss[i], cutset...)
	}
	return r
}

// TrimSliceSuffixes is the slice version of TrimSuffixes.
//
// TrimSliceSuffixes([]string{"test.go", "test.go.tmp"}, ".go", ".tmp") => []string{"test", "test.go"}
func TrimSliceSuffixes(ss []string, suffixes ...string) []string {
	r := make([]string, len(ss))
	for i := range ss {
		r[i] = TrimSuffixes(ss[i], suffixes...)
	}
	return r
}

// TrimSlicePrefixes is the slice version of TrimPrefixes.
//
// TrimSlicePrefixes([]string{"/tmp/file", "/tmp/tmp/file"}, "/tmp", "/tmp/tmp") => []string{"/file", "/tmp/file"}
func TrimSlicePrefixes(ss []string, prefixes ...string) []string {
	r := make([]string, len(ss))
	for i := range ss {
		r[i] = TrimPrefixes(ss[i], prefixes...)
	}
	return r
}

// HasPrefixes judge if `s` have at least one of elem in `prefixes` as prefix.
//
// HasPrefixes("asd", "ddd", "uuu") => false
//
// HasPrefixes("asd", "sd", "as") => true
//
// HasPrefixes("asd", "asd") => true
func HasPrefixes(s string, prefixes ...string) bool {
	for i := range prefixes {
		if strings.HasPrefix(s, prefixes[i]) {
			return true
		}
	}
	return false
}

// HasSuffixes judge if `s` have at least one of elem in `suffixes` as suffix.
//
// HasSuffixes("asd", "ddd", "d") => true
//
// HasSuffixes("asd", "sd") => true
//
// HasSuffixes("asd", "iid", "as") => false
func HasSuffixes(s string, suffixes ...string) bool {
	for i := range suffixes {
		if strings.HasSuffix(s, suffixes[i]) {
			return true
		}
	}
	return false
}

var (
	collapseWhitespaceRegex = regexp.MustCompile("[ \t\n\r]+")
)

// CollapseWhitespace replace continues space(collapseWhitespaceRegex) to one blank.
//
// CollapseWhitespace("only    one   space") => "only one space"
//
// CollapseWhitespace("collapse \n   all \t  sorts of \r \n \r\n whitespace") => "collapse all sorts of whitespace"
func CollapseWhitespace(s string) string {
	return collapseWhitespaceRegex.ReplaceAllString(s, " ")
}

// Center centering `s` according to total length.
//
// Center("a", 5) => "  a  "
//
// Center("ab", 5) => "  ab "
//
// Center("abc", 1) => "abc"
func Center(s string, length int) string {
	minus := length - len(s)
	if minus <= 0 {
		return s
	}
	right := minus / 2
	mod := minus % 2
	return strings.Join([]string{strings.Repeat(" ", right+mod), s, strings.Repeat(" ", right)}, "")
}

// Split split `s` by `sep`. If `omitEmptyOpt`=true, ignore empty string.
//
// Split("a|bc|12||3", "|") => []string{"a", "bc", "12", "", "3"}
//
// Split("a|bc|12||3", "|", true) => []string{"a", "bc", "12", "3"}
//
// Split("a,b,c", ":") => []string{"a,b,c"}
func Split(s string, sep string, omitEmptyOpt ...bool) []string {
	var omitEmpty bool
	if len(omitEmptyOpt) > 0 && omitEmptyOpt[0] {
		omitEmpty = true
	}
	parts := strings.Split(s, sep)
	if !omitEmpty {
		return parts
	}
	result := []string{}
	for _, v := range parts {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

var (
	linesRegex = regexp.MustCompile("\r\n|\n|\r")
)

// Lines split `s` by newline. If `omitEmptyOpt`=true, ignore empty string.
//
// Lines("abc\ndef\nghi") => []string{"abc", "def", "ghi"}
//
// Lines("abc\rdef\rghi") => []string{"abc", "def", "ghi"}
//
// Lines("abc\r\ndef\r\nghi\n") => []string{"abc", "def", "ghi", ""}
//
// Lines("abc\r\ndef\r\nghi\n", true) => []string{"abc", "def", "ghi"}
func Lines(s string, omitEmptyOpt ...bool) []string {
	lines := linesRegex.Split(s, -1)
	if len(omitEmptyOpt) == 0 || !omitEmptyOpt[0] {
		return lines
	}
	r := []string{}
	for i := range lines {
		if lines[i] != "" {
			r = append(r, lines[i])
		}
	}
	return r
}

// Join see also strings.Join,
// If omitEmptyOpt = true, ignore empty string inside `ss`.
func Join(ss []string, sep string, omitEmptyOpt ...bool) string {
	if len(omitEmptyOpt) == 0 || !omitEmptyOpt[0] {
		return strings.Join(ss, sep)
	}
	r := []string{}
	for i := range ss {
		if ss[i] != "" {
			r = append(r, ss[i])
		}
	}
	return strings.Join(r, sep)
}

// Contains check if `s` contains one of `substrs`.
//
// Contains("test contains.", "t c", "iii")  => true
//
// Contains("test contains.", "t cc", "test  ") => false
//
// Contains("test contains.", "iii", "uuu", "ont") => true
func Contains(s string, substrs ...string) bool {
	for i := range substrs {
		if strings.Contains(s, substrs[i]) {
			return true
		}
	}
	return false
}

// Equal judge whether `s` is equal to `other`. If ignorecase=true, judge without case.
//
// Equal("aaa", "AAA") => false
//
// Equal("aaa", "AaA", true) => true
func Equal(s, other string, ignorecase ...bool) bool {
	if len(ignorecase) == 0 || !ignorecase[0] {
		return strings.Compare(s, other) == 0
	}
	return strings.EqualFold(s, other)
}

// Map apply each funcs to each elem of `ss`.
//
// Map([]string{"1", "2", "3"}, func(s string) string {return Concat("X", s)}) => []string{"X1", "X2", "X3"}
//
// Map([]string{"Aa", "bB", "cc"}, ToLower, Title) => []string{"Aa", "Bb", "Cc"}
func Map(ss []string, fs ...func(s string) string) []string {
	r := []string{}
	for i := range ss {
		r = append(r, ss[i])
	}
	r2 := []string{}
	for _, f := range fs {
		for i := range r {
			r2 = append(r2, f(r[i]))
		}
		r = r2[:]
		r2 = []string{}
	}
	return r
}

// DedupSlice return a slice without repeating elements, and the elements are sorted in the order of their first appearance.
// If omitEmptyOpt = true, ignore empty string.
//
// DedupSlice([]string{"c", "", "b", "a", "", "a", "b", "c", "", "d"}) => []string{"c", "", "b", "a", "d"}
//
// DedupSlice([]string{"c", "", "b", "a", "", "a", "b", "c", "", "d"}, true) => []string{"c", "b", "a", "d"}
func DedupSlice(ss []string, omitEmptyOpt ...bool) []string {
	var omitEmpty bool
	if len(omitEmptyOpt) > 0 && omitEmptyOpt[0] {
		omitEmpty = true
	}
	result := make([]string, 0, len(ss))
	m := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		if s == "" && omitEmpty {
			continue
		}
		if _, ok := m[s]; ok {
			continue
		}
		result = append(result, s)
		m[s] = struct{}{}
	}
	return result
}

// DedupUint64Slice return a slice without repeating elements, and the elements are sorted in the order of their first appearance.
// If omitZeroOpt = true, ignore zero value elem.
//
// DedupUint64Slice([]uint64{3, 3, 1, 2, 1, 2, 3, 3, 2, 1, 0, 1, 2}) => []uint64{3, 1, 2, 0}
//
// DedupUint64Slice([]uint64{3, 3, 1, 2, 1, 2, 3, 3, 2, 1, 0, 1, 2}, true) => []uint64{3, 1, 2}
func DedupUint64Slice(ii []uint64, omitZeroOpt ...bool) []uint64 {
	var omitZero bool
	if len(omitZeroOpt) > 0 && omitZeroOpt[0] {
		omitZero = true
	}
	result := make([]uint64, 0, len(ii))
	m := make(map[uint64]struct{}, len(ii))
	for _, i := range ii {
		if i == 0 && omitZero {
			continue
		}
		if _, ok := m[i]; ok {
			continue
		}
		result = append(result, i)
		m[i] = struct{}{}
	}
	return result
}

// DedupInt64Slice ([]int64{3, 3, 1, 2, 1, 2, 3, 3, 2, 1, 0, 1, 2}, true) => []int64{3, 1, 2} .
func DedupInt64Slice(ii []int64, omitZeroOpt ...bool) []int64 {
	var omitZero bool
	if len(omitZeroOpt) > 0 && omitZeroOpt[0] {
		omitZero = true
	}
	result := make([]int64, 0, len(ii))
	m := make(map[int64]struct{}, len(ii))
	for _, i := range ii {
		if i == 0 && omitZero {
			continue
		}
		if _, ok := m[i]; ok {
			continue
		}
		result = append(result, i)
		m[i] = struct{}{}
	}
	return result
}

// IntersectionUin64Slice return the intersection of two uint64 slices, complexity O(m * n), to be optimized.
//
// IntersectionUin64Slice([]uint64{3, 1, 2, 0}, []uint64{0, 3}) => []uint64{3, 0}
//
// IntersectionUin64Slice([]uint64{3, 1, 2, 1, 0}, []uint64{1, 2, 0}) => []uint64{1, 2, 1, 0}
func IntersectionUin64Slice(s1, s2 []uint64) []uint64 {
	if len(s1) == 0 {
		return nil
	}
	if len(s2) == 0 {
		return s1
	}
	var result []uint64
	for _, i := range s1 {
		for _, j := range s2 {
			if i == j {
				result = append(result, i)
				break
			}
		}
	}
	return result
}

// IntersectionInt64Slice return the intersection of two int64 slices, complexity O(m * log(m)) .
//
// IntersectionInt64Slice([]int64{3, 1, 2, 0}, []int64{0, 3}) => []int64{3, 0}
//
// IntersectionInt64Slice([]int64{3, 1, 2, 1, 0}, []int64{1, 2, 0}) => []int64{1, 2, 1, 0}
func IntersectionInt64Slice(s1, s2 []int64) []int64 {
	m := make(map[int64]bool)
	nn := make([]int64, 0)
	for _, v := range s1 {
		m[v] = true
	}
	for _, v := range s2 {
		if _, ok := m[v]; ok {
			nn = append(nn, v)
		}
	}
	return nn
}

// RemoveSlice delete the elements of slice in `removes`.
//
// RemoveSlice([]string{"a", "b", "c", "a"}, "a") => []string{"b", "c"})
//
// RemoveSlice([]string{"a", "b", "c", "a"}, "b", "c") => []string{"a", "a"})
func RemoveSlice(ss []string, removes ...string) []string {
	m := make(map[string]struct{})
	for _, rm := range removes {
		m[rm] = struct{}{}
	}
	result := make([]string, 0, len(ss))
	for _, s := range ss {
		if _, ok := m[s]; ok {
			continue
		}
		result = append(result, s)
	}
	return result
}

// Exist check if elem exist in slice.
func Exist(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// NormalizeNewlines normalizes \r\n (windows) and \r (mac)
// into \n (unix).
//
// There are 3 ways to represent a newline.
//   Unix: using single character LF, which is byte 10 (0x0a), represented as “” in Go string literal.
//   Windows: using 2 characters: CR LF, which is bytes 13 10 (0x0d, 0x0a), represented as “” in Go string literal.
//   Mac OS: using 1 character CR (byte 13 (0x0d)), represented as “” in Go string literal. This is the least popular.
func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

var fontKinds = [][]int{{10, 48}, {26, 97}, {26, 65}}

// RandStr get random string.
func RandStr(size int) string {
	result := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		ikind := rand.Intn(3)
		scope, base := fontKinds[ikind][0], fontKinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

// ReverseSlice reverse slice.
//
// ReverseSlice([]string{"s1", "s2", "s3"} => []string{"s3", "s2", "s1}
func ReverseSlice(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
