/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
*/

package context

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	quiet   = kingpin.Flag("quiet", "Minimal/no output.").Short('q').Bool()
	verbose = kingpin.Flag("verbose", "Verbose output.").Short('v').Bool()

	regions = kingpin.Flag("region", "AWS region to process. (repeatable, default: all standard regions)").Short('r').Strings()
	tags    = kingpin.Flag("tag", "Tag to set/update on all selected resources. (repeatable)").Short('t').PlaceHolder("KEY=VALUE").Strings()
	rmTags  = kingpin.Flag("rm-tag", "Tag key to remove from all selected resources. (repeatable)").PlaceHolder("KEY").Strings()

	doEc2Instances = kingpin.Flag("ec2-instance", "Tag EC2 instances. (default: true)").Default("true").Bool()
	doEc2Amis      = kingpin.Flag("ec2-ami", "Tag EC2 AMIs. (default: true)").Default("true").Bool()
	doEc2Volumes   = kingpin.Flag("ec2-volume", "Tag EC2 volumes. (default: true)").Default("true").Bool()
	doEc2Snapshots = kingpin.Flag("ec2-snapshot", "Tag EC2 snapshots. (default: true)").Default("true").Bool()
	doEc2Vpcs      = kingpin.Flag("ec2-vpc", "Tag EC2 VPCs. (default: true)").Default("true").Bool()
	doEc2SGs       = kingpin.Flag("ec2-security-groups", "Tag EC2 security groups. (default: true)").Default("true").Bool()
	doEc2NIs       = kingpin.Flag("ec2-network-interfaces", "Tag EC2 network interfaces. (default: true)").Default("true").Bool()

	allRegions = []string{"us-east-1", "us-west-1", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "sa-east-1"}
)

type tagFlags struct {
	Ec2Instances      bool
	Ec2Amis           bool
	Ec2Volumes        bool
	Ec2Snapshots      bool
	Ec2Vpcs           bool
	Ec2SecurityGroups bool
	Ec2NetInterfaces  bool
}

type Context struct {
	AwsSession *session.Session

	Quiet   bool
	Verbose bool

	Regions []string
	Tags    map[string]string
	RmTags  []string

	TagFlags tagFlags
}

func New() Context {
	tagMap := make(map[string]string)
	for _, tag := range *tags {
		kv := strings.SplitN(tag, "=", 2)
		tagMap[kv[0]] = kv[1]
	}

	if len(*regions) == 0 {
		regions = &allRegions
	}

	return Context{
		AwsSession: session.New(),

		Quiet:   *quiet,
		Verbose: *verbose,

		Regions: *regions,
		Tags:    tagMap,
		RmTags:  *rmTags,

		TagFlags: tagFlags{
			Ec2Instances:      *doEc2Instances,
			Ec2Amis:           *doEc2Amis,
			Ec2Volumes:        *doEc2Volumes,
			Ec2Snapshots:      *doEc2Snapshots,
			Ec2Vpcs:           *doEc2Vpcs,
			Ec2SecurityGroups: *doEc2SGs,
			Ec2NetInterfaces:  *doEc2NIs,
		},
	}
}

func (c Context) Print(line string) {
	if !c.Quiet {
		fmt.Println(line)
	}
}

func (c Context) PrintVerbose(line string) {
	if c.Verbose && !c.Quiet {
		fmt.Println(line)
	}
}
