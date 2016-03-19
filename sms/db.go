package main

import (
	"log"
)

func addToDB(user_id int, phone, message string) int64 {
	stmt, err := db.Prepare("INSERT messages SET user_id=?, phone=?,message=?,date_create=NOW()")
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(user_id, phone, message)
	if err != nil {
		log.Println(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
	}
	log.Println(id)
	return id
}

func updateToDb(id int64) {
	// update
	stmt, err := db.Prepare("UPDATE messages set date_send=NOW() where id=?")
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Println(err)
	}
	log.Print(res)
}
