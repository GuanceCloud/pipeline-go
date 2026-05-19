# Cloud Security Patterns

Use this reference for SIEM/CSPM rules based on cloud audit logs such as AWS CloudTrail.

## Review Checklist

- Classify the rule as SIEM when it detects runtime account/API activity, and CSPM when it checks resource posture.
- For CloudTrail console login rules, keep `source`, cloud provider, event name, account ID, principal, login status, source IP, and user agent as evidence when available.
- For S3 ACL exposure rules, keep bucket name, region, principal ARN, source IP, user agent, and the ACL/public-principal match reason as evidence when available.
- Alias nested DQL fields to simple names used by the script, for example `last(\`@userIdentity.accountId\`) as accountId`.
- If grouping by a nested field, project that same field with an alias before reading it through `dql_series_get()`.
- `last()` collapses multiple events to one record per group. Use it for "latest state/activity" rules; use `count(*) BY status/account` or explicit event rows when the rule needs every attempt.
- Treat `@responseElements.ConsoleLogin == "Success"` as a successful root login. Treat non-success values as failed/unknown attempts unless the rule explicitly ignores them.
- Check each status/severity branch with mock DQL results before final delivery.

## AWS Root ConsoleLogin Template

This template detects active AWS root account console login events from CloudTrail. Successful root login is high severity; failed or unknown attempts are medium by default.

```python
event_result = dql("L('default')::`cloudtrail`:(last(`@responseElements.ConsoleLogin`) as login_status, last(`evtName`) as evtName, last(`@userIdentity.userName`) as userName, last(`@userIdentity.accountId`) as accountId) {`source` = 'cloudtrail' and `evtName` = 'ConsoleLogin' and `@userIdentity.type` = 'Root' and `eventType` != 'AwsServiceEvent'} BY `@userIdentity.accountId`")

if len(event_result["series"]) > 0 {
    loginStatuses = dql_series_get(event_result, "login_status")
    accountIds = dql_series_get(event_result, "accountId")
    userNames = dql_series_get(event_result, "userName")

    for i=0; i<len(loginStatuses); i+=1 {
        group = loginStatuses[i]
        for j=0; j<len(group); j+=1 {
            loginStatus = group[j]
            if loginStatus != nil {
                accountId = ""
                userName = "Root"

                if len(accountIds) > i && len(accountIds[i]) > j && accountIds[i][j] != nil {
                    accountId = accountIds[i][j]
                }
                if len(userNames) > i && len(userNames[i]) > j && userNames[i][j] != nil {
                    userName = userNames[i][j]
                }

                eventTitle = "failed login"
                status = "medium"
                if loginStatus == "Success" {
                    eventTitle = "successful login"
                    status = "high"
                }

                trigger(
                    result=1,
                    status=status,
                    dimension_tags={
                        "source": "cloudtrail",
                        "iaas": "aws",
                        "security_type": "cloud"
                    },
                    related_data={
                        "accountId": accountId,
                        "userName": userName,
                        "loginStatus": loginStatus,
                        "event_title": eventTitle,
                        "eventName": "ConsoleLogin"
                    }
                )
            }
        }
    }
}
```

Validation samples:

```sh
./bin/arbiter-check \
  --script ./rule.p \
  --dql-result '{"series":[[{"columns":{"login_status":"Success","evtName":"ConsoleLogin","userName":"root","accountId":"123456789012"},"tags":{"@userIdentity.accountId":"123456789012"}}]],"status_code":200}' \
  --check-dql \
  --require-trigger \
  --expect-status high

./bin/arbiter-check \
  --script ./rule.p \
  --dql-result '{"series":[[{"columns":{"login_status":"Failure","evtName":"ConsoleLogin","userName":"root","accountId":"123456789012"},"tags":{"@userIdentity.accountId":"123456789012"}}]],"status_code":200}' \
  --check-dql \
  --require-trigger \
  --expect-status medium
```

## AWS S3 Public Bucket ACL Template

This template detects CloudTrail `PutBucketAcl` events that grant public access to an S3 bucket through `AllUsers`, `AuthenticatedUsers`, `public-read`, or `public-read-write`.

```python
event_result = dql("L('default')::`cloudtrail`:(last(`evtName`) as eventName, last(`@userIdentity.arn`) as userArn, last(`region`) as region, last(`@requestParameters.bucketName`) as bucketName, last(`@sourceIPAddress`) as sourceIPAddress, last(`@userAgent`) as userAgent) {`source` = 'cloudtrail' and `evtName` = 'PutBucketAcl' and (`message` = queryString('AllUsers') or `message` = queryString('AuthenticatedUsers') or `message` = queryString('public-read') or `message` = queryString('public-read-write'))} BY `@requestParameters.bucketName`")

if len(event_result["series"]) > 0 {
    bucketNames = dql_series_get(event_result, "bucketName")
    userArns = dql_series_get(event_result, "userArn")
    regions = dql_series_get(event_result, "region")
    eventNames = dql_series_get(event_result, "eventName")
    sourceIPs = dql_series_get(event_result, "sourceIPAddress")
    userAgents = dql_series_get(event_result, "userAgent")

    for i=0; i<len(bucketNames); i+=1 {
        group = bucketNames[i]
        for j=0; j<len(group); j+=1 {
            bucketName = group[j]
            if bucketName != nil && bucketName != "" {
                region = ""
                userArn = ""
                eventName = "PutBucketAcl"
                sourceIPAddress = ""
                userAgent = ""

                if len(regions) > i && len(regions[i]) > j && regions[i][j] != nil {
                    region = regions[i][j]
                }
                if len(userArns) > i && len(userArns[i]) > j && userArns[i][j] != nil {
                    userArn = userArns[i][j]
                }
                if len(eventNames) > i && len(eventNames[i]) > j && eventNames[i][j] != nil {
                    eventName = eventNames[i][j]
                }
                if len(sourceIPs) > i && len(sourceIPs[i]) > j && sourceIPs[i][j] != nil {
                    sourceIPAddress = sourceIPs[i][j]
                }
                if len(userAgents) > i && len(userAgents[i]) > j && userAgents[i][j] != nil {
                    userAgent = userAgents[i][j]
                }

                trigger(
                    result=1,
                    status="high",
                    dimension_tags={
                        "bucketName": bucketName,
                        "source": "cloudtrail",
                        "region": region,
                        "iaas": "aws",
                        "security_type": "cloud"
                    },
                    related_data={
                        "eventName": eventName,
                        "userArn": userArn,
                        "sourceIPAddress": sourceIPAddress,
                        "userAgent": userAgent,
                        "bucketName": bucketName,
                        "reason": "AWS S3 bucket ACL grants public access"
                    }
                )
            }
        }
    }
}
```

Validation sample:

```sh
./bin/arbiter-check \
  --script ./rule.p \
  --dql-result '{"series":[[{"columns":{"eventName":"PutBucketAcl","userArn":"arn:aws:iam::123456789012:user/alice","region":"us-east-1","bucketName":"public-bucket","sourceIPAddress":"203.0.113.10","userAgent":"aws-cli/2"},"tags":{"@requestParameters.bucketName":"public-bucket"}}]],"status_code":200}' \
  --check-dql \
  --require-trigger \
  --expect-status high
```
