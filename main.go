package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v1"
	"github.com/awslabs/aws-sdk-go/aws"
    //"github.com/awslabs/aws-sdk-go/service/ec2"
    "github.com/awslabs/aws-sdk-go/service/elb"
)

var (
	elbId  = kingpin.Flag("elb", "Identifier for the Elastic Load Balancer.").Required().String()
	amiId  = kingpin.Flag("ami", "Identifier for the Image to be deployed.").Required().String()
	count  = kingpin.Flag("num", "The amount of instances to be delpoyed.").Default("1").String()
	tags   = kingpin.Flag("tag", "Tag the environment's for billing purposes.").Required().String()
	region = kingpin.Flag("region", "Deploy the images to a Region.").Required().String()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.Help = "AWS deployment tools."
	kingpin.Parse()

    clientElb := elb.New(&aws.Config{Region: *region})

    // Do we have an ELB with this ID.
    query := &elb.DescribeLoadBalancersInput{
    	LoadBalancerNames: []*string {elbId},
    }

    elbs, err := clientElb.DescribeLoadBalancers(query)
    if err != nil {
        panic(err)
    }

    // resp has all of the response data, pull out instance IDs:
    fmt.Println("> Number of reservation sets: ", len(elbs.LoadBalancerDescriptions))
    for idx, res := range elbs.LoadBalancerDescriptions {
        fmt.Println("  > Number of instances: ", len(res.Instances))
        for _, inst := range elbs.LoadBalancerDescriptions[idx].Instances {
            fmt.Println("    - Instance ID: ", *inst.InstanceID)
        }
    }
}
