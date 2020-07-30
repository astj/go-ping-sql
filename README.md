# go-ping-sql

SQL ping for mysql/postgres.

## prepare

If you want to test with RDS TLS connection, you need to download certificate file according to https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html .

## mysql

### without tls

```
go run main.go mysql "user:pass@tcp(your.database.region.rds.amazonaws.com:3306)/database"
```

### with tls required, but without specifying CA certificate file

This should fail all RDS because of lacking root CA file
```
go run main.go mysql "user:pass@tcp(your.database.region.rds.amazonaws.com:3306)/database?tls=true"
```

### with TLS required and CA certificate file

```
go run main.go --ssl-ca rds-ca-2019-root.pem mysql "user:pass@tcp(your.database.region.rds.amazonaws.com:3306)/database?tls=custom"
```

This should pass when:
- Target RDS instance has configured to use `rds-ca-2019` cert file
- One of following conditions are met:
  - prior to Go 1.15
  - Go 1.15 and target DB instance that was created or updated to the rds-ca-2019 certificate AFTER July 28, 2020

Which means, if the instance was created or updated to the rds-ca-2019 certificate prior to July 28, 2020, behavior will change at Go 1.15.

In such cases, you'll see following error messages:
```
go run main.go --ssl-ca rds-ca-2019-root.pem mysql "user:pass@tcp(your.database.region.rds.amazonaws.com:3306)/database?tls=custom"
2020/07/30 19:53:03 x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0
exit status 1
```

## postgres

### sslmode=require

This should pass, as long as the instance is using `rds-ca-2019`

```
PGPASSWORD=xxx go run main.go postgres "user=xxx dbname=xxx sslmode=require host=your.db.region.rds.amazonaws.com"
```

### sslmode=verify-ca

```
PGSSLROOTCERT=rds-ca-2019-root.pem PGPASSWORD=xxx go run main.go postgres "user=xxx dbname=xxx sslmode=verify-ca host=your.db.region.rds.amazonaws.com"
```

This should also pass, but this requires CA cert file by `PGSSLROOTCERT` env.
Otherwise you'll see:

```
PGPASSWORD=xxx go run main.go postgres "user=xxx dbname=xxx sslmode=verify-ca host=your.db.region.rds.amazonaws.com"
2020/07/31 05:10:27 x509: certificate signed by unknown authority
exit status 1
```

### sslmode=verify-full

```
PGSSLROOTCERT=rds-ca-2019-root.pem PGPASSWORD=xxx go run main.go postgres "user=xxx dbname=xxx sslmode=verify-full host=your.db.region.rds.amazonaws.com"
```

This will pass when the same condition as "with TLS required and CA certificate file" of MySQL.
