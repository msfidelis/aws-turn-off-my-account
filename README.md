# aws-turn-off-my-account
Lambda stack to turn off all account resources to avoid billing surprises

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
