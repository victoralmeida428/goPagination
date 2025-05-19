package pagination

import (
	"database/sql"
	"fmt"
)

type PaginationAbstract[E any] struct {
	pageSize   int
	pageNum    int
	conn       *sql.DB
	totalCount int
	orderBy    string
	rawQuery   string
	data       []E
	params     []interface{}
}

func (p PaginationAbstract[E]) GetPageSize() int {
	return p.pageSize
}

func (p *PaginationAbstract[E]) SetRawQuery(query string, args ...interface{}) {
	p.rawQuery = query
	p.params = args
}

func (p *PaginationAbstract[E]) SetCountByQuery() error {
	
	countQuery := fmt.Sprintf("select count(*) from (%s)", p.rawQuery)
	
	return p.conn.QueryRow(countQuery, p.params...).Scan(&p.totalCount)
}

func (p *PaginationAbstract[E]) SetTotalCount(i int) {
	
	p.totalCount = i
}

func (p PaginationAbstract[E]) GetQuery() string {
	limit := fmt.Sprintf(" limit %d offset %d", p.pageSize, p.calculateOffset())
	return p.rawQuery + p.orderBy + limit
}

func (p PaginationAbstract[E]) nextPage() *int {
	if p.pageNum*p.pageSize >= p.totalCount {
		return nil
	}
	page := p.pageNum + 1
	return &page
}

func (p *PaginationAbstract[E]) runSQL(scan func(e *[]E, rows *sql.Rows) error) error {
	query := p.GetQuery()
	
	rows, err := p.conn.Query(query, p.params...)
	if err != nil {
		return err
	}
	
	defer rows.Close()
	
	results := make([]E, 0)
	err = scan(&results, rows)
	
	p.data = results
	
	return err
	
}

func (p PaginationAbstract[E]) previousPage() *int {
	page := p.pageNum - 1
	if page < 1 {
		return nil
	}
	return &page
}

func (p PaginationAbstract[E]) JSON(scan func(e *[]E, rows *sql.Rows) error) (Envelope, error) {
	
	if err := p.SetCountByQuery(); err != nil {
		return nil, err
	}
	if err := p.validate(); err != nil {
		return nil, err
	}
	if err := p.runSQL(scan); err != nil {
		return nil, err
	}
	
	return Envelope{
		"data":          p.data,
		"next_page":     p.nextPage(),
		"count":         p.totalCount,
		"previous_page": p.previousPage(),
	}, nil
}

func (p *PaginationAbstract[E]) SetOrder(orders ...string) {
	var order string
	if len(orders) > 0 {
		order = " order by "
	}
	for _, i := range orders {
		order += i + ","
	}
	if len(order) > 0 {
		order = order[:len(order)-1]
	}
	p.orderBy = order
}

func (p PaginationAbstract[E]) calculateOffset() int {
	return (p.pageNum - 1) * p.pageSize
}

func (p *PaginationAbstract[E]) SetPageSize(pageSize int) {
	p.pageSize = pageSize
}

func (p *PaginationAbstract[E]) SetPage(page int) {
	p.pageNum = page
}
