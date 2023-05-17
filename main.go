package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws-apigateway/sdk/go/apigateway"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		policy := map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Action": []string{
						"sts:AssumeRole",
					},
					"Principal": map[string]interface{}{
						"Service": "lambda.amazonaws.com",
					},
					"Effect": "Allow",
					"Sid":    "",
				},
			},
		}
		policyStr, err := json.Marshal(policy)

		if err != nil {
			return err
		}

		echoHandlerRole, err := iam.NewRole(ctx, "echo-handler-role", &iam.RoleArgs{
			Description:      pulumi.String("Role used by the lambda"),
			AssumeRolePolicy: pulumi.String(string(policyStr)),
		})

		if err != nil {
			return err
		}

		echoHandler, err := lambda.NewFunction(ctx, "echo-handler-func", &lambda.FunctionArgs{
			Name:        pulumi.String("echo-handler"),
			Description: pulumi.String("Providing a message will echo the message"),
			Code:        pulumi.NewFileArchive("./echo"),
			Runtime:     pulumi.String("go1.x"),
			Role:        echoHandlerRole.Arn,
			Handler:     pulumi.String("echo"),
		})

		if err != nil {
			return err
		}

		getMethod := apigateway.MethodGET
		restAPI, err := apigateway.NewRestAPI(ctx, "pulumi-serverless-gateway", &apigateway.RestAPIArgs{
			Routes: []apigateway.RouteArgs{
				{
					Path:         "/echo/{thing}",
					Method:       &getMethod,
					EventHandler: echoHandler,
				},
			},
		})

		if err != nil {
			return err
		}

		ctx.Export("url", restAPI.Url)

		return nil
	})
}
