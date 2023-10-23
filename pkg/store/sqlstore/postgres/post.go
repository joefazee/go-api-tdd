package postgres

import (
	"github.com/joefazee/go-api-tdd/pkg/domain"
)

var (
	sqlCreatePost = `INSERT INTO posts 
    	(user_id, title, body) 
			VALUES ($1, $2, $3) RETURNING id, user_id, title, body, created_at`

	sqlDeletePostByID = `DELETE FROM posts WHERE id = $1`

	sqlDeleteAllPosts = `DELETE FROM posts`

	sqlGetAllUserPosts = `SELECT id, user_id, title, body, created_at FROM posts WHERE user_id = $1`
)

func (q *postgresStore) CreatePost(post *domain.Post) (*domain.Post, error) {
	err := q.db.QueryRow(sqlCreatePost, post.UserID, post.Title, post.Body).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Body,
		&post.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}
func (q *postgresStore) GetAllUserPosts(userID int64) ([]domain.Post, error) {
	rows, err := q.db.Query(sqlGetAllUserPosts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post

	for rows.Next() {
		post := domain.Post{}
		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Body,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
func (q *postgresStore) DeleteAllPosts() error {
	_, err := q.db.Exec(sqlDeleteAllPosts)
	if err != nil {
		return err
	}

	return nil
}
func (q *postgresStore) DeletePostByID(ID int64) error {
	_, err := q.db.Exec(sqlDeletePostByID, ID)
	if err != nil {
		return err
	}

	return nil
}
