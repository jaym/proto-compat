package rules

import (
	"fmt"
)

type Reporter interface {
	Fail(rule Descriptor, reason string)
}

type reportedFailure struct {
	rule   Descriptor
	reason string
}
type CLIReporter struct {
	failures []reportedFailure
}

func NewCLIReporter() *CLIReporter {
	return &CLIReporter{
		failures: []reportedFailure{},
	}
}

func (c *CLIReporter) Fail(rule Descriptor, reason string) {
	c.failures = append(c.failures, reportedFailure{rule, reason})
	fmt.Printf("[%s]: %s\n", rule.Name, reason)

}

func (c *CLIReporter) Summurize() int {
	fmt.Println()
	fmt.Printf("%d checks failed\n", len(c.failures))
	if len(c.failures) > 0 {
		return 1
	}
	return 0
}
