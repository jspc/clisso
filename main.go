package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"bitbucket.org/emindsys/onelogin-aws-cli/onelogin"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/howeyc/gopass"
)

func main() {
	// TODO Use Cobra?

	// Get CLI arguments
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v <app_id>\n", os.Args[0])
		os.Exit(1)
	}
	var appId = os.Args[1]

	// Get env vars
	var secret string = os.Getenv("ONELOGIN_CLIENT_SECRET")
	var id string = os.Getenv("ONELOGIN_CLIENT_ID")
	var principal = os.Getenv("ONELOGIN_PRINCIPAL_ARN")
	var role = os.Getenv("ONELOGIN_ROLE_ARN")

	if secret == "" {
		log.Fatal("The ONELOGIN_CLIENT_SECRET environment variable must bet set.")
	}

	if id == "" {
		log.Fatal("The ONELOGIN_CLIENT_ID environment variable must bet set.")
	}

	// Get OneLogin access token
	log.Println("Generating OneLogin access tokens")
	token, err := onelogin.GenerateTokens(id, secret)
	if err != nil {
		log.Fatal(err)
	}

	// Get credentials from the user
	fmt.Print("OneLogin username: ")
	var user string
	fmt.Scanln(&user)

	fmt.Print("OneLogin password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		log.Fatal("Couldn't read password from terminal")
	}

	// Generate SAML assertion
	log.Println("Generating SAML assertion")
	pSaml := onelogin.GenerateSamlAssertionParams{
		UsernameOrEmail: user,
		Password:        string(pass),
		AppId:           appId,
		Subdomain:       "emind",
	}

	rSaml, err := onelogin.GenerateSamlAssertion(token, &pSaml)
	if err != nil {
		log.Fatal(err)
	}

	st := rSaml.Data[0].StateToken
	// TODO Handle multiple devices
	deviceId := strconv.Itoa(rSaml.Data[0].Devices[0].DeviceId)

	fmt.Print("Please enter your OneLogin OTP: ")
	var otp string
	fmt.Scanln(&otp)

	// Verify MFA
	pMfa := onelogin.VerifyFactorParams{
		AppId:      appId,
		DeviceId:   string(deviceId),
		StateToken: st,
		OtpToken:   otp,
	}

	rMfa, err := onelogin.VerifyFactor(token, &pMfa)
	if err != nil {
		log.Fatal(err)
	}

	samlAssertion := rMfa.Data

	// Assume role
	pAssumeRole := sts.AssumeRoleWithSAMLInput{
		PrincipalArn:  aws.String(principal),
		RoleArn:       aws.String(role),
		SAMLAssertion: aws.String(samlAssertion),
	}

	sess := session.Must(session.NewSession())
	svc := sts.New(sess)

	resp, err := svc.AssumeRoleWithSAML(&pAssumeRole)
	if err != nil {
		log.Fatal(err)
	}

	keyId := *resp.Credentials.AccessKeyId
	secretKey := *resp.Credentials.SecretAccessKey
	sessionToken := *resp.Credentials.SessionToken

	// Set temporary credentials in environment
	// TODO Error if already set
	fmt.Println("Paste the following in your terminal:")
	fmt.Println()
	fmt.Printf("export AWS_ACCESS_KEY_ID=%v\n", keyId)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%v\n", secretKey)
	fmt.Printf("export AWS_SESSION_TOKEN=%v\n", sessionToken)
}
