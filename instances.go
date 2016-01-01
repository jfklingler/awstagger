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

func tagInstances(ctx context, session *session.Session, region string) {
	ctx.print("    Processing EC2 instances...")

	svc := ec2.New(session, &aws.Config{Region: aws.String(region)})

	instancesOut := getInstances(*svc)

	for idx := range instancesOut.Reservations {
		for _, instance := range instancesOut.Reservations[idx].Instances {
			ctx.printVerbose(fmt.Sprintf("    Processing instance %s...", *instance.InstanceId))

			tagsOut := getTags(*svc, *instance.InstanceId)

			for _, td := range tagsOut.Tags {
				ctx.printVerbose(fmt.Sprintf("      %s=%s", *td.Key, *td.Value))
			}
		}
	}
}

func getInstances(svc ec2.EC2) ec2.DescribeInstancesOutput {
	instancesOut, err := svc.DescribeInstances(nil)

	kingpin.FatalIfError(err, "Could not retrieve EC2 instances")

	return *instancesOut
}

func getTags(svc ec2.EC2, instanceID string) ec2.DescribeTagsOutput {
	resp, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{ // Required
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	})

	kingpin.FatalIfError(err, "Could not retrieve tags for EC2 instance %s", instanceID)

	return *resp
}
