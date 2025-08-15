## AWS

## Deploying docker images to ECR

- Create `deploy.yml`.
- Go to `github.com/marketplace`
- Select Actions
- Search for AWS ECR
- Select AWS ECR "Login" action
- Copy the AWS ECR tempalte for private registry

```yml
steps:
  - name: Login to Amazon ECR
    id: login-ecr
    uses: aws-actions/amazon-ecr-login@v2

  - name: Build, tag, and push docker image to Amazon ECR
    id: build-image
    env:
      REGISTRY: ${{ steps.login-ecr.outputs.registry }}
      REPOSITORY: my-ecr-repo
      IMAGE_TAG: ${{ github.sha }}
    run: |
      docker build -t $REGISTRY/$REPOSITORY:$IMAGE_TAG .
      docker push $REGISTRY/$REPOSITORY:$IMAGE_TAG
```


### AWS ECR

Click on `Create Repository`. Select the visibility (private/public), select the `mutability` for tags. Finally the `encryption`; which is **AES** by default.

### AWS IAM

Search for IAM. Its a web service that lets use securely access the AWS resources.

- `IAM users` can be one person or an application.
- `IAM user groups` can set permissions for multiple users at once (groups for admin, developers, tester etc.).
- `IAM role` is similar to `users` but it can be assumed by anyone who needs it.

> For now we need `IAM Users`. Create a user with the name `github-ci`.

But to create a user we first need to create a group.  

Create a group with the name deployment and the select the `AmazonEC2ContainerRegistryPowerUser` policy. This will give all the access except for deletion.  

> Make sure to add the same permission for the `github-ci` user as well.

### Github Secrets

Go to **Settings > Secrets** and select **Repository Secrets**.  

Enter the AWS_SECRET, AWS_REGION values etc.

## AWS RDS

Search for RDS and select `Choose a database creation method` followed by engine options i.e. `Postgres`.  

Select the version and the templates (Production, Dev/Test, Free Tier).  

Under `settings` provide the DB instance identifier and the user/password.  

> Since we are using Free Tier there is only type of `DB instance class` available i.e. `db.t3.micro` or `db.t4.micro`.

Default storage is 20GiB of SSD storage. 

### Connectivity

For now select a default VPC setting. DB subnet group as `default`.  

Create new `VPC security group (firewall)` so we can allow access to specified ports.  

The VPC which was created will only allow request from 1 inbound rule with your current IP address.

So change the source to anywhere i.e. `0.0.0.0/0`.

Finally, apply the DB schema using

> migrate -path db/migration -database "postgres://root:<AWS_DB_PASSWORD>@simple-bank.<AWS_DB_ARN>/simple_bank" -verbose up

## Store/retrieve production secrets

We must not keep the secrets/password/urls in the `app.env`. We can use the `AWS Secret Manager` service.

> **AWS Secret Manager > Secrets > Store new secret**

### Choose secret type

Use `key/value pairs` secret type.

A god PASETO token can be generated using `openssl` on your local machine

```sh
$ openssl rand -hex 64 | head -c 32
2e9afe82b8b042604aec19d47afbd548
```

### Storing IAM secrets

Configure AWS and store the secrets on your system using the `AWS CLI`.

```sh
aws configure
AWS Access Key ID [None]: xxx
AWS Secret Access Key [None]: xxx
Default region name [None]: xxx
Default output format [None]: json
```

These configs are present under `~/.aws`

```
ls ~/.aws
config		credentials
```

The `credentials` file has the secrets. Whereas, `config` file holds the region or output.

**NOTE**: `default` is name of the AWS profile; we can set multiple AWS profiles.

### Pulling the AWS secrets

We can use the **secret key name** or the **ARN** to pull the secret.

```sh
$ aws secretsmanager get-secret-value --secret-id simple_bank

# Need to login to the AWS ECR
An error occurred (AccessDeniedException) when calling the GetSecretValue operation: User: arn:aws:iam::[ACCOUNT_ID]:user/github-ci is not authorized to perform: secretsmanager:GetSecretValue on resource: simple_bank because no identity-based policy allows the secretsmanager:GetSecretValue action
```

Either we create a new user or give extra permissions to the existing (_github-ci_) user; we can do so by adding the `SecretsManagerReadWrite` permission to the existing user.

Now, `get-secret-value` will work.

```sh
$ aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text

{
  "ARN": "xxx",
  "Name": "simple_bank",
  "VersionId": "cd4a78fa-6552-49a6-a0ae-73e696d0f91f",
  "SecretString": "{}",
  "VersionStages": [
    "AWSCURRENT"
  ],
  "CreatedDate": "2025-08-15T07:19:47.971000+05:30"
}
```

`--query SecretString`: gets the secret payload
`--output text`: gets it in JSON format

### Accessing the secrets in our app

We need to populate the `app.env` from the CLI

```sh
$ aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]'

"DB_SOURCE=postgres://root:[PASSWORD]@[DB_HOST]:5432/simple_bank"
"DB_DRIVER=postgres"
"SERVER_ADDRESS=0.0.0.0:8080"
"TOKEN_SYMMETRIC_KEY=[YOUR_TOKEN_KEY]"
"ACCESS_TOKEN_DURATION=15m"
```

### Pulling images from ECR

We cannot pull the AWS ECR images without login

```sh
docker pull [ACCOUNT_ID].dkr.ecr.[REGION].amazonaws.com/simplebank:[IMAGE_TAG]
Error response from daemon: failed to resolve reference "044116941566.dkr.ecr.eu-north-1.amazonaws.com/simplebank:2bbd12a93928c263bed58d6bef617c6ba0834463": pull access denied, repository does not exist or may require authorization: authorization failed: no basic auth credentials
```

Use the ecr API

```
aws ecr get-login-password | docker login --username AWS --password-stdin [ACCOUNT_ID].dkr.ecr.[REGION].amazonaws.com

Login Succeeded
```

### Creating container

If we try to run the docker container, the startup will fail because of `DB_SOURCE` being empty.

To fix this issue we need to update `start.sh` and `source` the ENV variables before the migration starts

```sh
source /app/app.env
```

This is required because running the local images, the Go application can read the ENV directly from the OS environment w/o needing to load the life.

In the CI/CD the secrets are written to an `app.env` file inside the container during the build process.
The container doesn't have these an OS ENV variables.
