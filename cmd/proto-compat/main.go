package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jaym/proto-compat/pkg/checker"
	"github.com/jaym/proto-compat/pkg/checker/rules"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage:\n\t%s oldproto.pb newproto.pb\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	var descriptorSetOld descriptor.FileDescriptorSet
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(data, &descriptorSetOld)

	var descriptorSetNew descriptor.FileDescriptorSet
	data, err = ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}
	err = proto.Unmarshal(data, &descriptorSetNew)
	reporter := rules.NewCLIReporter()

	c := checker.New(&descriptorSetOld, &descriptorSetNew)
	c.Check(reporter)
	os.Exit(reporter.Summurize())
}
