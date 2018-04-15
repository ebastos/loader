package main

import (
	"database/sql"
	"flag"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type creator func(int, int, connection, *sync.WaitGroup)

type connection struct {
	user     string
	password string
	host     string
}

// realCreator connects to a MySQL database and created the giving number of DBs and their tables.
func realCreator(n int, tables int, c connection, wg *sync.WaitGroup) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/", c.user, c.password, c.host)
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		fmt.Println("Error connecting to DB!")
		wg.Done()
	}
	defer db.Close()

	// Create empty database
	dbName := fmt.Sprintf("db_%d", n)
	fmt.Printf("Starting job for DB %s\n", dbName)

	_, err = db.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		fmt.Printf("Error creating database %s: %s\n", dbName, err)
		wg.Done()
	}
	_, err = db.Exec("USE " + dbName)
	if err != nil {
		fmt.Printf("Error opening database %s: %s\n", dbName, err)
		wg.Done()
	}

	// Create empty tables. Forcing MYISAM engine.
	for i := 0; i < tables; i++ {
		createTable := fmt.Sprintf("CREATE TABLE table_%d ( id integer, data varchar(32) ) ENGINE = MYISAM", i)
		_, err = db.Exec(createTable)
		if err != nil {
			fmt.Printf("Error creating table table_%d on %s: %s\n", i, dbName, err)
			wg.Done()
		}

	}
	fmt.Printf("Job done for DB %s\n", dbName)
	wg.Done()
}

// fakeCreator allows for a dry-run and helps on troubleshooting and testing. May be moved to a test suite later.
func fakeCreator(n int, tables int, c connection, wg *sync.WaitGroup) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/", c.user, c.password, c.host)
	fmt.Println("Connecting with ", connectionString)
	dbName := fmt.Sprintf("db_%d", n)
	fmt.Printf("Starting job for DB %s\n", dbName)

	fmt.Printf("%d CREATE DATABASE \n", n)
	for i := 0; i <= tables; i++ {
		fmt.Printf("[%s]CREATE TABLE table_%d ( id integer, data varchar(32) )\n", dbName, i)
	}
	wg.Done()
}

// dbCreator keeps tab of database creation process per chunk.
func dbCreator(dbs, tables, created int, dbcreator creator, c connection) int {
	var wg sync.WaitGroup
	for db := 0; db < dbs; db++ {
		wg.Add(1)
		go dbcreator(db+created, tables, c, &wg)
	}
	wg.Wait()
	return created + dbs
}

func main() {
	dbs := flag.Int("dbs", 100, "Number of databases")
	tables := flag.Int("tables", 10, "Tables per database")
	chuncks := flag.Int("chunks", 10, "Number of databases to create per iteraction")

	user := flag.String("user", "root", "Username to connect with - defaults to root")
	password := flag.String("password", "", "Password to use - defaults to empty")
	host := flag.String("host", "localhost:3306", "Hostname and port to connect - defaults to localhost:3306")

	flag.Parse()

	var myConnection connection
	myConnection.host = *host
	myConnection.password = *password
	myConnection.user = *user

	fmt.Printf("Creating %d tables\n", *dbs**tables)

	// De-referencing stuff and creating control variables.
	d := *dbs
	t := *tables
	c := *chuncks
	created := 0

	// For real or test/dry-run?
	// TODO: Write a full test suite
	myCreator := realCreator

	// If creating all DBs fit in a single chunk, just do it.
	if d <= c {
		dbCreator(d, t, created, myCreator, myConnection)
	} else {

		// If not, let's do math and find how to slice it up.
		fullRun := d / c
		missing := d - fullRun*c
		fmt.Println(fullRun)
		fmt.Println(missing)

		for run := 0; run < fullRun; run++ {
			created = dbCreator(c, t, created, myCreator, myConnection)
			fmt.Println("Running full run ", run)
		}
		created = dbCreator(missing, t, created, myCreator, myConnection)
	}

	fmt.Printf("Finished creating %d tables\n", *dbs**tables)

}
