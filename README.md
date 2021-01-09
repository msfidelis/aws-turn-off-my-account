# aws-turn-off-my-account
Lambda stack to turn off all account resources to avoid billing surprises

## Instalation

### Using Serverless Framework 

* Clone this repo 

```bash
git clone ghttps://github.com/msfidelis/aws-turn-off-my-account.git
```

* Edit your preferences in `configs/prod.yml` and customize your cron rate on `serverless.yml`

* Deploy 

```bash
make deploy 
```

