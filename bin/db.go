package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

//ConnectTODatabase connection
func ConnectTODatabase() {
	db, err = sql.Open("postgres", "host="+DBConfigInterface.DBHost+" port="+strconv.Itoa(DBConfigInterface.DBPort)+" user="+DBConfigInterface.DBUser+" password="+DBConfigInterface.DBPass+" dbname="+DBConfigInterface.DBName+" sslmode=disable")
	if err != nil {
		errstr := "WriteToDB(): Couldn't connect to db:" + err.Error()
		fmt.Println(errstr)
	}

	// err = db.Ping()
	// if err != nil {
	// 	errstr := `Error: Could not establish a connection with the database`
	// 	fmt.Println(errstr + err.Error())
	// }
}

func istherenumber(plate_number string) (have bool) {

	query := `
	SELECT islastans
	FROM vehiclerouting 
	where plate_number  = '` + plate_number + `'`

	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}

	rows, err := db.Query(query)
	if err != nil {
		errstr := "istherenumber(): sql query problem: " + err.Error()
		PrintlnToLogFile(errstr)
		fmt.Println(errstr, query)
	}
	defer rows.Close()

	if rows.Next() {
		have = true
	} else {
		have = false
	}
	return have
}

func SelectFromDBEntrance(plate, ansid string) (ent int, lastans bool, fullplate string) {
	var islastans bool

	query := `SELECT islastans, plate 
FROM vehiclerouting 
where plate_number = ` + plate + ` and ansid = ` + ansid + ` order by idkey desc limit 1;
`

	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}

	rows, err := db.Query(query)
	if err != nil {
		errstr := "SelectFromDBEntrance(): sql query problem: " + err.Error()
		PrintlnToLogFile(errstr)
		fmt.Println(errstr, query)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&islastans, &fullplate)
		if err != nil {
			errstr := "SelectFromDBEntrance(): db rows reading problem: " + err.Error()
			PrintlnToLogFile(errstr)
			fmt.Println(errstr)
		}
	}
	if islastans == false {
		ent = 1
		lastans = false
	} else {
		ent = 0
		lastans = true
	}
	return ent, lastans, fullplate
}
func cutnumber(plate string) (cutplate string) {
	if _, err = strconv.Atoi(plate[:3]); err == nil {
		return plate[:3]
	} else if _, err = strconv.Atoi(plate[1:4]); err == nil {
		return plate[:3]
	} else if _, err = strconv.Atoi(plate[2:5]); err == nil {
		return plate[:3]
	} else {
		return plate
	}
}

//InsertToDBNewPlate add new vehicle information to DB
func InsertToDBNewPlate(ansId int, plate, dstpoint string, ent bool) error {
	smalPlate := cutnumber(plate)
	sqlStatement := `INSERT INTO public.vehiclerouting(ansid, plate, destination,  islastans, plate_number)
VALUES (` + strconv.Itoa(ansId) + `,'` + plate + `', '` + dstpoint + `', ` + strconv.FormatBool(ent) + `, ` + smalPlate + `);
`
	//Check connection
	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}
	_, err = db.Exec(sqlStatement)
	if err != nil {
		errstr := "InsertToDBNewPlate(): SQL query problem: " + err.Error()
		fmt.Println(errstr, sqlStatement)
		PrintlnToLogFile(`InsertToDBNewPlate(): ` + sqlStatement + `SQL query problem: ` + err.Error())
		return err
	}
	return nil
}

