# rds2conn
Generate database connection strings from RDS with grace and ease

## Installation

```
$ go get github.com/adhocteam/rdsconn
```

to upgrade your rdsconn version, do:

```
$ go get -u github.com/adhocteam/rdsconn
```

## Usage

### Displaying

```
$ rdsconn -d
```

### Other options/flags

Specify an AWS profile other than "default":

```
$ rdsconn --profile profilenamehere
```

## Notes

- Assumes you have at least one AWS profile configured. See [AWS docs for details](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html#cli-quick-configuration).