package handy

import "errors"

//ErrAssertion is an assertion failed error
var ErrAssertion = errors.New("assertion")

//Assert panic WhatWillPanic for ErrAssertion if value is false
func Assert(value bool, WhatWillPanic ...interface{}) {
	if !value {
		if len(WhatWillPanic) == 0 || WhatWillPanic[0] == nil {
			panic(ErrAssertion)
		}
		panic(WhatWillPanic[0])
	}
}
