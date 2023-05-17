# pulumi-serverless

Creating a Serverless REST API with Pulumi. I decided to test out pulumi's method for Lambda creation compared to CDK and put a simple lambda behind an api gateway.

## Deploy

- you will need to have go installed as well as pulumi cli setup
- `cd echo && make build`
- `pulumi up`
