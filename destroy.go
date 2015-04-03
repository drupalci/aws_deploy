package main

import (
    "fmt"

    "github.com/mitchellh/goamz/ec2"
    "github.com/mitchellh/goamz/elb"
    "github.com/mitchellh/multistep"
)

type StepDestroy struct{}

func (s *StepDestroy) Run(state multistep.StateBag) multistep.StepAction {
    clientEc2 := state.Get("client_ec2").(*ec2.EC2)
    clientElb := state.Get("client_elb").(*elb.ELB)
    elbId := state.Get("elb").(string)

    // This is the query that we will perform to find our load balancer.
    query := &elb.DescribeLoadBalancer{
        Names: []string{elbId},
    }

    // Get all the load balancers from AWS's API.
    bals, err := clientElb.DescribeLoadBalancers(query)
    Check(err)

    bal := bals.LoadBalancers[0]
    for _, instance := range bal.Instances {
        fmt.Println("Removing: ", instance.InstanceId)
        
        // Luckily we don't need to worry about deregistering from the Balancer first.
        _, err := clientEc2.TerminateInstances([]string{instance.InstanceId})
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
