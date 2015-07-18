package main

import (
	"fmt"
	"strings"

  "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/mitchellh/multistep"
)

type StepCreate struct{}

func (s *StepCreate) Run(state multistep.StateBag) multistep.StepAction {
	clientEc2 := state.Get("client_ec2").(*ec2.EC2)
	clientElb := state.Get("client_elb").(*elb.ELB)

	ami := state.Get("ami").(string)
	key := state.Get("key").(string)
	security := state.Get("security").(string)
	size := state.Get("size").(string)
  elbName := state.Get("elb").(string)
  amount := state.Get("amount").(int)

	// Spin up the instances.
	options := &ec2.RunInstancesInput{
		ImageID:        aws.String(ami),
		KeyName:        aws.String(key),
		InstanceType:   aws.String(size),
		MinCount:       aws.Long(int64(amount)),
		MaxCount:       aws.Long(int64(amount)),
		SecurityGroups: []*string{ aws.String(security) },
	}
	resp, err := clientEc2.RunInstances(options)
	Check(err)

	// Assign these to the correct ELB instance.
	for _, instance := range resp.Instances {
		fmt.Println("Creating: ", instance.InstanceID)
    add := &elb.RegisterInstancesWithLoadBalancerInput{
      Instances: []*elb.Instance{
        { InstanceID: instance.InstanceID,
        },
      },
      LoadBalancerName: aws.String(elbName),
    }
    _, err := clientElb.RegisterInstancesWithLoadBalancer(add)
		Check(err)

		// Tag the instances so we know what they are.
		tags := buildTags(state.Get("tags").(string))
		clientEc2.CreateTags(&ec2.CreateTagsInput{
      Resources: []*string{instance.InstanceID},
      Tags:      tags,
    })
	}

	return multistep.ActionContinue
}

func (s *StepCreate) Cleanup(multistep.StateBag) {
	// This is called after all the steps have run or if the runner is
	// cancelled so that cleanup can be performed.
}

func buildTags(t string) []*ec2.Tag {
	var tags []*ec2.Tag

	tSlice := strings.Split(t, ",")
	for _, tag := range tSlice {
		tagSplit := strings.Split(tag, "=")
		newTag := &ec2.Tag{
			Key:   aws.String(tagSplit[0]),
			Value: aws.String(tagSplit[1]),
		}
		tags = append(tags, newTag)
	}

	return tags
}
