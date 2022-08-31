package inter

type IDBServer interface {
	Exec(sql string) error
}
