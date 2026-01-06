package model

type CreateCategory struct {
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}

type UpdateCategory struct {
	Id   int32  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}
