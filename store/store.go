package store

import (
  "database/sql"
	"github.com/AntJanus/budget-arcade-bot/config"
)

func SaltySave(gameName string, userName string) error {
  db, err := sql.Open("sqlite3", config.DBName)
  checkErr(err)

  stmt, err : = db.Prepare("INSERT INTO saltiness(game_name, user_name) values (?,?)")
  checkErr(err)

  res, err := stmt.Exec(gameName, userName)
  checkErr(err)

}
