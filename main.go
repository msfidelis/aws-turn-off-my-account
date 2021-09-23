package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Response events.APIGatewayProxyResponse

var (
	ec2Svc *ec2.EC2
	elbSvc *elbv2.ELBV2
	rdsSvc *rds.RDS
	stsSvc *sts.STS
	elcSvc *elasticache.ElastiCache
)

func Handler(ctx context.Context) error {
	//Handle
	ec2Handle()
	albHandle()
	rdsHandle()
	elasticacheHandle()

	return nil
}

func main() {

	region := os.Getenv("REGION")

	//Prepare
	sess, err := getAWSSession(region)

	if err != nil {
		panic(err)
	}

	// Services
	ec2Svc = ec2.New(sess)
	elbSvc = elbv2.New(sess)
	rdsSvc = rds.New(sess)
	stsSvc = sts.New(sess)
	elcSvc = elasticache.New(sess)

	// Handle Lambda
	lambda.Start(Handler)
	// elasticacheHandle()
}

func ec2Handle() {
	fmt.Println("Searching for EC2 Instances")

	instances, err := getEc2Instances()

	if err != nil {
		panic(err)
	}

	terminateInstances(instances)

	fmt.Println("Searching for Snapshots")

	snapshots, err := getEc2Snapshots()

	if err != nil {
		fmt.Println(err)
	}

	terminateSnapShots(snapshots)

	fmt.Println("Searching for EBS Volumes")

	volumes, err := getEBSVolumes()

	if err != nil {
		panic(err)
	}

	terminateEBS(volumes)
}

func albHandle() {
	fmt.Println("Searching for ALB / ELB / NLBs Instances")

	instances, err := getLoadBalancersInstances()

	if err != nil {
		panic(err)
	}

	terminateLoadBalancers(instances)
}

func rdsHandle() {
	fmt.Println("Searching for RDS Instances")

	instances, err := getRDSInstances()

	if err != nil {
		panic(err)
	}

	clusters, err := getRDSClusters()

	if err != nil {
		panic(err)
	}

	terminateRDSInstances(instances)
	terminateRDSClusters(clusters)
}

func elasticacheHandle() {
	fmt.Println("Searching for Elasticache Clusters")

	clusters, err := getElasticacheClusters()

	if err != nil {
		panic(err)
	}

	terminateElasticacheClusters(clusters)

	fmt.Println("Searching for Elasticache Repliocation Groups")

	groups, err := getReplicationGroups()

	if err != nil {
		panic(err)
	}

	terminateReplicationGroups(groups)
}


func getElasticacheClusters() ([]*string, error) {

	var instances []*string

	input := &elasticache.DescribeCacheClustersInput{}
	result, err := elcSvc.DescribeCacheClusters(input)

	if err != nil {
		return nil, err
	}

	for _, cluster := range result.CacheClusters {
		if (*cluster.CacheClusterStatus == "available") {
			instances = append(instances, cluster.CacheClusterId)
		}
	}

	return instances, nil

}

func getReplicationGroups() ([]*string, error) {

	var instances []*string

	input := &elasticache.DescribeReplicationGroupsInput{}
	result, err := elcSvc.DescribeReplicationGroups(input)

	if err != nil {
		return nil, err
	}

	for _, rpg := range result.ReplicationGroups {
		instances = append(instances, rpg.ReplicationGroupId)
	}

	return instances, nil 
}


func getLoadBalancersInstances() ([]*string, error) {

	var instances []*string

	input := &elbv2.DescribeLoadBalancersInput{}
	result, err := elbSvc.DescribeLoadBalancers(input)

	if err != nil {
		return nil, err
	}

	for _, lb := range result.LoadBalancers {
		instances = append(instances, lb.LoadBalancerArn)
	}

	return instances, nil

}

func getEc2Instances() ([]*string, error) {

	var instances []*string

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("pending")},
			},
		},
	}

	result, err := ec2Svc.DescribeInstances(input)

	if err != nil {
		return nil, err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			instances = append(instances, instance.InstanceId)
		}
	}

	return instances, nil
}

func getEc2Snapshots() ([]*string, error) {

	var snapshots []*string

	account, err := getAWSAccount()

	if err != nil {
		return nil, err
	}

	filters := []*ec2.Filter{
		{
			Name: aws.String("owner-id"),
			Values: []*string{
				aws.String(account),
			},
		},
	}

	input := &ec2.DescribeSnapshotsInput{Filters: filters}

	result, err := ec2Svc.DescribeSnapshots(input)

	if err != nil {
		return nil, err
	}

	for _, snapshot := range result.Snapshots {
		snapshots = append(snapshots, snapshot.SnapshotId)
	}

	return snapshots, nil

}

