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

### Github Secrets

Go to Settings > Secrets and select `Repository Secrets`.
