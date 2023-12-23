package lambda_conf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Config struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
}

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
	getResult.Environment.Variables["REFRESH_TOKEN"] = aws.String(conf["RefreshToken"])

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

func GetEnv() (Config, error) {
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
		return Config{}, err
	}

	return Config{
		AccessToken:  *getResult.Environment.Variables["ACCESS_TOKEN"],
		RefreshToken: *getResult.Environment.Variables["REFRESH_TOKEN"],
		ClientID:     *getResult.Environment.Variables["CLIENT_ID"],
	}, nil
}
