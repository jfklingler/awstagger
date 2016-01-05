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

func getInstanceIds(svc *ec2.EC2) []*string {
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

func getAmiIds(svc *ec2.EC2) []*string {
	imagesOut, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Owners: []*string{&self},
	})

	kingpin.FatalIfError(err, "Could not retrieve EC2 instances")

	var imageIds []*string
	for _, image := range imagesOut.Images {
		imageIds = append(imageIds, image.ImageId)
	}

	return imageIds
}

func getVolumeIds(svc *ec2.EC2) []*string {
	volumesOut, err := svc.DescribeVolumes(&ec2.DescribeVolumesInput{})

	kingpin.FatalIfError(err, "Could not retrieve EC2 volumes")

	var volumeIds []*string
	for _, volume := range volumesOut.Volumes {
		volumeIds = append(volumeIds, volume.VolumeId)
	}

	return volumeIds
}

func getSnapshotIds(svc *ec2.EC2) []*string {
	ownerId := "owner-id"
	snapshotsOut, err := svc.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
		Filters: []*ec2.Filter{
			{
				Name: &ownerId,
				Values: []*string{&self},
			},
		},
	})

	kingpin.FatalIfError(err, "Could not retrieve EC2 snapshots")

	var snapshotIds []*string
	for _, snapshot := range snapshotsOut.Snapshots {
		snapshotIds = append(snapshotIds, snapshot.SnapshotId)
	}

	return snapshotIds
}

func getVpcIds(svc *ec2.EC2) []*string {
	vpcsOut, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{})

	kingpin.FatalIfError(err, "Could not retrieve EC2 VPCs")

	var vpcIds []*string
	for _, vpc := range vpcsOut.Vpcs {
		vpcIds = append(vpcIds, vpc.VpcId)
	}

	return vpcIds
}

func getSecurityGroupIds(svc *ec2.EC2) []*string {
	ownerId := "owner-id"
	securityGroups, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name: &ownerId,
				Values: []*string{&self},
			},
		},
	})

	kingpin.FatalIfError(err, "Could not retrieve EC2 security groups")

	var sgIds []*string
	for _, sg := range securityGroups.SecurityGroups {
		sgIds = append(sgIds, sg.GroupId)
	}

	return sgIds
}
