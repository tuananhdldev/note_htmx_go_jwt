package database

import (
	"database/sql"
	"fmt"
	"loa/user_content/types"
	"loa/user_content/utils"

	_ "github.com/lib/pq"
)


type DBInstance struct {
	db *sql.DB
}
var DB DBInstance
func NewStoreInstance() error {
    env:= utils.GetEnv()
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
        env.User, env.Password, env.URL, env.DB, env.PORT)
	fmt.Println(connStr)
	db, err := sql.Open("postgres",connStr)	
	if err != nil {
		fmt.Println("error connecting to database"+err.Error())
		return err
	}
	fmt.Println("connection succeeded")
	DB = DBInstance{
		db: db,
	}
	return nil
}
func (db *DBInstance) InsertDoc(req *types.CreateDocumentRequest) error {
      	query := `INSERT INTO content (hash, author, title, content)
		Values($1,$2,$3,$4)`
		_, err := db.db.Exec(query, req.Hash, req.Author, req.Title, req.Content)
        return err
}
func (db *DBInstance) UpdateDoc(req *types.UpdateDocumentRequest) error{
	query := `UPDATE content SET title = $1, content = $2, update_at=NOW() WHERE hash = $3`
	_,err := db.db.Exec(query, req.Title, req.Content, req.Hash)
	return err
}
//delete Document
func (db *DBInstance) DeleteDoc(hash string) error {
	query := `Delete From content Where hash=$1`
	_,err := db.db.Exec(query,hash)
	return err
}

func (db *DBInstance) GetByHash(hash string) (*types.Document, error){

	query := `SELECT * FROM content WHERE hash=$1`
	var document types.Document
	rows := db.db.QueryRow(query,hash)
	err:= rows.Scan(&document.Hash,&document.Author,&document.Title,&document.Content,&document.CreatedAt,&document.UpdatedAt)
	if  err != nil {
		return &document,err
	}
	return &document, err
}
func (db *DBInstance) GetAll() ([]types.Document, error){
	query := `SELECT * from content`
	var documents []types.Document
	rows,err := db.db.Query(query)
	if err !=nil{
		return documents, err
	}
	for rows.Next(){
		var document types.Document
		err := rows.Scan(&document.Hash,&document.Author, &document.Title, &document.Content,&document.CreatedAt)
		if err !=nil{
			return documents,err
		}
       documents = append(documents, document)   	
	}
   return documents,nil
}
//USER SET
func (db *DBInstance) CreateUser(req types.CreateUserRequest) error{
	query := `INSERT INTO users (username,password) VALUES($1,$2)`
	hashedpassword, err := utils.HashPassword(req.Password)
	if err !=nil{
		return err
	}
	_, err = db.db.Exec(query, req.Username, hashedpassword)
	return err
}
//get user
func (db *DBInstance) GetUser(username string) (*types.User, error){
    query := `Select * from users where username=$1`
	var user types.User
	rows := db.db.QueryRow(query, username)
	if rows == nil {
		return &user, nil
	}
	err:= rows.Scan(&user.Username, &user.Password,&user.CreatedAt)
	if err != nil{
			
		return &user, err
	}
	
	return &user, nil
}
//get user docs
func (db *DBInstance) GetUserDocs(username string)([]types.Document, error){
	query := `
	SELECT * from content WHERE author=$1
	`

	var documents []types.Document

	rows, err := db.db.Query(query, username)
	if err != nil {
		return documents, err
	}

	for rows.Next() {
		var document types.Document
		err := rows.Scan(&document.Hash, &document.Author, &document.Title, &document.Content, &document.CreatedAt, &document.UpdatedAt)
		if err != nil {
			return documents, err
		}

		documents = append(documents, document)
	}

	return documents, nil
}


func (pq *DBInstance) DeleteUser(username string) error {
	query := `
	DELETE FROM users WHERE username=$1
	`
	_, err := pq.db.Exec(query, username)
	return err
}

func (pq *DBInstance) CreateTable() error {
	user_query := `CREATE TABLE IF NOT EXISTS users (
		username Text Primary key,
		password Text not null,
		created_at TIMESTAMP Default Now()
		)`
	if _,err :=pq.db.Exec(user_query);	
	 err != nil {
		return err
	 }
    content_query := `CREATE TABLE IF NOT EXISTS content(
         hash TEXT PRIMARY KEY,
		 author TEXT References users(username),
		 title TEXT NOT NULL Default 'Untitled',
		 content TEXT,
		 created_at TIMESTAMP Default NOW(),
		 updated_at TIMESTAMP Default NOW()

	)`
  if _ , err := pq.db.Exec(content_query); err != nil {
      return err
  }

	return nil
}