package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/service/elb"
	"github.com/mitchellh/multistep"
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	elbId    = kingpin.Flag("elb", "Identifier for the Elastic Load Balancer.").Required().String()
	amiId    = kingpin.Flag("ami", "Identifier for the Image to be deployed.").Required().String()
	key      = kingpin.Flag("key", "The ssh key used to access the environments.").Required().String()
	size     = kingpin.Flag("size", "The size of the instance.").Default("t2.medium").String()
	tags     = kingpin.Flag("tags", "Tag the environment's for billing purposes.").Required().String()
	region   = kingpin.Flag("region", "Deploy the images to a Region.").Default("us-east-1").String()
	security = kingpin.Flag("security", "The security group.").Default("default").String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "AWS deployment tools."
	kingpin.Parse()

	state := new(multistep.BasicStateBag)

	// This allows us to share our client connections while in each of the steps.
  state.Put("client_elb", elb.New(&aws.Config{Region: *region}))
  state.Put("client_ec2", ec2.New(&aws.Config{Region: *region}))

	// Standard configuration that has been passed in via the CLI.
	state.Put("elb", *elbId)
	state.Put("ami", *amiId)
	state.Put("key", *key)
	state.Put("size", *size)
	state.Put("region", *region)
	state.Put("security", *security)
	state.Put("tags", *tags)

	steps := []multistep.Step{
		&StepDestroy{}, // Remove the existing hosts from the Load balancer.
		&StepCreate{},  // Create some EC2 instances and ensure they are ready to be deployed.
	}
	runner := &multistep.BasicRunner{Steps: steps}
	runner.Run(state)
}
