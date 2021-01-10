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
)

type Response events.APIGatewayProxyResponse

var (
	ec2Svc *ec2.EC2
	elbSvc *elbv2.ELBV2
	rdsSvc *rds.RDS
)

func Handler(ctx context.Context) error {
	//Handle
	ec2Handle()
	albHandle()
	rdsHandle()

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

	// Handle Lambda
	lambda.Start(Handler)
}

func ec2Handle() {
	fmt.Println("Searching for EC2 Instances")

	instances, err := getEc2Instances()

	if err != nil {
		panic(err)
	}

	terminateInstances(instances)
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

func terminateLoadBalancers(instances []*string) {

	for _, instance := range instances {

		params := &elbv2.DeleteLoadBalancerInput{
			LoadBalancerArn: instance,
		}

		resp, err := elbSvc.DeleteLoadBalancer(params)

		if err != nil {
			fmt.Printf("Failed to terminate lb", err)
		}

		fmt.Println(resp)

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

func getAWSSession(region string) (*session.Session, error) {
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}

	awsConfig = awsConfig.WithCredentialsChainVerboseErrors(true)
	return session.NewSession(awsConfig)
}
