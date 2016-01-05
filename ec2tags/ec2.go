/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
*/

package ec2tags

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jfklingler/awstagger/context"
)

func Process(ctx context.Context, region string) {
	svc := ec2.New(ctx.AwsSession, &aws.Config{Region: aws.String(region)})

	switch {

	case ctx.TagFlags.Ec2Instances:
		ctx.Print("  Processing EC2 instances...")
		applyTags(ctx, svc, getInstanceIds(svc))
		fallthrough

	case ctx.TagFlags.Ec2Amis:
		ctx.Print("  Processing EC2 AMIs...")
		applyTags(ctx, svc, getAmiIds(svc))
		fallthrough

	case ctx.TagFlags.Ec2Volumes:
		ctx.Print("  Processing EC2 volumes...")
		applyTags(ctx, svc, getVolumeIds(svc))
		fallthrough

	case ctx.TagFlags.Ec2Snapshots:
		ctx.Print("  Processing EC2 snapshots...")
		applyTags(ctx, svc, getSnapshotIds(svc))
		fallthrough

	case ctx.TagFlags.Ec2Vpcs:
		ctx.Print("  Processing EC2 VPCs...")
		applyTags(ctx, svc, getVpcIds(svc))
		fallthrough

	case ctx.TagFlags.Ec2SecurityGroups:
		ctx.Print("  Processing EC2 security groups...")
		applyTags(ctx, svc, getSecurityGroupIds(svc))
		fallthrough

	case ctx.TagFlags.Ec2NetInterfaces:
		ctx.Print("  Processing EC2 network interfaces...")
		applyTags(ctx, svc, getNetIfaceIds(svc))
	}
}

func applyTags(ctx context.Context, svc *ec2.EC2, resourceIds []*string) {
	updateTags(ctx, *svc, resourceIds)
	deleteTags(ctx, *svc, resourceIds)
	printTags(ctx, *svc, resourceIds)
}

func updateTags(ctx context.Context, svc ec2.EC2, instanceIds []*string) {
	if len(ctx.Tags) <= 0 {
		return
	}

	resp, err := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: instanceIds,
		Tags:      tagArgsToEc2Tags(ctx.Tags),
	})

	kingpin.FatalIfError(err, "Could not update tags for EC2 instances %s", instanceIds)

	fmt.Println(resp)
}

func deleteTags(ctx context.Context, svc ec2.EC2, instanceIds []*string) {
	if len(ctx.RmTags) <= 0 {
		return
	}

	resp, err := svc.DeleteTags(&ec2.DeleteTagsInput{
		Resources: instanceIds,
		Tags:      rmtagArgsToEc2Tags(ctx.RmTags),
	})

	kingpin.FatalIfError(err, "Could not delete tags for EC2 instances %s", instanceIds)

	fmt.Println(resp)
}

func getTags(svc ec2.EC2, instanceIds []*string) ec2.DescribeTagsOutput {
	resp, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{ // Required
				Name:   aws.String("resource-id"),
				Values: instanceIds,
			},
		},
	})

	kingpin.FatalIfError(err, "Could not retrieve tags for EC2 instances %s", instanceIds)

	return *resp
}

func printTags(ctx context.Context, svc ec2.EC2, instanceIds []*string) {
	if ctx.Verbose {
		tagsOut := getTags(svc, instanceIds)
		lastID := ""

		for _, td := range tagsOut.Tags {
			if lastID != *td.ResourceId {
				ctx.PrintVerbose(fmt.Sprintf("    Instance %s", *td.ResourceId))
				lastID = *td.ResourceId
			}
			ctx.PrintVerbose(fmt.Sprintf("      %s=%s", *td.Key, *td.Value))
		}
	}
}

func rmtagArgsToEc2Tags(tags []string) []*ec2.Tag {
	var ec2Tags []*ec2.Tag

	for _, k := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key: &k,
		})
	}

	return ec2Tags
}

func tagArgsToEc2Tags(tags map[string]string) []*ec2.Tag {
	var ec2Tags []*ec2.Tag

	for k, v := range tags {
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   &k,
			Value: &v,
		})
	}

	return ec2Tags
}
