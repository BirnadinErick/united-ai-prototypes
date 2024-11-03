/*
Copyright 2024 Birnadin Erick, Anton Robert Clive

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type Template struct {
	templates *template.Template
}

type Project struct {
	Id       int
	Name     string
	Lead     string
	Deadline string
	By       string
	Desc     string
	Cover    string
}

type Interest struct {
	Email    string `json:"email" form:"email"`
	Course   string `json:"course" form:"course"`
	IsMember bool   `json:"is_member" form:"is_member"`
	Project  uint8  `json:"project" form:"project"`
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./uai.db")
	if err != nil {
		return nil, err
	}

	// init tables
	initStmt := `create table if not exists interests(
		id integer primary key autoincrement,
		email text,
		is_member integer,
		project integer,
		course string
	);
	create table if not exists projects(
		id integer not null primary key,
		name text,
		lead text,
		deadline text,
		by text,
		desc text,
		cover text
	);`
	_, err = db.Exec(initStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {

	e := echo.New()
	e.HideBanner = true

	// setup logger
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("could not open log file:", err)
		fmt.Print("failed to open logger")
	}
	defer logFile.Close()
	e.Logger.SetOutput(logFile)
	e.Use(middleware.Logger())

	// bind templates
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// static asset
	if ex, err := os.Executable(); err != nil {
		e.Logger.Fatal(err)
	} else {
		dir := filepath.Dir(ex)
		e.File("style.css", filepath.Join(dir, "style.css"))
	}

	// init db
	// creates interests and projects tables
	db, err := InitDB()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	e.GET("/", func(c echo.Context) error {
		// get projects
		var projects []Project
		rows, err := db.Query("select id, name, lead, deadline, by, desc, cover from projects")
		if err != nil {
			e.Logger.Debug("couldn't query")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}
		defer rows.Close()

		for rows.Next() {
			var i int
			var n, l, d, b, desc, cover string

			err = rows.Scan(&i, &n, &l, &d, &b, &desc, &cover)
			if err != nil {
				e.Logger.Debug("failed to scan into variables")
				e.Logger.Fatal(err)
				return echo.ErrInternalServerError
			}

			projects = append(projects, Project{i, n, l, d, b, desc, cover})
		}

		// check for any silent errors
		err = rows.Err()
		if err != nil {
			e.Logger.Debug("err occured during row scans")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}

		return c.Render(http.StatusOK, "index", projects)
	})

	e.GET("/declare-interest/:project_id", func(c echo.Context) error {
		var project Project
		var i int
		var n, d, b, l, desc string
		p, err := strconv.Atoi(c.Param("project_id"))
		if err != nil {
			e.Logger.Debug("wrong project id")
			e.Logger.Fatal(err)
			return echo.ErrBadRequest
		}

		stmt, err := db.Prepare("select id, name, lead, deadline, by, desc from projects where id = ?")
		if err != nil {
			e.Logger.Debug("couldn't prepare query")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}
		defer stmt.Close()

		err = stmt.QueryRow(p).Scan(&i, &n, &l, &d, &b, &desc)
		if err != nil {
			e.Logger.Debug("couldn't bind results")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}
		project = Project{i, n, l, d, b, desc, ""}

		return c.Render(http.StatusOK, "interest-form", project)
	})

	e.POST("/interest", func(c echo.Context) error {
		i := new(Interest)
		if err := c.Bind(i); err != nil {
			return echo.ErrBadRequest
		}

		// record the interest
		tx, err := db.Begin()
		if err != nil {
			e.Logger.Debug("failed to begin a db tx")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}

		stmt, err := tx.Prepare("insert into interests(email, course, is_member, project) values(?, ?, ?, ?)")
		if err != nil {
			e.Logger.Debug("failed to prepare insert stmt")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}
		defer stmt.Close()

		_, err = stmt.Exec(i.Email, i.Course, i.IsMember, i.Project)
		if err != nil {
			e.Logger.Debug("failed to execute stmt")
			e.Logger.Fatal(err)
			log.Fatal(err)
			return echo.ErrInternalServerError
		}

		err = tx.Commit()
		if err != nil {
			e.Logger.Debug("failed to commit tx")
			e.Logger.Fatal(err)
			return echo.ErrInternalServerError
		}

		return c.Render(http.StatusAccepted, "interest-ok", "")
	})

	// render by default binds to 10000 PORT
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "10000"
	}
	e.Logger.Fatal(e.Start("0.0.0.0:" + PORT))
}
