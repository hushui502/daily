package dbops

func AddUserCredential(loginName string, pwd string) error {
	stmIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES(?, ?)")
	if err != nil {
		return err
	}

	_, err = stmIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}

	defer stmIns.Close()

	return nil
}

func Get() {
	stm, err := dbConn.Prepare()
}
