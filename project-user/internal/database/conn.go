package database

type DbConn interface {
	Begin()
	RollBack()
	Commit()
}
