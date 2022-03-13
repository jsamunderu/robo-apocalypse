package survivordb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type SurvivorDB struct {
	DBName               string
	DB                   *sql.DB
	dataDefinitionlStmt  *sql.Stmt
	createStmt           *sql.Stmt
	selectStmt           *sql.Stmt
	selectInfectedStmt   *sql.Stmt
	countInfectedStmt    *sql.Stmt
	selectByIdNumberStmt *sql.Stmt
	updateLocationStmt   *sql.Stmt
	updateResourceStmt   *sql.Stmt
	updateInfectedStmt   *sql.Stmt
}

const (
	ddlSQL = `CREATE TABLE IF NOT EXISTS Survivors (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT,
	age INTEGER,
	gender TEXT,
	id_number TEXT,
	longitude TEXT,
	latitude TEXT,
	water TEXT,
	food TEXT,
	medication TEXT,
	ammunition TEXT,
	infected INTEGER,
	last_ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
	);`
	createSQL           = `INSERT INTO Survivors (name, age, gender, id_number, longitude, latitude, water, food, medication, ammunition, infected) VALUES(?,?,?,?,?,?,?,?,?,?,?);`
	selectSQL           = `SELECT name, age, gender, id_number, longitude, latitude, water, food, medication, ammunition, infected, last_ts FROM Survivors;`
	selectByIdNumberSQL = `SELECT name, age, gender, id_number, longitude, latitude, water, food, medication, ammunition, infected, last_ts FROM Survivors  WHERE id_number = ?;`
	selectInfectedSQL   = `SELECT name, age, gender, id_number, longitude, latitude, water, food, medication, ammunition, infected, last_ts FROM Survivors  WHERE infected = ?;`
	countInfectedSQL    = `SELECT count(*) FROM Survivors  WHERE infected = ?;`

	updateLocationSQL = `UPDATE Survivors SET longitude = ?, latitude = ?, last_ts = CURRENT_TIMESTAMP WHERE id_number = ?`
	updateResourceSQL = `UPDATE Survivors SET water = ?, food = ?, medication = ?, ammunition = ?, last_ts = CURRENT_TIMESTAMP WHERE id_number = ?`
	updateInfectedSQL = `UPDATE Survivors SET infected = 1, last_ts = CURRENT_TIMESTAMP WHERE id_number = ?`
)

func Open(dbName string) *SurvivorDB {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return nil
	}

	return &SurvivorDB{
		DBName: dbName,
		DB:     db,
	}
}

func (s *SurvivorDB) Setup() error {
	ddlStmt, err := s.DB.Prepare(ddlSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	s.dataDefinitionlStmt = ddlStmt

	_, err = s.dataDefinitionlStmt.Exec()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   ddlSQL,
		}).Info("Sql error")
		return err
	}

	createStmt, err := s.DB.Prepare(createSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	selectStmt, err := s.DB.Prepare(selectSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	selectInfectedStmt, err := s.DB.Prepare(selectInfectedSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	countInfectedStmt, err := s.DB.Prepare(countInfectedSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	selectByIdNumberStmt, err := s.DB.Prepare(selectByIdNumberSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	updateLocationStmt, err := s.DB.Prepare(updateLocationSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	updateResourceStmt, err := s.DB.Prepare(updateResourceSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}

	updateInfectedStmt, err := s.DB.Prepare(updateInfectedSQL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
		}).Info("Sql error")
		return err
	}
	s.createStmt = createStmt
	s.selectStmt = selectStmt
	s.selectInfectedStmt = selectInfectedStmt
	s.countInfectedStmt = countInfectedStmt
	s.selectByIdNumberStmt = selectByIdNumberStmt
	s.updateLocationStmt = updateLocationStmt
	s.updateResourceStmt = updateResourceStmt
	s.updateInfectedStmt = updateInfectedStmt
	return nil
}

// Save inserts a survivor into the Survivors table
func (s *SurvivorDB) Save(survivor *Survivor) error {
	_, err := s.createStmt.Exec(survivor.Name,
		survivor.Age,
		survivor.Gender,
		survivor.IdNumber,
		survivor.Longitude,
		survivor.Latitude,
		survivor.Water,
		survivor.Food,
		survivor.Medication,
		survivor.Ammunition,
		survivor.Infected)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   createSQL,
		}).Info("Sql error")
		return err
	}

	return nil
}

