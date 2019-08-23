package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretKeySeparator is the separator for secret name and key
const SecretKeySeparator = "/"

var awsRegion = "us-west-2" // default aws region

func main() {

	var (
		region = flag.String("region", "", "AWS Region to use")
		n      = flag.Bool("n", false, "No newlines when printing secret")
		all    = flag.Bool("all", false, "Returns a list of all values for a given secret name")
	)
	flag.Parse()

	log.SetPrefix("[get_secret] ")
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatal("invalid arguments, kindly pass the secret & key as a single parameter")
	}

	parseRegion(*region)

	lastArg := getLastArg(os.Args)

	var name string
	if *all {
		name = lastArg

		ss, err := fetchAllSecrets(name)
		if err != nil {
			log.Fatalf("failed to fetch all secrets for [name = %s] [AWS REGION = %s ] \n\t%v",
				name, awsRegion, err,
			)
		}

		if *n {
			fmt.Print(strings.Join(ss, " "))
			return
		}

		fmt.Println(strings.Join(ss, "\n"))
		return
	}

	name, key, err := parseNameKey(lastArg)
	if err != nil {
		log.Fatal("invalid arguments, secret name/key must be passed last")
	}
	v, err := fetchSecretValue(name, key)
	if err != nil {
		log.Fatalf(
			"failed to fetch secret for [name = %s ] [ key = %s ] [ AWS REGION = %s] \n\t%v",
			name, key, awsRegion, err,
		)
	}
	if *n {
		fmt.Print(v)
		return
	}
	fmt.Println(v)
}

// name/key => name,key  nam/e/key => nam/e,key
func parseNameKey(name string) (string, string, error) {
	if !strings.Contains(name, SecretKeySeparator) {
		return "", "", errors.New("secret name must contain a '/' to separate name/key")
	}

	i := strings.LastIndex(name, SecretKeySeparator)

	n := name[:i]
	k := name[i+1:]

	if n == "" || k == "" {
		return "", "", errors.New("missing secret name / secret key")
	}

	return n, k, nil
}

func fetchSecretValue(name, key string) (string, error) {
	values, err := getSecretValues(name)
	if err != nil {
		return "", fmt.Errorf("failed to get_secret [%s] : %v", name, err)
	}
	vInt, ok := values[key]
	if !ok {
		return "", fmt.Errorf("missing requested key in secret")
	}

	switch v := vInt.(type) {
	case int, string, bool:
		return fmt.Sprintf("%s", v), nil
	}

	return "", fmt.Errorf("secret value is an unsupported data type (should be number/string/bool)")
}

func fetchAllSecrets(name string) ([]string, error) {
	if name == "" {
		return nil, errors.New("missing secret name")
	}
	values, err := getSecretValues(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get_secret [%s] : %v", name, err)
	}
	var pairs []string
	for k, v := range values {
		k := fmt.Sprintf("%s", k)
		s := fmt.Sprintf("%s", v)

		if strings.Contains(k, " ") || strings.Contains(s, " ") {
			log.Printf("warning for secret [%s]: skipping secret key '%s' (or its value) contains spaces", name, k)
			continue
		}
		pairs = append(pairs,
			fmt.Sprintf(`%s=%s`, k, s),
		)
	}
	return pairs, nil
}

func parseRegion(region string) {
	if region != "" {
		awsRegion = region
		return
	}

	if r := os.Getenv("AWS_REGION"); r != "" {
		awsRegion = r
	}
}

func getSecretValues(name string) (map[string]interface{}, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	secret, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, err
	}

	values := make(map[string]interface{})

	err = json.Unmarshal([]byte(*secret.SecretString), &values)
	if err != nil {
		return nil, fmt.Errorf("unable to parse secret values as json: %v", err)
	}

	return values, nil
}

func getLastArg(args []string) string {
	last := args[len(args)-1:][0]

	for _, v := range []string{"us-", "-"} {
		if strings.HasPrefix(last, v) {
			return ""
		}
	}
	return last
}
