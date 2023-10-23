package postgres

import (
	"github.com/joefazee/go-api-tdd/pkg/domain"
	"testing"
)

func TestCreatePost(t *testing.T) {

	pStore := NewPostgresStore(testDB)

	user := &domain.User{
		Email:    "test@test.com",
		Password: "ess",
		Name:     "John Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	post := &domain.Post{
		UserID: createdUser.ID,
		Title:  "test",
		Body:   "test",
	}

	createdPost, err := pStore.CreatePost(post)
	if err != nil {
		t.Fatal(err)
	}

	if createdPost.ID == 0 {
		t.Errorf("want id not to be zero")
	}

	if post.Title != createdPost.Title {
		t.Errorf("expected %q; got %q", post.Title, createdPost.Title)
	}

	err = pStore.DeletePostByID(createdPost.ID)
	if err != nil {
		t.Fatal(err)
	}

	_ = pStore.DeleteUserByID(createdUser.ID)
}

func TestDeleteAllPosts(t *testing.T) {
	pStore := NewPostgresStore(testDB)

	user := &domain.User{
		Email:    "testss@test.com",
		Password: "ess",
		Name:     "John Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	post := &domain.Post{
		UserID: createdUser.ID,
		Title:  "sss test",
		Body:   "test",
	}

	createdPost, err := pStore.CreatePost(post)
	if err != nil {
		t.Fatal(err)
	}

	err = pStore.DeleteAllPosts()
	if err != nil {
		t.Fatal(err)
	}

	ps, err := pStore.GetAllUserPosts(createdPost.ID)
	if err != nil {
		t.Errorf("expect no error")
	}

	if len(ps) != 0 {
		t.Errorf("expected no posts for this user after a call to DeleteAllPosts")
	}

	_ = pStore.DeleteUserByID(createdUser.ID)
}

func TestGetAllUserPosts(t *testing.T) {

	pStore := NewPostgresStore(testDB)

	user := &domain.User{
		Email:    "testss@test.com",
		Password: "ess",
		Name:     "John Doe",
	}

	createdUser, err := pStore.CreateUser(user)
	if err != nil {
		t.Fatal(err)
	}

	_, err = pStore.CreatePost(&domain.Post{
		UserID: createdUser.ID,
		Title:  "sss sss test",
		Body:   "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = pStore.CreatePost(&domain.Post{
		UserID: createdUser.ID,
		Title:  "sss sss sstest",
		Body:   "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	ps, err := pStore.GetAllUserPosts(createdUser.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(ps) != 2 {
		t.Errorf("want user posts to be 2; got %d", len(ps))
	}

	_ = pStore.DeleteAllPosts()
}
