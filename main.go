package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [options] {instance id|private IPv4 address|name}
Options:
  -v	        be verbose
  -P, --profile	AWS profile to be used for the session (optional)
  -d, --database        show rds db instances and connection info
`, filepath.Base(os.Args[0]))
	os.Exit(1)
}

var verboseFlag bool
var listDB bool
var profile string

func debugf(format string, args ...interface{}) {
	if verboseFlag {
		log.Printf(format, args...)
	}
}

func printError(err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		log.Println("Error:", awsErr.Code(), awsErr.Message())
	} else {
		log.Println("Error:", err.Error())
	}
	os.Exit(1)
}

func init() {
	// default key path to home dir, inherit if env var if set
	flag.BoolVar(&verboseFlag, "v", false, "be verbose")

	const (
		defaultProfile = "default"
		profileDesc    = "the named profile to use when creating a new AWS session"
	)
	flag.StringVar(&profile, "P", defaultProfile, profileDesc)
	flag.StringVar(&profile, "profile", defaultProfile, profileDesc)

	flag.BoolVar(&listDB, "d", false, "List DB Information")
	flag.BoolVar(&listDB, "database", false, "List DB info")
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")

	flag.Usage = usage
	flag.Parse()

	debugf("using AWS profile: %s", profile)
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			CredentialsChainVerboseErrors: &verboseFlag,
		},
		Profile: profile,
	})
	if err != nil {
		printError(err)
	}
	if listDB {
		dbSession := rds.New(sess)
		var params *rds.DescribeDBInstancesInput
		params = &rds.DescribeDBInstancesInput{}
		resp, err := dbSession.DescribeDBInstances(params)
		if err != nil {
			printError(err)

		}
		printDBS(rdsToDb(resp.DBInstances))
	}
}
