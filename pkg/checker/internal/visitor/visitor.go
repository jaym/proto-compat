package visitor

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jaym/proto-compat/pkg/checker/internal/types"
)

type VisitCallback func(types.FullyQualifiedProtoTypeName, *descriptor.DescriptorProto) error

func AllMessages(root *descriptor.FileDescriptorSet, visitCallback VisitCallback) error {
	visited := make(map[types.FullyQualifiedProtoTypeName]bool)

	for _, fdProto := range root.GetFile() {
		for _, msg := range fdProto.GetMessageType() {
			err := allMessages(msg, types.FullyQualifiedProtoTypeName(fdProto.GetPackage()), visited, visitCallback)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func allMessages(root *descriptor.DescriptorProto, visitPath types.FullyQualifiedProtoTypeName, visited map[types.FullyQualifiedProtoTypeName]bool, visitCallback VisitCallback) error {
	newVisitPath := types.ProtoTypeName(visitPath, root.GetName())

	_, exist := visited[newVisitPath]
	if exist {
		return nil
	}

	visited[newVisitPath] = true

	for _, msg := range root.GetNestedType() {
		err := allMessages(msg, newVisitPath, visited, visitCallback)
		if err != nil {
			return err
		}
	}

	return visitCallback(newVisitPath, root)
}
