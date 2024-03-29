package pkg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"pbmysql-go/dbproto"
	"testing"
)

func GetMysqlConfig() *mysql.Config {
	file, err := os.Open("db.json")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	decoder := json.NewDecoder(file)
	jsonConfig := JsonConfig{}
	err = decoder.Decode(&jsonConfig)
	if err != nil {
		log.Fatal(err)
	}
	return NewMysqlConfig(jsonConfig)
}

func TestCreateTable(t *testing.T) {
	pbMySqlDB := NewPb2DbTables()
	pbMySqlDB.AddMysqlTable(&dbproto.GolangTest{})

	mysqlConfig := GetMysqlConfig()
	conn, err := mysql.NewConnector(mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()
	pbMySqlDB.SetDB(db, mysqlConfig.DBName)
	pbMySqlDB.UseDB()
	result, err := db.Exec(pbMySqlDB.GetCreateTableSql(&dbproto.GolangTest{}))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}

func TestAlterTable(t *testing.T) {
	pbMySqlDB := NewPb2DbTables()
	pbMySqlDB.AddMysqlTable(&dbproto.GolangTest{})

	mysqlConfig := GetMysqlConfig()
	conn, err := mysql.NewConnector(mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()

	pbMySqlDB.SetDB(db, mysqlConfig.DBName)
	pbMySqlDB.UseDB()

	pbMySqlDB.AlterTableAddField(&dbproto.GolangTest{})
}

func TestLoadSave(t *testing.T) {
	pbMySqlDB := NewPb2DbTables()
	pbsave := &dbproto.GolangTest{
		Id:      1,
		GroupId: 1,
		Ip:      "127.0.0.1",
		Port:    3306,
		Player: &dbproto.Player{
			PlayerId: 111,
			Name:     "foo\\0bar,foo\\nbar,foo\\rbar,foo\\Zbar,foo\\\"bar,foo\\\\bar,foo\\'bar",
		},
	}
	pbMySqlDB.AddMysqlTable(pbsave)
	mysqlConfig := GetMysqlConfig()
	conn, err := mysql.NewConnector(mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()
	pbMySqlDB.SetDB(db, mysqlConfig.DBName)
	pbMySqlDB.UseDB()

	pbMySqlDB.Save(pbsave)

	pbload := &dbproto.GolangTest{}
	pbMySqlDB.LoadOneByKV(pbload, "id", "1")
	if !proto.Equal(pbsave, pbload) {
		log.Fatal("pb not equal")
	}
}

func TestLoadSaveList(t *testing.T) {
	pbMySqlDB := NewPb2DbTables()
	pbsavelist := &dbproto.GolangTestList{
		TestList: []*dbproto.GolangTest{
			{
				Id:      1,
				GroupId: 1,
				Ip:      "127.0.0.1",
				Port:    3306,
				Player: &dbproto.Player{
					PlayerId: 111,
					Name:     "foo\\0bar,foo\\nbar,foo\\rbar,foo\\Zbar,foo\\\"bar,foo\\\\bar,foo\\'bar",
				},
			},
			{
				Id:      2,
				GroupId: 1,
				Ip:      "127.0.0.1",
				Port:    3306,
				Player: &dbproto.Player{
					PlayerId: 111,
					Name:     "foo\\0bar,foo\\nbar,foo\\rbar,foo\\Zbar,foo\\\"bar,foo\\\\bar,foo\\'bar",
				},
			},
		},
	}
	pbMySqlDB.AddMysqlTable(&dbproto.GolangTest{})
	mysqlConfig := GetMysqlConfig()
	conn, err := mysql.NewConnector(mysqlConfig)
	if err != nil {
		log.Fatal(err)
	}
	db := sql.OpenDB(conn)
	defer db.Close()
	pbMySqlDB.SetDB(db, mysqlConfig.DBName)
	pbMySqlDB.UseDB()

	pbloadlist := &dbproto.GolangTestList{}
	pbMySqlDB.LoadList(pbloadlist)
	if !proto.Equal(pbsavelist, pbloadlist) {
		fmt.Println(pbsavelist.String())
		fmt.Println(pbloadlist.String())
		log.Fatal("pb not equal")
	}
}
