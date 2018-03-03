package checker

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jaym/proto-compat/pkg/checker/rules"
	"github.com/pkg/errors"
)

type checker struct {
	oldProto *descriptor.FileDescriptorSet
	newProto *descriptor.FileDescriptorSet
}

func New(oldProto, newProto *descriptor.FileDescriptorSet) *checker {
	return &checker{
		oldProto: oldProto,
		newProto: newProto,
	}
}

func (c *checker) Check(reporter rules.Reporter) error {
	for ruleDescriptor, ruleFunc := range rules.Rules {
		err := ruleFunc(reporter, c.oldProto, c.newProto)
		if err != nil {
			errors.Wrapf(err, "Failed to check rule %s", ruleDescriptor.Name)
		}
	}
	return nil
}
