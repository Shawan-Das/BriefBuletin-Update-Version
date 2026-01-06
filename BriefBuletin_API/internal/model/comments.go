package model

type CreateComment struct {
	ArticleID int32  `db:"article_id" json:"article_id"`
	UserName  string `db:"user_name" json:"user_name"`
	UserEmail string `db:"user_email" json:"user_email"`
	Content   string `db:"content" json:"content"`
}
