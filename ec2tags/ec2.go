/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
*/

package ec2tags

import (
	"fmt"
	"math"

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
		processInstances(svc, &ctx.BatchSize, applyTags(ctx, svc))
		fallthrough

	case ctx.TagFlags.Ec2Amis:
		ctx.Print("  Processing EC2 AMIs...")
		processAmis(svc, ctx.BatchSize, applyTags(ctx, svc))
		fallthrough

	case ctx.TagFlags.Ec2Volumes:
		ctx.Print("  Processing EC2 volumes...")
		processVolumes(svc, ctx.BatchSize, applyTags(ctx, svc))
		fallthrough

	case ctx.TagFlags.Ec2Snapshots:
		ctx.Print("  Processing EC2 snapshots...")
		processSnapshots(svc, ctx.BatchSize, applyTags(ctx, svc))
		fallthrough

	case ctx.TagFlags.Ec2Vpcs:
		ctx.Print("  Processing EC2 VPCs...")
		processVpcs(svc, ctx.BatchSize, applyTags(ctx, svc))
		fallthrough

	case ctx.TagFlags.Ec2SecurityGroups:
		ctx.Print("  Processing EC2 security groups...")
		processSecurityGroups(svc, ctx.BatchSize, applyTags(ctx, svc))
		fallthrough

	case ctx.TagFlags.Ec2NetInterfaces:
		ctx.Print("  Processing EC2 network interfaces...")
		processNetInterfaces(svc, ctx.BatchSize, applyTags(ctx, svc))
	}
}

func applyTags(ctx context.Context, svc *ec2.EC2) func([]*string) {
	return func (resourceIds []*string) {
		nextWindow := func (ids []*string) ([]*string, []*string) {
			w := int(math.Min(float64(len(ids)), float64(200)))
			return ids[0:w], ids[w:]
		}

		for thisRound, remaining := nextWindow(resourceIds); len(thisRound) > 0; thisRound, remaining = nextWindow(remaining){
			updateTags(ctx, *svc, thisRound)
			deleteTags(ctx, *svc, thisRound)
			printTags(ctx, *svc, thisRound)
		}
	}
}

func updateTags(ctx context.Context, svc ec2.EC2, resourceIds []*string) {
	if len(ctx.Tags) <= 0 {
		return
	}

	_, err := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: resourceIds,
		Tags:      tagArgsToEc2Tags(ctx.Tags),
	})

	kingpin.FatalIfError(err, "Could not update tags for EC2 resources %s", resourceIds)
}

func deleteTags(ctx context.Context, svc ec2.EC2, resourceIds []*string) {
	if len(ctx.RmTags) <= 0 {
		return
	}

	_, err := svc.DeleteTags(&ec2.DeleteTagsInput{
		Resources: resourceIds,
		Tags:      rmtagArgsToEc2Tags(ctx.RmTags),
	})

	kingpin.FatalIfError(err, "Could not delete tags for EC2 resources %s", resourceIds)
}

func getTags(svc ec2.EC2, resourceIds []*string) ec2.DescribeTagsOutput {
	resp, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{ // Required
				Name:   aws.String("resource-id"),
				Values: resourceIds,
			},
		},
	})

	kingpin.FatalIfError(err, "Could not retrieve tags for EC2 resources %s", resourceIds)

	return *resp
}

func printTags(ctx context.Context, svc ec2.EC2, instanceIds []*string) {
	if ctx.Verbose {
		tagsOut := getTags(svc, instanceIds)
		lastID := ""

		for _, td := range tagsOut.Tags {
			if lastID != *td.ResourceId {
				ctx.PrintVerbose(fmt.Sprintf("    Resource ID %s", *td.ResourceId))
				lastID = *td.ResourceId
			}
			ctx.PrintVerbose(fmt.Sprintf("      %s=%s", *td.Key, *td.Value))
		}
	}
}

func rmtagArgsToEc2Tags(tags []string) []*ec2.Tag {
	var ec2Tags []*ec2.Tag

	for _, k := range tags {
		kx := k
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key: &kx,
		})
	}

	return ec2Tags
}

func tagArgsToEc2Tags(tags map[string]string) []*ec2.Tag {
	var ec2Tags []*ec2.Tag

	for k, v := range tags {
		kx, vx := k, v
		ec2Tags = append(ec2Tags, &ec2.Tag{
			Key:   &kx,
			Value: &vx,
		})
	}

	return ec2Tags
}
