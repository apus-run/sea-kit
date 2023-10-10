package pagination

type Pager interface {

	// PageNumber will return the page number
	PageNumber() int

	// PageSize will return the page size
	PageSize() int

	// TotalPages will return the number of total pages
	TotalPages() int

	Total() int

	// Data will return the data
	Data() []any

	// DataSize will return the size of data.
	// Usually it's len(Data())
	DataSize() int

	// HasNext will return whether has next page
	HasNext() bool

	// HasData will return whether this page has data.
	HasData() bool

	Limit() int
	Offset() int

	// Sort get sort field
	Sort() string
}

// Pagination is the default implementation of Page interface
type Pagination struct {
	// pageNumber 当前页
	pageNumber int `json:"page_number,optional,default=1" form:"page_number,optional,default=1"`
	// pageSize 分页数
	pageSize int `json:"page_size,optional,default=10" form:"page_size,optional,default=10"`
	// sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
	sort string `json:"sort"`
	// total means total page count
	total int `json:"total"`
	// data 数据
	data []any `json:"data"`
	// totalPages 总页数
	totalPages int  `json:"total_pages,omitempty"`
	hasNext    bool `json:"has_next,omitempty"`

	keyword string `json:"keyword,optional" form:"keyword,optional"`
}

// Option is config option.
type Option func(*Pagination)

// DefaultOptions .
func DefaultOptions() *Pagination {
	return &Pagination{
		pageNumber: 1,
		pageSize:   10,
	}
}

func Apply(opts ...Option) *Pagination {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// WithPageSize .
func WithPageSize(pageSize int) Option {
	return func(o *Pagination) {
		o.pageSize = pageSize
	}
}

// WithPageNumber .
func WithPageNumber(pageNumber int) Option {
	return func(o *Pagination) {
		o.pageNumber = pageNumber
	}
}

// WithSort .
func WithSort(sort string) Option {
	return func(o *Pagination) {
		o.sort = sort
	}
}

// WithTotal .
func WithTotal(total int) Option {
	return func(o *Pagination) {
		o.total = total
	}
}

// WithData .
func WithData(data []any) Option {
	return func(o *Pagination) {
		o.data = data
	}
}
