package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/goamz/ec2"
	"github.com/mitchellh/goamz/elb"
	"github.com/mitchellh/multistep"
)

type StepCreate struct{}

func (s *StepCreate) Run(state multistep.StateBag) multistep.StepAction {
	clientEc2 := state.Get("client_ec2").(*ec2.EC2)
	clientElb := state.Get("client_elb").(*elb.ELB)

	ami := state.Get("ami").(string)
	size := state.Get("size").(string)
	amount := state.Get("amount").(int)

	// Spin up the instances.
	options := ec2.RunInstances{
		ImageId:      ami,
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

		// Tag the instances so we know what they are.
		tags := buildTags(state.Get("tags").(string))
		clientEc2.CreateTags([]string{instance.InstanceId}, tags)
	}

	return multistep.ActionContinue
}

func (s *StepCreate) Cleanup(multistep.StateBag) {
	// This is called after all the steps have run or if the runner is
	// cancelled so that cleanup can be performed.
}

func buildTags(t string) []ec2.Tag {
	var tags []ec2.Tag

    tSlice := strings.Split(t, ",")
    for _, tag := range tSlice {
        tagSplit := strings.Split(tag, "=")
        newTag := &ec2.Tag{
        	Key:   tagSplit[0],
			Value: tagSplit[0],
        }
        tags = append(tags, *newTag)
    }

    return tags
}
