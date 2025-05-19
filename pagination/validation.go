package pagination

import (
	"errors"
)

func (p PaginationAbstract[E]) validate() error {
	
	if p.pageNum < 1 {
		return errors.New("pagination pageNum cannot be less than 1")
	}
	if p.totalCount+p.pageSize <= p.pageSize*p.pageNum {
		
		return errors.New("pagination pageNum and totalCount must be greater than pageSize")
	}
	
	return nil
}
