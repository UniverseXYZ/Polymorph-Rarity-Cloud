package structs

type QueryParams struct {
	Take      string `schema:"take"`
	Page      string `schema:"page"`
	SortField string `schema:"sortField"`
	SortDir   string `schema:"sortDir"`
	Select    string `schema:"select"`
	Filter    string `schema:"filter"`
	Search    string `schema:"search"`
}
