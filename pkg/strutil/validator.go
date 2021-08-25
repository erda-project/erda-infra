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

package strutil

import (
	"fmt"
	"regexp"
	"unicode"
)

// Validator defines a validator function.
// User can extend validator in their own packages.
type Validator func(s string) error

// Validate validate `s` with composed validators and return error if have
func Validate(s string, validators ...Validator) error {
	for _, v := range validators {
		if err := v(s); err != nil {
			return err
		}
	}
	return nil
}

// MinLenValidator verify if `s` meets the minimum length requirement.
func MinLenValidator(minLen int) Validator {
	return func(s string) error {
		if len(s) < minLen {
			if minLen == 1 {
				return fmt.Errorf("cannot be empty")
			}
			return fmt.Errorf("less than min length: %d", minLen)
		}
		return nil
	}
}

// MaxLenValidator check whether `s` exceeds the maximum length.
func MaxLenValidator(maxLen int) Validator {
	return func(s string) error {
		if len(s) > maxLen {
			return fmt.Errorf("over max length: %d", maxLen)
		}
		return nil
	}
}

// MaxRuneCountValidator check max rune count.
func MaxRuneCountValidator(maxLen int) Validator {
	return func(s string) error {
		if len([]rune(s)) > maxLen {
			return fmt.Errorf("over max rune count: %d", maxLen)
		}
		return nil
	}
}

var envKeyRegexp = regexp.MustCompilePOSIX(`^[a-zA-Z_]+[a-zA-Z0-9_]*$`)

// EnvKeyValidator check whether `s` meets the linux env key specification.
var EnvKeyValidator Validator = func(s string) error {
	valid := envKeyRegexp.MatchString(s)
	if !valid {
		return fmt.Errorf("illegal env key, validated by regexp: %s", envKeyRegexp.String())
	}
	return nil
}

// EnvValueLenValidator check whether `s` exceeds the maximum length of linux env value.
var EnvValueLenValidator = MaxLenValidator(128 * 1024)

// NoChineseValidator check whether `s` contains Chinese characters.
var NoChineseValidator Validator = func(s string) error {
	var chineseCharacters []string
	for _, runeValue := range s {
		if unicode.Is(unicode.Han, runeValue) {
			chineseCharacters = append(chineseCharacters, string(runeValue))
		}
	}
	if len(chineseCharacters) > 0 {
		return fmt.Errorf("found %d chinese characters: %s", len(chineseCharacters),
			Join(chineseCharacters, " ", true))
	}
	return nil
}

// AlphaNumericDashUnderscoreValidator regular expression verification, can only:
// - start with uppercase and lowercase letters or numbers
// - supports uppercase and lowercase letters, numbers, underscores, underscores, and dots
var AlphaNumericDashUnderscoreValidator Validator = func(s string) error {
	exp := `^([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$`
	valid := regexp.MustCompile(exp).MatchString(s)
	if !valid {
		return fmt.Errorf("valid regexp: %s", exp)
	}
	return nil
}
