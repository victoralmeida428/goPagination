package pagination

type IPagination interface {
	GetPageSize() int
	SetPageSize(int)
	SetPage(int)
	SetTotalCount(int)
	SetRawQuery(string, ...interface{})
	SetCountByQuery() error
	GetQuery() string
	JSON(data []interface{}) (Envelope, error)
	SetOrder(...string)
}
