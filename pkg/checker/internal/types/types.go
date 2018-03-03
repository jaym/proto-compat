package types

import (
	"fmt"
)

type FullyQualifiedProtoTypeName string

func ProtoTypeName(base FullyQualifiedProtoTypeName, name string) FullyQualifiedProtoTypeName {
	if string(base) == "" {
		return FullyQualifiedProtoTypeName(name)
	}
	return FullyQualifiedProtoTypeName(fmt.Sprintf("%s.%s", string(base), name))
}
