package postgres

type Condition struct {
	Name string
	Pagination
}
type Pagination struct {
	Limit  uint64
	Offset uint64
}
