package main

//go get -u github.com/go-sql-driver/mysql
//go run RestApi.go
// cd GinGonic
import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	//_ "github.com/go-sql-driver/mysql"
)

// func createConnection() *sql.DB {
// 	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", "postgres", "root", "people_db")

// 	//var err error
// 	db, err := sql.Open("postgres", connectionString)
// 	//db, err := sql.Open("postgres", "postgres://postgres:7046365527@localhost/postgres?sslmode=disable")

// 	if err != nil {
// 		panic(err)
// 	}
// 	// check the connection
// 	err = db.Ping()

// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Successfully connected!")
// 	return db
// }
func main() {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", "postgres", "root", "people_db")

	//var err error
	db, err := sql.Open("postgres", connectionString)
	//db, err := sql.Open("postgres", "postgres://postgres:7046365527@localhost/postgres?sslmode=disable")

	if err != nil {
		panic(err)
	}
	defer db.Close()
	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	//return db

	// db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/product_db")
	// if err != nil {
	// 	// fmt.Print(err.Error())
	// 	fmt.Println("Error creating DB:", err)
	// 	fmt.Println("To verify, db is:", db)
	// }
	// defer db.Close()
	// fmt.Println("Successfully  Connected to MYSQl")
	// // make sure connection is available
	// err = db.Ping()
	// if err != nil {
	// 	fmt.Print(err.Error())

	type People struct {
		Id      int    `db:"ID" json:"id"`
		Name    string `db:"Name" json:"Name"`
		Address string `db:"Address" json:"Address"`
		Age     int    `db:"age" json:"age"`
	}

	router := gin.Default()

	// GET a people detail
	router.GET("/people/:id", func(c *gin.Context) {
		var people People

		//strconv.Atoi("-42")
		// id := c.Query("id")
		// id1,err = strconv.Atoi(id)
		// id := c.Query("id")
		// id := c.PostForm("id")
		// id := c.Params.ByName("id")
		// id := c.PostForm("id")
		id := c.Param("id")
		//id, _ := strconv.Atoi(id1)
		fmt.Println((id))
		rows, err := db.Query("select * from people where id = $1;", id)
		//rows, err := db.Query("select * from people ;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&people.Id, &people.Name, &people.Address, &people.Age)
			// peoples = append(peoples, people)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": people,
			"count":  1,
		})
	})

	// GET all peoples
	router.GET("/peoples", func(c *gin.Context) {
		var (
			people  People
			peoples []People
		)
		rows, err := db.Query("select * from people;")
		if err != nil {
			fmt.Print(err.Error())
		}
		for rows.Next() {
			err = rows.Scan(&people.Id, &people.Name, &people.Address, &people.Age)
			peoples = append(peoples, people)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
			"result": peoples,
			"count":  len(peoples),
		})

		// c.HTML(
		// 	// Set the HTTP status to 200 (OK)
		// 	http.StatusOK,
		// 	// Use the index.html template
		// 	"index.html",
		// 	// Pass the data that the page uses
		// 	gin.H{
		// 		"title":   "Home Page",
		// 		"payload": peoples,
		// 	},
		// )
	})

	// POST new people details
	router.POST("/people", func(c *gin.Context) {
		var buffer bytes.Buffer
		var people People
		c.Bind(&people)
		// id, err := strconv.Atoi(c.PostForm("id"))
		// fmt.Println("hello", id)
		// //id := c.PostForm("id")
		// Name := c.PostForm("Name")
		// Address := c.PostForm("Address")
		// Age, err := strconv.Atoi(c.PostForm("Age"))
		id := people.Id

		//id := c.PostForm("id")
		Name := people.Name
		Address := people.Address
		Age := people.Age
		//Age := c.PostForm("Age")

		//UPDATED CODE
		// sqlStatement := `INSERT INTO people (id,name, location, age) VALUES ($1, $2, $3,$4) RETURNING ID`
		sqlStatement := `INSERT INTO people (id,Name, Address, age) VALUES ($1, $2, $3,$4) `

		// the inserted id will store in this id
		//var id int64

		// execute the sql statement
		// Scan function will save the insert id in the id
		// err := db.QueryRow(sqlStatement, user.ID, user.Name, user.Location, user.Age).Scan(&id)
		//(1) WAY
		//err := db.QueryRow(sqlStatement, id, Name, Address, Age).Scan(&id)

		if err != nil {
			//log.Fatalf("Unable to execute the query. %v", err)
			fmt.Print(err.Error())
		}

		//fmt.Printf("Inserted a single record %v", id)

		// return the inserted id
		//return id

		//(2) WAY
		stmt, err := db.Prepare(sqlStatement)
		if err != nil {
			fmt.Print(err.Error())
		}
		// _, err = stmt.Exec(&id, &Name, &Address, &Age)
		_, err = stmt.Exec(id, Name, Address, Age)

		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Printf("Inserted a single record %v", id)

		// Fastest way to append strings
		//buffer.WriteString(id)
		buffer.WriteString(" ")
		buffer.WriteString(Name)
		buffer.WriteString(" ")
		buffer.WriteString(Address)
		buffer.WriteString(" ")

		// buffer.WriteString(strconv.Itoa(Age))
		//buffer.WriteString(Age)
		//defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s %ssuccessfully created", Name, name),
		})
	})

	// PUT - update a people details
	router.PUT("/people/:id", func(c *gin.Context) {
		var buffer bytes.Buffer
		// id := c.Query("id")
		// Name := c.Query("Name")
		// Address := c.Query("Address")
		// Age := c.Query("Age")

		id := c.Param("id")
		var people People
		c.Bind(&people)
		// id := people.Id

		//id := c.PostForm("id")
		Name := people.Name
		Address := people.Address
		Age := people.Age
		stmt, err := db.Prepare("update people set Name= $1, Address= $2,Age=$3 where id= $4;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(Name, Address, Age, id)
		if err != nil {
			fmt.Print(err.Error())
		}

		// Fastest way to append strings
		buffer.WriteString(Name)
		buffer.WriteString(" ")
		buffer.WriteString(Address)
		// buffer.WriteString(" ")
		// buffer.WriteString(Age)
		defer stmt.Close()
		name := buffer.String()

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated to %s", name),
		})
	})

	// Delete resources
	router.DELETE("/people/:id", func(c *gin.Context) {
		// id := c.Query("id")

		var people People
		c.Bind(&people)
		// id := people.Id
		id := c.Param("id")
		stmt, err := db.Prepare("delete from people where id= $1;")
		if err != nil {
			fmt.Print(err.Error())
		}
		_, err = stmt.Exec(id)
		if err != nil {
			fmt.Print(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted user: %s", id),
		})
	})
	router.Run(":9000")
}
