package postgres

type Condition struct {
	Name  string
	Epoch []uint64
	Pagination
}
type Pagination struct {
	Limit  uint64
	Offset uint64
}
type Aggregate int8

const (
	Day = Aggregate(iota)
	Month
	Week
	Year
)

func SearchAggregate(name string) Aggregate {
	switch name {
	case "month":
		return Month
	case "week":
		return Week
	case "year":
		return Year
	default:
		return Day
	}
}
