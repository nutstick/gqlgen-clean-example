package logging

import (
	"regexp"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type testPrinter struct {
	fxtest.TB
}

// NewTestPrinter returns a fx.Printer that logs to the testing TB.
func NewTestPrinter(t fxtest.TB) fx.Printer {
	return &testPrinter{t}
}

func (p *testPrinter) Printf(format string, args ...interface{}) {
	if regexp.MustCompile(`PROVIDE|INVOKE|RUNNING`).MatchString(format) {
	} else {
		p.Logf(format, args...)
	}
}
