# AWS Deploy

A toolset for managing the DrupalCI stack on AWS

## Installation

### Requirements

* Vagrant 1.6+

### Setup

Spin up the Vagrant host

```bash
$ vagrant up
```

## Building

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
