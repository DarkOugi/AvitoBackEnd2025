package main

import (
	initDB "avito/initDB"
	pg "avito/pkg/db"
	"fmt"
)

// для инита базы данных
func main() {
	db, err := pg.Connect("localhost", "5432", "avito", "0000", "avitodb")
	if err != nil {
		fmt.Printf("Connect to db %s\n", err.Error())
	}
	err = initDB.CreateTables(db)
	if err != nil {
		fmt.Printf("Error create table  %s\n", err.Error())
	}
	err = initDB.CreateIndex(db)
	if err != nil {
		fmt.Printf("Error create index  %s\n", err.Error())
	}
	err = initDB.InsertValue(db)
	if err != nil {
		fmt.Printf("Error insert to table  %s\n", err.Error())
	}

}

//func main() {
//	// create, insert := initdb.Init_mech()
//	// db, err := pg.Connect("localhost", "5432", "avito", "0000", "avitodb")
//	// if err != nil {
//	// 	fmt.Printf("Connect to db %s\n", err.Error())
//	// } else {
//	// 	_, err = db.Exec(create)
//	// 	if err != nil {
//	// 		fmt.Printf("Create table err: %s\n", err.Error())
//	// 	} else {
//	// 		fmt.Printf("Create table success\n")
//	// 	}
//	// 	in, err := db.Exec(insert)
//	// 	if err != nil {
//	// 		fmt.Printf("Create rows err: %s\n", err.Error())
//	// 		panic(err)
//	// 	} else {
//	// 		fmt.Printf("Insert table %d rows success\n", in)
//	// 	}
//	// 	defer db.Close()
//	// }
//	// q := token.MetaJWT("Denis.Zhilin")
//	// fmt.Println(q)
//	fmt.Println(pg.HashPassword("asdadqwq212e1d2wd"))
//
//}
