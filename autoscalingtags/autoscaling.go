/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
*/

package autoscalingtags

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/jfklingler/awstagger/context"
)

func Process(ctx context.Context, region string) {
	if !ctx.TagFlags.AutoScalingGroups {
		return
	}

	svc := autoscaling.New(ctx.AwsSession, &aws.Config{Region: aws.String(region)})


	processGroups := func (out *autoscaling.DescribeAutoScalingGroupsOutput, lastPage bool) bool {
		var names []*string

		for _, group := range out.AutoScalingGroups {
			names = append(names, group.AutoScalingGroupName)
		}

		updateTags(svc, makeNewTags(ctx.Tags, ctx.TagFlags.AsgPropogate, names))
		deleteTags(svc, makeDeleteTags(ctx.RmTags, names))
		printTags(ctx, svc)

		return !lastPage
	}

	ctx.Print("  Processing auto-scaling groups...")
	svc.DescribeAutoScalingGroupsPages(&autoscaling.DescribeAutoScalingGroupsInput{}, processGroups)
}

func updateTags(svc *autoscaling.AutoScaling, tags []*autoscaling.Tag) {
	_, err := svc.CreateOrUpdateTags(&autoscaling.CreateOrUpdateTagsInput{
		Tags: tags,
	})

	kingpin.FatalIfError(err, "Could not update tags for auto-scaling groups: %s", tags)
}

func deleteTags(svc *autoscaling.AutoScaling, tags []*autoscaling.Tag) {
	_, err := svc.DeleteTags(&autoscaling.DeleteTagsInput{
		Tags: tags,
	})

	kingpin.FatalIfError(err, "Could not delete tags for auto-scaling groups: %s", tags)
}

func printTags(ctx context.Context, svc *autoscaling.AutoScaling) {
	if !ctx.Verbose {
		return
	}
	resp, err := svc.DescribeTags(&autoscaling.DescribeTagsInput{})

	kingpin.FatalIfError(err, "Could not retrieve tags for auto-scaling groups")

	lastID := ""
	for _, td := range resp.Tags {
		if lastID != *td.ResourceId {
			ctx.PrintVerbose(fmt.Sprintf("    Resource ID %s", *td.ResourceId))
			lastID = *td.ResourceId
		}
		ctx.PrintVerbose(fmt.Sprintf("      %s=%s (propogate: %t)", *td.Key, *td.Value, *td.PropagateAtLaunch))
	}
}

func makeNewTags(tags map[string]string, prop bool, asgs []*string) []*autoscaling.Tag {
	astags := make([]*autoscaling.Tag, 0)

	for _, asg := range asgs {
		for k, v := range tags {
			vx := v
			astags = append(astags, makeTag(k, &vx, &prop, asg))
		}
	}

	return astags
}

func makeDeleteTags(tags []string, asgs []*string) []*autoscaling.Tag {
	astags := make([]*autoscaling.Tag, 0)

	for _, asg := range asgs {
		for _, k := range tags {
			astags = append(astags, makeTag(k, nil, nil, asg))
		}
	}

	return astags
}

func makeTag(key string, value *string, prop *bool, asg *string) *autoscaling.Tag {
	asgt := "auto-scaling-group"
	return &autoscaling.Tag{
		ResourceId: asg,
		Key: &key,
		Value: value,
		PropagateAtLaunch: prop,
		ResourceType: &asgt,
	}
}
