package pagination

import (
	"errors"
	"strings"
)

var _ Pager = (*Pagination)(nil)

// New .
func New(options ...Option) *Pagination {
	opts := Apply(options...)

	remain := opts.total % opts.pageSize
	totalPages := opts.total / opts.pageSize
	if remain > 0 {
		totalPages++
	}

	hasNext := opts.total-opts.pageNumber-opts.pageSize > 0

	return &Pagination{
		pageNumber: opts.pageNumber,
		pageSize:   opts.pageSize,
		sort:       getSort(opts.sort),
		total:      opts.total,
		data:       opts.data,
		totalPages: totalPages,
		hasNext:    hasNext,
		keyword:    opts.keyword,
	}
}

func (p *Pagination) Limit() int {
	if p.pageSize < 1 {
		return 10
	}
	return p.pageSize
}
func (p *Pagination) Offset() int {
	if p.pageNumber == 0 {
		p.pageNumber = 1
	}
	if p.pageSize < 1 {
		p.pageSize = 10
	}
	offset := (p.pageNumber - 1) * p.pageSize
	return offset
}

func (p *Pagination) Sort() string {
	return p.sort
}

type OrderBy struct {
	OrderKey string `json:"order_key"`
	Sort     string `json:"sort"`
}

func (p *Pagination) PageNumber() int {
	return p.pageNumber
}

func (p *Pagination) PageSize() int {
	return p.pageSize
}

func (p *Pagination) TotalPages() int {
	return p.totalPages
}

func (p *Pagination) Total() int {
	return p.total
}

func (p *Pagination) Data() []interface{} {
	return p.data
}

func (p *Pagination) DataSize() int {
	return len(p.data)
}

func (p *Pagination) HasNext() bool {
	return p.hasNext
}

func (p *Pagination) HasData() bool {
	return p.DataSize() > 0
}

func (p *Pagination) Valid() error {
	if p.pageNumber == 0 {
		p.pageNumber = 1
	}
	if p.pageSize == 0 {
		p.pageSize = 10
	}

	if p.pageNumber < 0 {
		return errors.New("current MUST be larger than 0")
	}

	if p.pageSize < 0 {
		return errors.New("invalid pageSize")
	}
	return nil
}

// convert to mysql sort, each column name preceded by a '-' sign, indicating descending order, otherwise ascending order, example:
//
//	columnNames="name" means sort by name in ascending order,
//	columnNames="-name" means sort by name descending,
//	columnNames="name,age" means sort by name in ascending order, otherwise sort by age in ascending order,
//	columnNames="-name,-age" means sort by name descending before sorting by age descending.
func getSort(columnNames string) string {
	columnNames = strings.Replace(columnNames, " ", "", -1)
	if columnNames == "" {
		return "id DESC"
	}

	names := strings.Split(columnNames, ",")
	strs := make([]string, 0, len(names))
	for _, name := range names {
		if name[0] == '-' && len(name) > 1 {
			strs = append(strs, name[1:]+" DESC")
		} else {
			strs = append(strs, name+" ASC")
		}
	}

	return strings.Join(strs, ", ")
}
