# Loader

A simple MySQL/MariaDB table generator to load-test a server.
Will create a given number of databases and tables per database.

It uses go routines to run creation in parallel. You can control the number of parallel jobs via ```--chunks```.


### Usage


```
$ ./loader --help
Usage of ./loader:
  -chunks int
        Number of databases to create per iteraction (default 10)
  -dbs int
        Number of databases (default 100)
  -host string
        Hostname and port to connect - defaults to localhost:3306 (default "localhost:3306")
  -password string
        Password to use - defaults to empty
  -tables int
        Tables per database (default 10)
  -user string
        Username to connect with - defaults to root (default "root")

```

### Example

```
$ ./loader --dbs=3 --tables=12 --user=admin --password=password1 
Creating 36 tables
Starting job for DB db_2
Starting job for DB db_0
Starting job for DB db_1
Job done for DB db_1
Job done for DB db_2
Job done for DB db_0
Finished creating 36 tables
```