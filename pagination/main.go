package pagination

import "database/sql"

type Pagination[E any] struct {
	PaginationAbstract[E]
}

func New[E any](pageSize, pageNum int, db *sql.DB) *Pagination[E] {
	abstract := PaginationAbstract[E]{pageSize: pageSize, pageNum: pageNum, conn: db}
	return &Pagination[E]{abstract}
}
