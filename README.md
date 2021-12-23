# arnlike

This package exports ArnLike, a function implementing the ArnLike comparison as per [the AWS IAM documentation](
https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_ARN)

```go
arnlike.ArnLike("arn:aws:iam::000000000000:role/some-role", "arn:aws:iam::000000000000:role/some-*") // true, nil
```
