/*
Copyright 2016 John Klingler
Licensed under the MIT License (MIT)
*/

package ec2tags

import (
	"github.com/aws/aws-sdk-go/service/ec2"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	self = "self"
)

func processInstances(svc *ec2.EC2, pageSize *int64, apply func([]*string)) {
	err := svc.DescribeInstancesPages(&ec2.DescribeInstancesInput{
		MaxResults: pageSize,
	},
		func(instancesOut *ec2.DescribeInstancesOutput, lastPage bool) bool {
			var instanceIds []*string
			for idx := range instancesOut.Reservations {
				for _, instance := range instancesOut.Reservations[idx].Instances {
					instanceIds = append(instanceIds, instance.InstanceId)
				}
			}

			apply(instanceIds)

			return !lastPage
		})

	kingpin.FatalIfError(err, "Could not retrieve EC2 instances")
}

func processAmis(svc *ec2.EC2, pageSize int64, apply func([]*string)) {
	imagesOut, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Owners: []*string{&self},
	})

	kingpin.FatalIfError(err, "Could not retrieve EC2 instances")

	var imageIds []*string
	for _, image := range imagesOut.Images {
		imageIds = append(imageIds, image.ImageId)
	}

	apply(imageIds)
}

func processVolumes(svc *ec2.EC2, pageSize int64, apply func([]*string)) {
	err := svc.DescribeVolumesPages(&ec2.DescribeVolumesInput{},
		func(volumesOut *ec2.DescribeVolumesOutput, lastPage bool) bool {
			var volumeIds []*string
			for _, volume := range volumesOut.Volumes {
				volumeIds = append(volumeIds, volume.VolumeId)
			}

			apply(volumeIds)

			return !lastPage
		})

	kingpin.FatalIfError(err, "Could not retrieve EC2 volumes")
}

func processSnapshots(svc *ec2.EC2, pageSize int64, apply func([]*string)) {
	ownerId := "owner-id"

	err := svc.DescribeSnapshotsPages(&ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name: &ownerId,
				Values: []*string{&self},
			},
		},
	}, func(snapshotsOut *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
		var snapshotIds []*string
		for _, snapshot := range snapshotsOut.Snapshots {
			snapshotIds = append(snapshotIds, snapshot.SnapshotId)
		}

		apply(snapshotIds)

		return !lastPage
	})

	kingpin.FatalIfError(err, "Could not retrieve EC2 snapshots")
}

func processVpcs(svc *ec2.EC2, pageSize int64, apply func([]*string)) {
	vpcsOut, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{})

	kingpin.FatalIfError(err, "Could not retrieve EC2 VPCs")

	var vpcIds []*string
	for _, vpc := range vpcsOut.Vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
	}

	apply(vpcIds)
}

func processSecurityGroups(svc *ec2.EC2, pageSize int64, apply func([]*string)) {
	securityGroups, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})

	kingpin.FatalIfError(err, "Could not retrieve EC2 security groups")

	var sgIds []*string
	for _, sg := range securityGroups.SecurityGroups {
		sgIds = append(sgIds, sg.GroupId)
	}

	apply(sgIds)
}

func processNetInterfaces(svc *ec2.EC2, pageSize int64, apply func([]*string)) {
	networkInterfaces, err := svc.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{})

	kingpin.FatalIfError(err, "Could not retrieve EC2 network interfaces")

	var niIds []*string
	for _, ni := range networkInterfaces.NetworkInterfaces {
		niIds = append(niIds, ni.NetworkInterfaceId)
	}

	apply(niIds)
}
