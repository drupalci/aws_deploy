package main

import (
	"fmt"
  "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/mitchellh/multistep"
)

type StepDestroy struct{}

func (s *StepDestroy) Run(state multistep.StateBag) multistep.StepAction {
	clientEc2 := state.Get("client_ec2").(*ec2.EC2)
	clientElb := state.Get("client_elb").(*elb.ELB)
	elbName := state.Get("elb").(string)

  query := &elb.DescribeLoadBalancersInput{
    LoadBalancerNames: []*string{
      aws.String(elbName),
    },
  }

	// Get all the load balancers from AWS's API.
	bals, err := clientElb.DescribeLoadBalancers(query)
	Check(err)

  bal := bals.LoadBalancerDescriptions[0]
	for _, instance := range bal.Instances {
		fmt.Println("Removing: ", instance.InstanceID)

		// Luckily we don't need to worry about deregistering from the Balancer first.
    params := &ec2.TerminateInstancesInput{
      InstanceIDs: []*string{
        instance.InstanceID,
      },
    }
		_, err := clientEc2.TerminateInstances(params)
		Check(err)
	}

	// This is so the other commands know how many instances were destroyed.
	state.Put("amount", len(bal.Instances))

	return multistep.ActionContinue
}

func (s *StepDestroy) Cleanup(multistep.StateBag) {
	// This is called after all the steps have run or if the runner is
	// cancelled so that cleanup can be performed.
}
