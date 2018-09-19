package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/rds"
)

type Endpoint struct {
	Port    int
	Address string
}
type DBInstance struct {
	Endpoint           Endpoint
	Name, Size, Engine string
	ConnectionString   string
	Account            string
	CreatedAt          string
}

const TimeFormat = "2006-01-02 15:04"

func createConnectionString(instance *rds.DBInstance) string {
	switch *instance.Engine {
	case "postgres":
		return fmt.Sprintf("postgresql://%v:PASSWORD@%v:%v/%v", *instance.MasterUsername, *instance.Endpoint.Address, *instance.Endpoint.Port, *instance.DBName)
	default:
		return ""
	}

}

func rdsToDb(instances []*rds.DBInstance) []*DBInstance {
	var dbs []*DBInstance
	for _, instance := range instances {
		db := DBInstance{
			Endpoint: Endpoint{
				Port:    int(*instance.Endpoint.Port),
				Address: *instance.Endpoint.Address,
			},
			Name:             *instance.DBInstanceIdentifier,
			Size:             *instance.DBInstanceClass,
			Engine:           *instance.Engine,
			ConnectionString: createConnectionString(instance),
			Account:          *instance.MasterUsername,
			CreatedAt:        instance.InstanceCreateTime.Format(TimeFormat),
		}
		dbs = append(dbs, &db)
	}
	return dbs

}
func printConnectionCommand(instance *DBInstance) {
	var resp string
	switch instance.Engine {
	case "mysql":
		resp += fmt.Sprintln("Connection Command:")
		resp += fmt.Sprintf("mysql --host %v -u %v --port %v -p", instance.Endpoint.Address, instance.Account, instance.Endpoint.Port)
		resp += fmt.Sprintln("")
		break
	case "postgres":
		resp = fmt.Sprintln("Connection String")
		resp += fmt.Sprintf("%v", instance.ConnectionString)
		resp += fmt.Sprintln("")
		resp += fmt.Sprintln("Connection Command:")
		resp += fmt.Sprintf("psql --host %v -U %v --port %v -W", instance.Endpoint.Address, instance.Account, instance.Endpoint.Port)
		resp += fmt.Sprintln("")
		break
	}
	fmt.Println(resp)
}

func printDBS(dbs []*DBInstance) {
	writer := tabwriter.NewWriter(os.Stdout, 4, 4, 4, ' ', tabwriter.TabIndent)
	fmt.Fprintln(writer, "Number\tDBName\tEngine\tCreatedAt\tEndpoint")
	fmt.Fprintln(writer, "------\t---------\t----------\t----------\t---------")
	for i, instance := range dbs {
		fmt.Fprintf(writer, "[%v]\t%s\t%s\t%s\t%s\t\n", i+1, instance.Name, instance.Engine, instance.CreatedAt, instance.Endpoint.Address)
	}
	writer.Flush()
	fmt.Fprintln(os.Stdout, "Select a Number to get the connection command")
	var instanceNo string
	_, err := fmt.Scanln(&instanceNo)
	if err != nil {
		log.Fatal(err)
	}
	var index int
	if index, err = strconv.Atoi(instanceNo); err != nil {
		fmt.Println("please select a number")
		log.Fatal(err)
	}
	if index > len(dbs) {
		fmt.Println("please select an index in range")
	} else {
		instance := dbs[index-1]
		printConnectionCommand(instance)
	}
	os.Exit(1)
}
