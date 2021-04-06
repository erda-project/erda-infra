// Copyright 2021 Terminus
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

package urlenc

import "net/url"

// URLValuesUnmarshaler .
type URLValuesUnmarshaler interface {
	UnmarshalURLValues(prefix string, vals url.Values) error
}

// URLValuesMarshaler is the interface implemented by types that
// can marshal themselves into valid url.Values.
type URLValuesMarshaler interface {
	MarshalURLValues(prefix string, out url.Values) error
}
