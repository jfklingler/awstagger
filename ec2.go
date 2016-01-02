/*
Copyright 2016 John Klingler
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"gopkg.in/alecthomas/kingpin.v2"
)

func TagInstances(ctx context, session *session.Session, region string) {
	ctx.print("  Processing EC2 instances...")

	svc := ec2.New(session, &aws.Config{Region: aws.String(region)})

	instanceIds := getInstanceIds(*svc)

	updateTags(ctx, *svc, instanceIds)
	deleteTags(ctx, *svc, instanceIds)
	printTags(ctx, *svc, instanceIds)
}

func TagAmis(ctx context, session *session.Session, region string) {
	ctx.print("  Processing AMIs...")

	svc := ec2.New(session, &aws.Config{Region: aws.String(region)})

	amiIds := getAmiIds(*svc)

	updateTags(ctx, *svc, amiIds)
	deleteTags(ctx, *svc, amiIds)
	printTags(ctx, *svc, amiIds)
}

func getInstanceIds(svc ec2.EC2) []*string {
	instancesOut, err := svc.DescribeInstances(nil)

	kingpin.FatalIfError(err, "Could not retrieve EC2 instances")

	var instanceIds []*string
	for idx := range instancesOut.Reservations {
		for _, instance := range instancesOut.Reservations[idx].Instances {
			instanceIds = append(instanceIds, instance.InstanceId)
		}
	}

	return instanceIds
}

func getAmiIds(svc ec2.EC2) []*string {
	// This seems idiotic. How can I simply create a []*string literal???
	var owners []*string
	self := "self"
	owners = append(owners, &self)

	imagesOut, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Owners: []*string{&self},
	})

	kingpin.FatalIfError(err, "Could not retrieve EC2 instances")

	var imageIds []*string
	for idx := range imagesOut.Images {
		imageIds = append(imageIds, imagesOut.Images[idx].ImageId)
	}

	return imageIds
}

func updateTags(ctx context, svc ec2.EC2, instanceIds []*string) {
	if len(ctx.tags) <= 0 {
		return
	}

	resp, err := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: instanceIds,
		Tags:      tagArgsToEc2Tags(ctx.tags),
	})

	kingpin.FatalIfError(err, "Could not update tags for EC2 instances %s", instanceIds)

	fmt.Println(resp)
}

func deleteTags(ctx context, svc ec2.EC2, instanceIds []*string) {
	if len(ctx.rmTags) <= 0 {
		return
	}

	resp, err := svc.DeleteTags(&ec2.DeleteTagsInput{
		Resources: instanceIds,
		Tags:      rmtagArgsToEc2Tags(ctx.rmTags),
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

func printTags(ctx context, svc ec2.EC2, instanceIds []*string) {
	if ctx.verbose {
		tagsOut := getTags(svc, instanceIds)
		lastID := ""

		for _, td := range tagsOut.Tags {
			if lastID != *td.ResourceId {
				ctx.printVerbose(fmt.Sprintf("    Instance %s", *td.ResourceId))
				lastID = *td.ResourceId
			}
			ctx.printVerbose(fmt.Sprintf("      %s=%s", *td.Key, *td.Value))
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
