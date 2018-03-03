package rules

import "github.com/golang/protobuf/protoc-gen-go/descriptor"

type Descriptor struct {
	Name       string
	Descripton string
}

type RuleRunFunc func(reporter Reporter, oldProto *descriptor.FileDescriptorSet, newProto *descriptor.FileDescriptorSet) error

var Rules = map[Descriptor]RuleRunFunc{
	reserveRemoveFieldsRuleDescriptor: CheckRemovedFieldsAreReserved,
}
