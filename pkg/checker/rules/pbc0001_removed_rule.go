package rules

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jaym/proto-compat/pkg/checker/internal/types"
	"github.com/jaym/proto-compat/pkg/checker/internal/visitor"
)

var reserveRemoveFieldsRuleDescriptor = Descriptor{
	Name:       "PBC0001",
	Descripton: "Removed fields must be reserved to prevent future reuse",
}

type reserveRemovedFieldsRule struct {
	oldProto *descriptor.FileDescriptorSet
	newProto *descriptor.FileDescriptorSet
}

type fullyQualifiedFieldPath string

type fieldInfo struct {
	name     string
	reserved bool
}

func CheckRemovedFieldsAreReserved(reporter Reporter, oldProto *descriptor.FileDescriptorSet, newProto *descriptor.FileDescriptorSet) error {
	rule := &reserveRemovedFieldsRule{oldProto, newProto}
	return rule.Run(reporter)
}

func (r *reserveRemovedFieldsRule) Run(reporter Reporter) error {
	oldProtoFields := gatherFieldInfo(r.oldProto)
	newProtoFields := gatherFieldInfo(r.newProto)

	for k, v := range oldProtoFields {
		_, exist := newProtoFields[k]
		if !exist {
			reporter.Fail(reserveRemoveFieldsRuleDescriptor, fmt.Sprintf("Field %s[%s] has been removed but not reserved", k, v.name))
		}
	}

	return nil
}

func gatherFieldInfo(d *descriptor.FileDescriptorSet) map[fullyQualifiedFieldPath]fieldInfo {
	fieldInfoMap := make(map[fullyQualifiedFieldPath]fieldInfo)

	visitor.AllMessages(d, func(visitPath types.FullyQualifiedProtoTypeName, msg *descriptor.DescriptorProto) error {
		for _, field := range msg.GetField() {
			fieldInfoMap[fieldPath(visitPath, field.GetNumber())] = fieldInfo{
				name: field.GetName(),
			}
		}

		// Mark reserved fields
		for _, reserved := range msg.GetReservedRange() {
			for i := reserved.GetStart(); i <= reserved.GetEnd(); i++ {
				fieldInfoMap[fieldPath(visitPath, i)] = fieldInfo{
					reserved: true,
				}
			}

		}
		return nil
	})

	return fieldInfoMap
}

func fieldPath(base types.FullyQualifiedProtoTypeName, fieldNumber int32) fullyQualifiedFieldPath {
	return fullyQualifiedFieldPath(fmt.Sprintf("%s/%d", base, fieldNumber))
}
