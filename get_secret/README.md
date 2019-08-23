# get_secret

A helper for fetching secrets from AWS secrets manager.

- See how to [use get_secret on CircleCI](#use-get_secret-on-circle-ci)
- Find cross-platform binaries in the Releases page
- Alternative approach: AWS CLI [get-secret-value](https://docs.aws.amazon.com/cli/latest/reference/secretsmanager/get-secret-value.html#examples)

## Basic Usage

Assuming a secret is stored as:

```
Secret Name: mysystem/prod

Secret:
{
    "url": "https://mysystem.example.com",
    "key": "My_Awsome_key"
}
```

This tool requires that AWS Credentials have been configured.
eg:

- via ~/.aws/credentials
- or `export AWS_ACCESS_KEY_ID=xyz ACCESS_SECRET_ACCESS_KEY=abc`

This tool will fetch the values based on the keys:

```
get_secret -region us-west-2 mysystem/prod/url

# returns
https://mysystem.example.com
```

```
get_secret -region us-west-2 mysystem/prod/key

# returns
My_Awsome_key
```

## Get All Secrets Values

```
get_secret --all -region us-west-2 mysystem/prod

# returns all the values
key=My_Awsome_key
url= https://mysystem.example.com
```

Get all Secret values in a single line (no newline spacing)
```
AWS_REGION=us-west-2
get_secret --all -n mysystem/prod

# returns all the values
key=My_Awsome_key url=https://mysystem.example.com
```

That's useful if you want to export the secrets
```
export $(get_secret --all -n mysystem/uat)
```



## Syntax

Fetching single values for a key:

- get_secret -region AWS_REGION SECRET_NAME/SECRET_JSON_KEY

Fetching all values for a secret:

- get_secret --all -region AWS_REGION SECRET_NAME

### Tips

- Region can also be specified by exporting the environment variable `AWS_REGION`

```
export AWS_REGION=us-east-1
get_secret mysecret/a_key
```

- You can print the secret without newline spacing (trim newline) by passing `-n` switch
  - Useful for parsing in shell scripts

```
export AWS_REGION=us-east-1
get_secret mysecret/a_key
```

- Only supports fetching of values that are either numbers, strings or boolean

--

That's all

---

### Use get_secret on Circle-CI

`get_secret` is useful in sourcing credentials for an application on Circle CI. The example below illustrates how to use it to populate ElasticBeanstalk environment variables

- Ensure the AWS credentials are configured
- Ensure you've exported a Github Personal Access Token that can fetch the binary

```
export AWS_PROFILE=us-west-2

# download get_secret binary
curl -vLJ -H 'Accept: application/octet-stream' \
-o get_secret \
-H "Authorization: token ${GITHUB_PAT}" \
'https://api.github.com/repos/peteretelej/tools/releases/assets/13954806'

chmod +x get_secret

# use to fetch secrets
DB_HOST=$(./get_secret mysystem/proddb/host)
DB_USER=$(./get_secret mysystem/proddb/user)
DB_PASS=$(./get_secret mysystem/proddb/pass)
DB_NAME=$(./get_secret mysystem/proddb/name)

eb setenv DB_HOST=$DB_HOST DB_USER=$DB_USER DB_PASS=$DB_PASS DB_NAME=$DB_NAME -e beanstalk-env

# OR Use all secrets

eb setenv $(./get_secret --all -n mysystem/proddb) -e beanstalk-env

eb deploy beanstalk-env
```

- If you want to update the release binary asset number (ie binary version), you can check for the latest binary release ID by running:

```
export GITHUB_PAT=_GITHUB_PERSONAL_ACCESS_TOKEN_

curl -H "Authorization: token ${GITHUB_PAT}" \
'https://api.github.com/repos/peteretelej/tools/releases/latest'
```