// UpdateLocation updates a survivor location in the Survivors table
func (s *SurvivorDB) UpdateLocation(idNumber string, longitude, latitude float64) error {
	_, err := s.updateLocationStmt.Exec(
		longitude,
		latitude,
		idNumber,
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   updateLocationSQL,
		}).Info("Sql error")
		return err
	}

	return nil
}

// UpdateResource updates a survivor resouce in the Survivors table
func (s *SurvivorDB) UpdateResource(idNumber string, water float64, food, medication string, ammunition int) error {
	_, err := s.updateResourceStmt.Exec(
		water,
		food,
		medication,
		ammunition,
		idNumber,
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   updateResourceSQL,
		}).Info("Sql error")
		return err
	}

	return nil
}

// UpdateResource updates a survivor resouce in the Survivors table
func (s *SurvivorDB) UpdateInfected(idNumber string) error {
	_, err := s.updateInfectedStmt.Exec(idNumber)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   updateResourceSQL,
		}).Info("Sql error")
		return err
	}

	return nil
}

// GetAllSurvivors selects all survivors stored in the Survivors table
func (s *SurvivorDB) GetAllSurvivors() []Survivor {
	rows, err := s.selectStmt.Query()
	defer rows.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   selectSQL,
		}).Info("Sql error")
		return nil
	}

	survivors := []Survivor{}
	for rows.Next() {
		survivor := Survivor{}
		err = rows.Scan(&survivor.Name,
			&survivor.Age,
			&survivor.Gender,
			&survivor.IdNumber,
			&survivor.Longitude,
			&survivor.Latitude,
			&survivor.Water,
			&survivor.Food,
			&survivor.Medication,
			&survivor.Ammunition,
			&survivor.Infected,
			&survivor.LastUpdateTime)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error": err,
				"sql":   selectSQL,
			}).Info("Sql error")
			return nil
		}
		survivors = append(survivors, survivor)
	}
	err = rows.Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   selectSQL,
		}).Info("Sql error")
		return nil
	}

	return survivors
}

// GetInfectedSurvivors selects all infected or uninfected survivors stored in the Survivors table
func (s *SurvivorDB) GetSurvivors(infected bool) []Survivor {
	rows, err := s.selectInfectedStmt.Query(infected)
	defer rows.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   selectInfectedSQL,
		}).Info("Sql error")
		return nil
	}

	survivors := []Survivor{}
	for rows.Next() {
		survivor := Survivor{}
		err = rows.Scan(&survivor.Name,
			&survivor.Age,
			&survivor.Gender,
			&survivor.IdNumber,
			&survivor.Longitude,
			&survivor.Latitude,
			&survivor.Water,
			&survivor.Food,
			&survivor.Medication,
			&survivor.Ammunition,
			&survivor.Infected,
			&survivor.LastUpdateTime)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Error": err,
				"sql":   selectInfectedSQL,
			}).Info("Sql error")
			return nil
		}
		survivors = append(survivors, survivor)
	}
	err = rows.Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   selectSQL,
		}).Info("Sql error")
		return nil
	}

	return survivors
}

// CountSurvivors count all infected or uninfected survivors stored in the Survivors table
func (s *SurvivorDB) CountSurvivors(infected bool) int {
	rows, err := s.countInfectedStmt.Query(infected)
	defer rows.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   countInfectedSQL,
		}).Info("Sql error")
		return -1
	}

	count := 0
	if rows.Next() {
		err = rows.Scan(&count)
	} else {
		return -1
	}
	err = rows.Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   countInfectedSQL,
		}).Info("Sql error")
		return -1
	}

	return count
}

// GetSurvivor selects all survivors stored in the Survivors table
func (s *SurvivorDB) GetSurvivor(idNumber string) *Survivor {
	rows, err := s.selectByIdNumberStmt.Query(idNumber)
	defer rows.Close()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   selectByIdNumberSQL,
		}).Info("Sql error")
		return nil
	}

	survivor := Survivor{}
	if rows.Next() {
		err = rows.Scan(&survivor.Name,
			&survivor.Age,
			&survivor.Gender,
			&survivor.IdNumber,
			&survivor.Longitude,
			&survivor.Latitude,
			&survivor.Water,
			&survivor.Food,
			&survivor.Medication,
			&survivor.Ammunition,
			&survivor.Infected,
			&survivor.LastUpdateTime)
	} else {
		return nil
	}
	err = rows.Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error": err,
			"sql":   selectSQL,
		}).Info("Sql error")
		return nil
	}

	return &survivor
}
