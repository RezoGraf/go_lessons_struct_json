package lib_db

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"lessons/models"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func DBConn(dbString string) (db *sqlx.DB) {
	if _, err := os.Stat(dbString); errors.Is(err, fs.ErrNotExist) {
		os.Create(dbString)
	}
	db, err := sqlx.Open("sqlite3", dbString)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	return db
}

func CreateTables(db *sqlx.DB) error {
	createaAditionalinfoTable := `
	CREATE TABLE product (id INTEGER PRIMARY KEY AUTOINCREMENT, title CHAR (50) UNIQUE NOT NULL, tags INTEGER UNIQUE NOT NULL, description VARCHAR (250), price REAL, additionalinfo INTEGER UNIQUE);
	CREATE TABLE additionalinfo ("key" INTEGER, title VARCHAR (250), comment VARCHAR (250));
	CREATE TABLE tags ("key" INTEGER, id INTEGER, tag VARCHAR (250));
	`
	_, err := db.Exec(createaAditionalinfoTable)
	return err
}

func DBInit(dbString string) (db *sqlx.DB, err error) {
	db = DBConn(dbString)
	err = CreateTables(db)
	return
}

func GetIDFromDB(db *sqlx.DB) (int, error) {
	var p int
	err := db.Get(&p, "SELECT count(id) FROM product;")
	p = p + 1
	return p, err
}

func IfExistsTitleFromDB(db *sqlx.DB, title string) (int, error) {
	var p int
	err := db.Get(&p, "SELECT count(id) FROM product where title=?;", title)
	return p, err
}

func InsertProductDB(db *sqlx.DB, k *models.Product) error {

	p, err := GetIDFromDB(db)
	if err != nil {
		return err
	}
	insertProduct := ""
	insertProduct += fmt.Sprintf("INSERT INTO additionalinfo (\"key\", title, comment) VALUES (%d, '%s', '%s');\n", p, k.Additionalinfo.Title, k.Additionalinfo.Comment)
	for _, v := range k.Tags {
		insertProduct += fmt.Sprintf("INSERT INTO tags (\"key\", id, tag) VALUES (%d, %d, '%s');\n", p, v.ID, v.Tag)
	}
	insertProduct += fmt.Sprintf("INSERT INTO product (id, title, tags, description, price, additionalinfo) VALUES (%[1]d, '%[2]s', %[1]d, '%[3]s', '%[4]s', %[1]d);\n", p, *k.Title, k.Description, k.Price)
	tx := db.MustBegin()
	tx.MustExec(insertProduct)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func DeleteProductDB(db *sqlx.DB, title string) error {
	var p int
	err := db.Get(&p, "SELECT tags FROM product where title=?;", title)
	deleteProduct := fmt.Sprintf("DELETE from additionalinfo where key=%d;", p)
	deleteProduct += fmt.Sprintf("DELETE from tags where key=%d;", p)
	deleteProduct += fmt.Sprintf("delete from product where id=%d;", p)
	tx := db.MustBegin()
	tx.MustExec(deleteProduct)
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func CheckTitleExists(db *sqlx.DB, k *models.Product) bool {
	existsQuery, err := IfExistsTitleFromDB(db, *k.Title)
	if err != nil {
		log.Println(err)
		return false
	}
	if existsQuery > 0 {
		return true
	}
	return false
}

func ListProducts(db *sqlx.DB) ([]models.Product, error) {
	var pp []models.Product
	listProductQuery := "SELECT title,description,price from product;"
	rows, err := db.Queryx(listProductQuery)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	for rows.Next() {
		var p models.Product
		// var t models.Tag
		var tt []models.Tag
		err := rows.StructScan(&p)
		if err != nil {
			fmt.Println(err)
			panic(err.Error())
		}
		var id int
		queryID := fmt.Sprintf("SELECT tags FROM product where title=\"%s\";", *p.Title)
		err = db.Get(&id, queryID)
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		var ai models.Additionalinfo
		queryAI := fmt.Sprintf("SELECT title, comment from additionalinfo where key=%d;", id)
		err = db.Get(&ai, queryAI)
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		p.Additionalinfo.Title = ai.Title
		p.Additionalinfo.Comment = ai.Comment

		selectTagsQuery := fmt.Sprintf("SELECT id, tag from tags where key=%d;", id)
		err = db.Select(&tt, selectTagsQuery)
		p.Tags = tt
		pp = append(pp, p)
	}
	return pp, err
}

