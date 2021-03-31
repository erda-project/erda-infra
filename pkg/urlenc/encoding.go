// Author: recallsong
// Email: songruiguo@qq.com

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
