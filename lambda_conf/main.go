package lambda_conf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func UpdateEnv(conf map[string]string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := lambda.New(sess, &aws.Config{Region: aws.String("ap-northeast-1")})

	// Get current environment variables
	getInput := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String("test2"),
	}

	getResult, err := svc.GetFunctionConfiguration(getInput)
	if err != nil {
		return err
	}

	// Update ACCESS_TOKEN
	getResult.Environment.Variables["ACCESS_TOKEN"] = aws.String(conf["AccessToken"])

	// Update environment variables
	updateInput := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String("test2"),
		Environment:  &lambda.Environment{Variables: getResult.Environment.Variables},
	}

	_, err = svc.UpdateFunctionConfiguration(updateInput)
	if err != nil {
		return err
	}

	return nil
}