func UpdateDBinformation(plate, direction, dstpoint string, islastans bool, ansid int) error {
	var sqlStatement string
	if islastans == true {
		log.Println(`Updating destination point to exit plate: `, plate)
		UpdateDestination(plate, `1020`)
		sqlStatement = `UPDATE public.vehiclerouting
		SET   ansid=` + strconv.Itoa(ansid) + `, destination='1020', direction='` + direction + `', islastans='true'
		WHERE  plate_number= ` + plate + `;
		`
	} else {
		sqlStatement = `UPDATE public.vehiclerouting
		SET   ansid=` + strconv.Itoa(ansid) + `, destination='` + dstpoint + `', direction='` + direction + `' 
		WHERE plate_number= ` + plate + `;
		`
	}
	// Check connection
	// err = db.Ping()
	// if err != nil {
	// 	ConnectTODatabase()
	// }
	_, err = db.Exec(sqlStatement)
	if err != nil {
		errstr := "UpdateDBinformation(): SQL query problem: " + err.Error()
		fmt.Println(errstr, sqlStatement)
		PrintlnToLogFile(`UpdateDBinformation(): ` + sqlStatement + `SQL query problem: ` + err.Error())
		return err
	}
	return nil
}

func InsertNewDestination(plate, dest, plate_number string) {
	sqlStatement := `
	INSERT INTO vehicle_u_destination(
		 plate, destination, plate_number )
VALUES ( '` + plate + `', '` + dest + `', ` + plate_number + `);
	`

	//Check connection
	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}
	_, err = db.Exec(sqlStatement)
	if err != nil {
		errstr := "InsertNewDestination(): SQL query problem: " + err.Error()
		fmt.Println(errstr, sqlStatement)
		PrintlnToLogFile(`InsertNewDestination(): ` + sqlStatement + `SQL query problem: ` + err.Error())
	}
}

func SelectDestFromLocalDB(plate string) (destination string, err error) {
	query := `SELECT  destination
	FROM vehicle_u_destination
	WHERE plate_number  ='` + plate + `'
	order by id desc limit 1 ;`

	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}

	rows, err := db.Query(query)
	if err != nil {
		errstr := "SelectDestFromLocalDB(): sql query problem: " + err.Error()
		PrintlnToLogFile(errstr)
		fmt.Println(errstr, query)
		return destination, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&destination)
		if err != nil {
			errstr := "SelectDestFromLocalDB(): db rows reading problem: " + err.Error()
			PrintlnToLogFile(errstr)
			fmt.Println(errstr)
			return destination, err
		}
	}
	return destination, nil
}

func UpdateDestination(plate, destination string) {
	sqlStatement := `	
	UPDATE vehicle_u_destination
  	SET destination='` + destination + `'
 	WHERE plate='` + plate + `'
	`

	//Check connection
	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}
	_, err = db.Exec(sqlStatement)
	if err != nil {
		errstr := "UpdateDestination(): SQL query problem: " + err.Error()
		fmt.Println(errstr, sqlStatement)
		PrintlnToLogFile(`UpdateDestination(): ` + sqlStatement + `SQL query problem: ` + err.Error())
	}
}

func SelectFullPlate(cutplate string) (fullplate string, err error) {
	sqlselect := `
	SELECT  plate
	FROM public.vehicle_u_destination
	where plate_number =` + cutplate + `
	order by id desc limit 1`

	err = db.Ping()
	if err != nil {
		ConnectTODatabase()
	}

	rows, err := db.Query(sqlselect)
	if err != nil {
		errstr := "SelectFullPlate(): sql query problem: " + err.Error()
		PrintlnToLogFile(errstr)
		fmt.Println(errstr, sqlselect)
		return fullplate, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&fullplate)
		if err != nil {
			errstr := "SelectFullPlate(): db rows reading problem: " + err.Error()
			PrintlnToLogFile(errstr)
			fmt.Println(errstr)
			return fullplate, err
		}
	}
	return fullplate, nil
}

func CleanDB() {
	for {
		time.Sleep(24 * time.Hour)
		sql1 := `delete from  public.vehicle_u_destination`
		sql2 := `delete from public.vehiclerouting`
		_, err = db.Exec(sql1)
		_, err = db.Exec(sql2)
		if err != nil {
			log.Printf(`couldn't clean db`, err)
		}
	}
}
