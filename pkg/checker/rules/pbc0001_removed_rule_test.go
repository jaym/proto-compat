package rules

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type mockReporter struct {
	t     *testing.T
	fatal func(Descriptor, string)
}

func NewMockReporter(t *testing.T) *mockReporter {
	return &mockReporter{
		t: t,
	}
}

func (m *mockReporter) Fail(rule Descriptor, reason string) {
	if m.fatal != nil {
		m.fatal(rule, reason)
		return
	}
	m.t.Fatal("Unexpected call")
}

func TestNoChangesPBC0001(t *testing.T) {
	t.Run("simple message", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "generic/nochanges/simple")
		reporter := NewMockReporter(t)

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)
	})

	t.Run("nested messages", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "generic/nochanges/nested")
		reporter := NewMockReporter(t)

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)
	})

	t.Run("multiple messages", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "generic/nochanges/multiple")
		reporter := NewMockReporter(t)

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)
	})
}

func TestSimpleMessagePBC0001(t *testing.T) {
	t.Run("reports error when removed field is not reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/simple/incorrect")
		reporter := NewMockReporter(t)

		called := false
		calledPtr := &called
		reporter.fatal = func(rule Descriptor, reason string) {
			*calledPtr = true
			assert.Equal(t, "PBC0001", rule.Name)
			assert.NotEqual(t, "", reason)
		}

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)

		assert.True(t, *calledPtr)
	})

	t.Run("does not report an error when a field is removed and marked as reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/simple/correct")
		reporter := NewMockReporter(t)

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)
	})
}

func TestNestedMessagePBC0001(t *testing.T) {
	t.Run("reports error when removed field is not reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/nested/incorrect")
		reporter := NewMockReporter(t)

		called := false
		calledPtr := &called
		reporter.fatal = func(rule Descriptor, reason string) {
			*calledPtr = true
			assert.Equal(t, "PBC0001", rule.Name)
			assert.NotEqual(t, "", reason)
		}

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)

		assert.True(t, *calledPtr)
	})

	t.Run("does not report an error when a field is removed and marked as reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/nested/correct")
		reporter := NewMockReporter(t)

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)
	})
}

func TestMultipleMessagePBC0001(t *testing.T) {
	t.Run("reports error when removed field is not reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/multiple/incorrect")
		reporter := NewMockReporter(t)

		called := false
		calledPtr := &called
		reporter.fatal = func(rule Descriptor, reason string) {
			*calledPtr = true
			assert.Equal(t, "PBC0001", rule.Name)
			assert.NotEqual(t, "", reason)
		}

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)

		assert.True(t, *calledPtr)
	})

	t.Run("does not report an error when a field is removed and marked as reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/multiple/correct")
		reporter := NewMockReporter(t)

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)
	})
}

func TestOneofMessagePBC0001(t *testing.T) {
	t.Run("reports error when removed field is not reserved", func(t *testing.T) {
		oldProtos, newProtos := loadDescriptors(t, "pbc0001/oneof/incorrect")
		reporter := NewMockReporter(t)

		called := false
		calledPtr := &called
		reporter.fatal = func(rule Descriptor, reason string) {
			*calledPtr = true
			assert.Equal(t, "PBC0001", rule.Name)
			assert.NotEqual(t, "", reason)
		}

		CheckRemovedFieldsAreReserved(reporter, oldProtos, newProtos)

		assert.True(t, *calledPtr)
	})
}

func loadDescriptors(t *testing.T, name string) (*descriptor.FileDescriptorSet, *descriptor.FileDescriptorSet) {
	t.Helper()
	basePath := path.Join("testdata", name)
	oldPath := path.Join(basePath, "old.pb")
	newPath := path.Join(basePath, "new.pb")

	var descriptorSetOld descriptor.FileDescriptorSet
	data, err := ioutil.ReadFile(oldPath)
	if err != nil {
		t.Fatal("Could not load test data")
	}
	err = proto.Unmarshal(data, &descriptorSetOld)
	if err != nil {
		t.Fatal("Could not load test data")
	}

	var descriptorSetNew descriptor.FileDescriptorSet
	data, err = ioutil.ReadFile(newPath)
	if err != nil {
		t.Fatal("Could not load test data")
	}
	err = proto.Unmarshal(data, &descriptorSetNew)
	if err != nil {
		t.Fatal("Could not load test data")
	}

	return &descriptorSetOld, &descriptorSetNew
}
