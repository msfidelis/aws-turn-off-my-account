# AWS, Turn off my Account, please 

Lambda stack to turn off and destroy all resources from your personal AWS Account to avoid billing surprises

![Billing Modaf&*@&#](/.github/img/card.jpeg)

## Resources Roadmap

* EC2 :white_check_mark:
* EBS and Snapshots :white_check_mark:
* ALB, ELB, NLB :white_check_mark: 
* RDS Instances and Clusters :white_check_mark: 
* Elasticache Clusters and Replication Groups :white_check_mark: 
* Elastic IP's
* DocumentDB
* NAT Gateways

## Installation

### Using Serverless Framework 

* Clone this repo 

```bash
cd $GOPATH/src
git clone https://github.com/msfidelis/aws-turn-off-my-account.git
cd aws-turn-off-my-account
```

* Edit your preferences in `configs/prod.yml` and customize your cron rate on `serverless.yml`

* Deploy 

```bash
make deploy 
```

### Using console 

> Guenta ae 


### TODO

* Release ZIP 

* Console setup 

* Cloudformation Setup

* Tests 

* IAM Closed permissions

* Logs
