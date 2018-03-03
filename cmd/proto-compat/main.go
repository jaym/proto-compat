package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jaym/proto-compat/pkg/checker"
	"github.com/jaym/proto-compat/pkg/checker/rules"
)

func main() {
	var descriptorSetOld descriptor.FileDescriptorSet
	data, err := ioutil.ReadFile("/tmp/proto/foo")
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(data, &descriptorSetOld)

	var descriptorSetNew descriptor.FileDescriptorSet
	data, err = ioutil.ReadFile("/tmp/proto/bar")
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(data, &descriptorSetNew)
	reporter := rules.NewCLIReporter()

	c := checker.New(&descriptorSetOld, &descriptorSetNew)
	c.Check(reporter)
	os.Exit(reporter.Summurize())
}