func getEBSVolumes() ([]*string, error) {
	var volumes []*string

	result, err := ec2Svc.DescribeVolumes(nil)

	if err != nil {
		return nil, err
	}

	for _, volume := range result.Volumes {
		volumes = append(volumes, volume.VolumeId)
	}

	return volumes, nil
}

func getRDSInstances() ([]*string, error) {
	var instances []*string

	result, err := rdsSvc.DescribeDBInstances(nil)

	if err != nil {
		return nil, err
	}

	for _, rds := range result.DBInstances {
		instances = append(instances, rds.DBInstanceIdentifier)
	}

	return instances, nil
}

func getRDSClusters() ([]*string, error) {
	var instances []*string

	result, err := rdsSvc.DescribeDBClusters(nil)

	if err != nil {
		return nil, err
	}

	for _, rds := range result.DBClusters {
		instances = append(instances, rds.DBClusterIdentifier)
	}

	return instances, nil
}

func terminateInstances(instances []*string) {

	if len(instances) == 0 {
		fmt.Println("No more EC2 instances to destroy")
	} else {

		params := &ec2.TerminateInstancesInput{
			InstanceIds: instances,
		}

		resp, err := ec2Svc.TerminateInstances(params)

		if err != nil {
			fmt.Printf("Failed to terminate instance", err)
		}

		for _, ti := range resp.TerminatingInstances {
			fmt.Printf("Instance: %s \n\nStatus: %s", *ti.InstanceId, ti.CurrentState.String())
		}

	}

}

func terminateSnapShots(snapshots []*string) {
	for _, snapshot := range snapshots {
		input := &ec2.DeleteSnapshotInput{
			SnapshotId: snapshot,
		}

		_, err := ec2Svc.DeleteSnapshot(input)

		if err != nil {
			fmt.Printf("Failed to terminate snapshot", err)
		}
	}
}

func terminateEBS(volumes []*string) {
	for _, volume := range volumes {
		fmt.Println(*volume)

		params := &ec2.DeleteVolumeInput{
			VolumeId: volume,
		}

		_, err := ec2Svc.DeleteVolume(params)

		if err != nil {
			fmt.Printf("Failed to terminate EBS", err)
		}
	}
}

func terminateLoadBalancers(instances []*string) {
	for _, instance := range instances {

		params := &elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: instance,
		}

		_, err := elbSvc.DeleteLoadBalancer(params)

		if err != nil {
			fmt.Printf("Failed to terminate lb", err)
		}

	}
}

func terminateRDSInstances(instances []*string) {
	for _, instance := range instances {

		params := &rds.DeleteDBInstanceInput{
			DBInstanceIdentifier: instance,
			SkipFinalSnapshot:    aws.Bool(true),
		}

		_, err := rdsSvc.DeleteDBInstance(params)

		if err != nil {
			fmt.Printf("Failed to terminate RDS", err)
		}

	}
}

func terminateRDSClusters(instances []*string) {
	for _, instance := range instances {

		params := &rds.DeleteDBClusterInput{
			DBClusterIdentifier: instance,
			SkipFinalSnapshot:   aws.Bool(true),
		}

		_, err := rdsSvc.DeleteDBCluster(params)

		if err != nil {
			fmt.Printf("Failed to terminate RDS Cluster", err)
		}

	}
}

func terminateElasticacheClusters(clusters []*string) {
	for _, cluster := range clusters {
		params := &elasticache.DeleteCacheClusterInput{
			CacheClusterId: cluster,
		}

		_, err := elcSvc.DeleteCacheCluster(params)

		if err != nil {
			fmt.Printf("Failed to terminate Replication Group", err)
		}

	}
}

func terminateReplicationGroups(groups []*string) {
	for _, group := range groups {
		params := &elasticache.DeleteReplicationGroupInput{
			ReplicationGroupId: group,
		}

		_, err := elcSvc.DeleteReplicationGroup(params)

		if err != nil {
			fmt.Printf("Failed to terminate Cache Cluster", err)
		}

	}
}

func getAWSSession(region string) (*session.Session, error) {
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}

	awsConfig = awsConfig.WithCredentialsChainVerboseErrors(true)
	return session.NewSession(awsConfig)
}

func getAWSAccount() (string, error) {
	callerInput := &sts.GetCallerIdentityInput{}
	output, err := stsSvc.GetCallerIdentity(callerInput)

	if err != nil {
		return "", err
	}

	return *output.Account, nil
}
