AWS Deploy
==========

A toolset for managing deployments of AWS (AMI image) instances behind a load balancer.

## Usage

```bash
# Setup credentials.
export AWS_ACCESS_KEY=myaccesskey
export AWS_SECRET_KEY=mysecretkey

# Run the command.
aws_deploy --elb=balancer \
           --ami=ami-123456 \
           --key=my_key \
           --region=ap-southeast-2 \
           --tags="Name=Results,Environment=Production"
```

## Roadmap

* Check new instances are "InService" prior to removing the old one's out of the cluster.
* Run a script as part of the deploy.
* Tell AWS Autoscaling about the new AMI to avoid regressions.

## Building

### Requirements

* Vagrant 1.6+

#### Setup

Spin up the Vagrant host

```bash
$ vagrant up
```

All of these commands should be run within the vm:

```bash
$ vagrant ssh
$ cd /opt/golang/src/github.com/drupalci/aws_deploy
$ make
```

This will create a binary:

```
bin/aws_deploy
```
