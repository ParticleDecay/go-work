package database

import (
	"database/sql"
	"fmt"

	"github.com/ParticleDecay/go-work/pkg/tmux"
	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	log "github.com/sirupsen/logrus"
)

// Project represents an instance of a project environment.
type Project struct {
	Name   string
	Path   string
	Goroot string
	Gopath string
}

// Database allows access to the database.
type Database struct {
	Conn *sql.DB
	Path string
}

// Open returns a connection to the database.
func (d *Database) Open() {
	if d.Path == "" {
		log.Fatal("Cannot open database without a path to database file")
	}
	log.Debugf("Using database file %s", d.Path)
	db, err := sql.Open("sqlite3", d.Path)
	if err != nil {
		log.Fatal(err)
	}
	d.Conn = db
}

// Prep creates the database and table if it does not already exist.
func (d *Database) Prep() bool {
	if d.Conn == nil {
		d.Open()
	}
	_, err := d.Conn.Exec(`CREATE TABLE IF NOT EXISTS projects (
							name VARCHAR(100),
							path VARCHAR(255),
							goroot VARCHAR(255),
							gopath VARCHAR (255));`)
	return err == nil
}

// GetProjects retrieves projects from the database with the specified name
func (d *Database) GetProjects(name string) []*Project {
	if d.Conn == nil {
		d.Open()
	}
	var rows *sql.Rows
	if name != "" {
		query, err := d.Conn.Prepare(`SELECT * FROM projects WHERE name=?`)
		if err != nil {
			log.Fatal(err)
		}
		rows, err = query.Query(name)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		rows, err = d.Conn.Query(`SELECT * FROM projects`)
		if err != nil {
			log.Fatal(err)
		}
	}

	var results []*Project
	for rows.Next() {
		p := &Project{}
		err := rows.Scan(&p.Name, &p.Path, &p.Goroot, &p.Gopath)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, p)
	}
	rows.Close()

	return results
}

// NewProject adds a new project to the database.
func (d *Database) NewProject(name, path, goroot, gopath string) bool {
	if d.Conn == nil {
		d.Open()
	}
	update, err := d.Conn.Prepare(`INSERT INTO projects VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	updated, err := update.Exec(name, path, goroot, gopath)
	if err != nil {
		log.Fatal(err)
	}
	affected, err := updated.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return affected > 0
}

// SelectProject selects a project for launching an environment.
func (d *Database) SelectProject(name string) {
	projects := d.GetProjects(name)
	var selectedProject *Project

	// If multiple projects match, have user select one.
	if len(projects) > 1 {
		templates := promptui.SelectTemplates{
			Active:   `{{ .Path | cyan | bold }}`,
			Inactive: `{{ .Path | cyan }}`,
			Selected: `{{ "✔" | green | bold }} {{ .Path | cyan }}`,
		}
		prompt := promptui.Select{
			Label:     "Existing projects:",
			Items:     projects,
			Templates: &templates,
		}
		i, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}
		selectedProject = projects[i]
	} else if len(projects) == 1 {
		selectedProject = projects[0]
	} else {
		log.Fatal(fmt.Sprintf("'%s' project not found. Try adding it with the 'add' command.", name))
	}
	log.Debugf("Selected project '%s' at %s", selectedProject.Name, selectedProject.Path)

	tmux.LaunchEnvironment(selectedProject.Name, selectedProject.Path, selectedProject.Goroot, selectedProject.Gopath)
}
