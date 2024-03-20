package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

type Note struct {
	ID          string `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
	Amount      int    `db:"amount" json:"amount"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
}

type ErrorStd struct {
	Message string `json:"message"`
}

type NoDataStd struct {
	Message string `json:"message"`
}

func main() {
	db, err := sqlx.Connect("postgres", "user=notes-misbakh-be-prod dbname=db-notes-misbakh-be-prod sslmode=disable password=tN8Rmp3$8KKBTdcx host=db port=5432")
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(100)
	defer db.Close()

	//schema := `CREATE TABLE notes (
	//				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	//				description TEXT NOT NULL,
	//				amount integer NOT NULL,
	//				created_at timestamptz NOT NULL
	//      	);`

	// execute a query on the server
	//result, err := db.Exec(schema)
	//if err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(&result)

	////add data
	//log.Print(time.Now().Format(time.RFC3339), 23)
	//tx := db.MustBegin()
	//tx.MustExec("INSERT INTO notes (description, amount, created_at) VALUES ($1, $2, $3)", "Beli Makan", 9000, time.Now())
	//tx.MustExec("INSERT INTO notes (description, amount, created_at) VALUES ($1, $2, $3)", "Beli Jajan", 10000, time.Now())
	//tx.Commit()

	//// log data
	//rows, err := db.Queryx("SELECT * FROM notes")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//for rows.Next() {
	//	var notes Note
	//	err := rows.StructScan(&notes)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	log.Println(notes)
	//}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!!!")
	})

	e.POST("/notes", func(c echo.Context) error {
		var notes Note
		if err := c.Bind(&notes); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}

		rows, err := db.Queryx("INSERT INTO notes (description, amount, created_at) VALUES ($1, $2, $3) RETURNING id", notes.Description, notes.Amount, notes.CreatedAt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}

		var id string
		if rows.Next() {
			rows.Scan(&id)
		}
		notes.ID = id

		return c.JSON(http.StatusOK, notes)
	})
	e.GET("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")

		rows, err := db.Queryx("SELECT * FROM notes WHERE id = $1", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}
		var notes Note

		for rows.Next() {
			err := rows.StructScan(&notes)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
			}
		}

		if notes.ID == "" {
			return c.JSON(http.StatusNotFound, ErrorStd{Message: "Data not found"})
		}

		return c.JSON(http.StatusOK, notes)
	})
	e.GET("/notes", func(c echo.Context) error {
		rows, err := db.Queryx("SELECT * FROM notes")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}
		var notes []Note

		for rows.Next() {
			var note Note
			err := rows.StructScan(&note)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
			}
			notes = append(notes, note)
		}

		return c.JSON(http.StatusOK, notes)

	})
	e.PATCH("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")
		var notes Note
		notes.ID = id
		if err := c.Bind(&notes); err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}

		queryBuilder := "UPDATE notes SET "
		if notes.Description != "" {
			queryBuilder += "description = '" + notes.Description + "', "
		}
		if notes.Amount != 0 {
			queryBuilder += "amount = " + strconv.Itoa(notes.Amount) + ", "
		}
		queryBuilder = queryBuilder[:len(queryBuilder)-2]
		queryBuilder += " WHERE id = '" + id + "'"
		log.Print(queryBuilder)
		_, err := db.Queryx(queryBuilder)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, NoDataStd{Message: "Data updated"})
	})
	e.DELETE("/notes/:id", func(c echo.Context) error {
		id := c.Param("id")
		_, err := db.Queryx("DELETE FROM notes WHERE id = $1", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorStd{Message: err.Error()})
		}

		return c.JSON(http.StatusOK, NoDataStd{Message: "Data deleted"})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
