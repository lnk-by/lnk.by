package campaign

type Campaign struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
