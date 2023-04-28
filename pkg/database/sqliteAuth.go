package database

import (
	"database/sql"
	//"fmt"
	"log"

	//"github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	//"fmt"
	"time"
)

type Cookie struct {
	Id        int
	Ipaddress string
	Phpsessid string
	Time      string
}


// This function retrieves the session cookie from the sqlite database
//if 8 mins have passed since last sessioncookie was retrieved, fetch new, else return the last one.
func GetSqliteKeyIfNotExpired(db *sql.DB, ipaddress string) (cookie string, err error) {

	/*var sqliteDatabase *sql.DB

	sqliteDatabase, err = sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Print(err)
		return "", err
	}*/
	if (!RowExists(db, ipaddress)) {
		return "",nil
	}
	Hosts, err := GetCookieDB(db, ipaddress)
	if err != nil {
		log.Print(err)
		return "", err
	}

	var c string
	mins := time.Minute*time.Duration(9)
	for i := range Hosts {
		if Hosts[i].Ipaddress == ipaddress {
			now := time.Now().Format(time.RFC3339)
			parsed, _ := time.Parse(time.RFC3339, now)
			parsed2, err := time.Parse(time.RFC3339, Hosts[i].Time)
			if err != nil {
				log.Print(err)
				return "", nil
			}
			if parsed2.Add(mins).After(parsed) == true {
				c = Hosts[i].Phpsessid
				return c, nil
			} else {
				return "", nil
			}
		} else {
			return "", nil
		}
	}
	return c, nil
}
// Here starts functions concerning sessioncookies
// We are passing db reference connection from main to our method with other parameters
func InsertAuth(db *sql.DB, ipaddress string, phpsessid string, time string) error{
	log.Println("Inserting session cookie ...")
	insertAuthSQL := `INSERT INTO authentication(ipaddress, phpsessid, time) VALUES (?, ?, ?)`

	statement, err := db.Prepare(insertAuthSQL) // Prepare statement.
                                                   // This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(ipaddress, phpsessid, time)
	if err != nil {
		return err
	}
	return nil
}
func GetCookieDB(db *sql.DB, ipaddress string) ([]*Cookie, error){
	row, err := db.Query("SELECT * FROM authentication")
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer row.Close()

	var c []*Cookie
	for row.Next() {
			p := &Cookie{}
			if err := row.Scan(&p.Id, &p.Ipaddress, &p.Phpsessid, &p.Time); err != nil{
			}
			if (p.Ipaddress == ipaddress) {
				c = append(c, p)

			}
	}
	/*for i :=range c{
		fmt.Println(c[i].Ipaddress, c[i].Phpsessid, c[i].Time)
	}*/
	return c,err
}

func CreateAuthTable(db *sql.DB) error {
	createAuthTableSQL := `CREATE TABLE IF NOT EXISTS authentication (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"ipaddress" TEXT,
		"phpsessid" TEXT,
		"time" TEXT
	  );` // SQL Statement for Create Table

	//log.Println("Create table...")
	statement, err := db.Prepare(createAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	statement.Exec() // Execute SQL Statements
	//log.Println("table created")
	return nil
}

func DropTable(db *sql.DB) error{
	dropAuthTableSQL := `DROP TABLE IF EXISTS authentication`
	statement, err := db.Prepare(dropAuthTableSQL) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}// Execute SQL Statements
	log.Print("table dropped")
	return nil
}



func RowExists(db * sql.DB, ip string) bool {
    sqlStmt := `SELECT ipaddress FROM authentication WHERE ipaddress = ?`
    err := db.QueryRow(sqlStmt, ip).Scan(&ip)
    if err != nil {
        if err != sql.ErrNoRows {
            // a real error happened! you should change your function return
            // to "(bool, error)" and return "false, err" here
            log.Print(err)
        }

        return false
    }

    return true
}

func Update(db *sql.DB,  phpsessid string, time string, ipaddress string) {
	stmt, err := db.Prepare("UPDATE authentication set phpsessid=?, time=? WHERE ipaddress=?")
	if err != nil {
	 //log.Fatal(err)
	 log.Print("update",err)
	}
	res, err := stmt.Exec(phpsessid, time, ipaddress)
	if err != nil {
	 //log.Fatal(err)	 log.Print("update",err)

	}
	affected, _ := res.RowsAffected()
	log.Printf("Affected rows %d", affected)
   }


