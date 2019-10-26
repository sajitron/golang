//* showing one-to-many & many-to-one relationships with Go

package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type Post struct {
	Id			int
	Content	string
	Author	string
	Comments	[]Comment
}

type Comment struct {
	Id			int
	Content	string
	Author	string
	Post		*Post
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
}

//* create a single comment
func (comment *Comment) Create() (err error) {
	if comment.Post == nil {
		err = errors.New("Post not found")
		return
	}
	err = Db.QueryRow("insert into comments (content, author, post_id) values ($1, $2, $3) returning id", comment.Content, comment.Author, comment.Post.Id).Scan(&comment.Id)
	return
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	post.Comments = []Comment{}
	err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)

	rows, err := Db.Query("select id, content, author from comments")
	if err != nil {
		return
	}
	for rows.Next() {
		comment := Comment{Post: &post}
		err = rows.Scan(&comment.Id, &comment.Content, &comment.Author)
		if err != nil {
			return
		}
		post.Comments = append(post.Comments, comment)
	}
	rows.Close()
	return
}

func (post *Post) Create() (err error) {
	err = Db.QueryRow("insert into posts (content, author) values ($1, $2) returning id", post.Content, post.Author).Scan(&post.Id)
	return
}

func main() {
	post := Post{Content: "You killed my father!", Author: "Inigo Montoya"}
	post.Create()

	comment1 := Comment{Content: "Now you have to pay!", Author: "Kevin Pearson", Post: &post}
	comment1.Create()

	comment2 := Comment{Content: "Great Movie!", Author: "Mandy Patinkin", Post: &post}
	comment2.Create()

	readPost, _ := GetPost(post.Id)

	fmt.Println(readPost)
	fmt.Println(readPost.Comments)
	fmt.Println(readPost.Comments[0].Post)
}