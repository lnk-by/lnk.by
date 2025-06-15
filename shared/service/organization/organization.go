package organization

type Organization struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
