package comm

import (
	"fmt"
)

type ErrorGroupT struct {
	ellors    []error
	dumpStack bool
}

type ErrorGroup = *ErrorGroupT

func NewErrorGroup(dumpStack bool) ErrorGroup {
	return &ErrorGroupT{
		ellors:    make([]error, 0, 8),
		dumpStack: dumpStack,
	}
}

func (me ErrorGroup) HasError() bool {
	return me.AmountOfErrors() > 0
}

func (me ErrorGroup) AmountOfErrors() int {
	return len(me.ellors)
}

func (me ErrorGroup) Error() string {
	msgs := make([]string, 0, len(me.ellors)*2)
	for i, err := range me.ellors {
		if i > 0 && me.dumpStack {
			msgs = append(msgs, "")
		}
		msg := fmt.Sprintf("error #%d - %s", i+1, err.Error())
		msgs = append(msgs, msg)
	}

	return fmt.Sprintf("%d errors totally:\n%s", len(me.ellors), JoinedLines(msgs...))
}

func (me ErrorGroup) Add(err error) {
	if err != nil {
		me.ellors = append(me.ellors, err)
	}
}

func (me ErrorGroup) AddAll(that ErrorGroup) {
	if that != nil && len(that.ellors) > 0 {
		me.ellors = append(me.ellors, that.ellors...)
	}
}
