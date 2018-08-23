package vocaber

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/jmoiron/sqlx"
	"log"
)

type VocabItem struct {
	Id int `json:"id"`
	Value string `json:"value"`
	Created time.Time `json:"created"`
	Knownit int`json:"knownit"`
}

var dbUrl string = "nautilis:nautilis123@tcp(127.0.0.1:3306)/vocabulary?parseTime=true"

func getDb()(*sqlx.DB, error){
	db, err := sqlx.Open("mysql", dbUrl)
	if err != nil{
		log.Fatal(err)
		return nil, err;
	}
	return db, nil;
}

func getTableName() string{
	return "vocab_item"
}

func Save(item *VocabItem) error{
	db, err:= getDb()
	defer db.Close()
	if err != nil{
		return err
	}
	sql := " INSERT INTO " + getTableName() + " (value, created, knownit) VALUES (:value, :created, :knownit)"
	tx := db.MustBegin()
	tx.NamedExec(sql, item)
	tx.Commit()
	return nil
}

func Count(startDate time.Time, endDate time.Time) (int, error){
	db, err := getDb()
	defer db.Close()
	if err != nil {
		return 1, err
	}
	sql := " SELECT count(1) FROM " + getTableName() + " WHERE created >= ? AND created <= ?"
	var count int
	err = db.Get(&count, sql, startDate, endDate)
	if err != nil{
		return 1, err
	}
	return count, nil
}

func Know(id int) (bool, error){
	db, err := getDb()
	defer db.Close()
	if err != nil {
		return false, err;
	}

	sql := " UPDATE " + getTableName() + " SET knownit = knownit + 1 WHERE id = ? "
	db.MustExec(sql, id)
	return true, nil
}

func GetNoMaster() ([]VocabItem,error) {
	db, err := getDb()
	defer db.Close()
	if err != nil{
		return nil, err
	}
	sql := " SELECT * FROM " + getTableName() + " WHERE knownit < 10 "
	items := []VocabItem{}
	err = db.Select(&items, sql)
	if err != nil{
		return nil, err
	}
	return items, nil
}

func GetByDate(date time.Time) ([]VocabItem, error){
	db, err := getDb()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	dayBegin := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	dayEnd := dayBegin.Add(time.Hour * 24)
	sql := " SELECT * FROM " + getTableName() + " WHERE created >= ? AND created <= ?"
	items := []VocabItem{}
	err = db.Select(&items, sql, dayBegin, dayEnd)
	if err != nil{
		return nil, err
	}
	return items, nil
}

func Delete(id int)(bool, error){
	db, err := getDb()
	defer db.Close()
	if err != nil {
		return false, err
	}

	sql := "delete from " + getTableName() + " where id = ? "
	db.MustExec(sql, id)
	return true, nil
}
