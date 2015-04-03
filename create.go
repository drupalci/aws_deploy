package main

import (
	"fmt"

	"github.com/mitchellh/goamz/ec2"
	"github.com/mitchellh/goamz/elb"
	"github.com/mitchellh/multistep"
)

type StepCreate struct{}

func (s *StepCreate) Run(state multistep.StateBag) multistep.StepAction {
	clientEc2 := state.Get("client_ec2").(*ec2.EC2)
	clientElb := state.Get("client_elb").(*elb.ELB)

	ami := state.Get("ami").(string)
	key := state.Get("key").(string)
	size := state.Get("size").(string)
	amount := state.Get("amount").(int)

	// Spin up the instances.
	options := ec2.RunInstances{
		ImageId:      ami,
		KeyName:      key,
		InstanceType: size,
		MinCount:     amount,
		MaxCount:     amount,
	}
	resp, err := clientEc2.RunInstances(&options)
	Check(err)

	// Assign these to the correct ELB instance.
	for _, instance := range resp.Instances {
		fmt.Println("Creating: ", instance.InstanceId)
		add := &elb.RegisterInstancesWithLoadBalancer{
			LoadBalancerName: *elbId,
			Instances:        []string{instance.InstanceId},
		}
		_, err = clientElb.RegisterInstancesWithLoadBalancer(add)
		Check(err)
	}

	return multistep.ActionContinue
}

func (s *StepCreate) Cleanup(multistep.StateBag) {
	// This is called after all the steps have run or if the runner is
	// cancelled so that cleanup can be performed.
}
