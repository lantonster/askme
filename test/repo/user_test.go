package repo

import (
	"context"
	"strings"
	"testing"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/lantonster/askme/pkg/errors"
	"github.com/lantonster/askme/pkg/reason"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	repo := repo.NewUserRepo(data)

	user := &model.User{Id: 1, CreatedAt: 1, UpdatedAt: 1, Username: "user001"}
	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	var users []*model.User
	db.Model(&model.User{}).Find(&users)
	assert.Len(t, users, 1)
	assert.Equal(t, user, users[0])
}

func TestGetUserByEmail(t *testing.T) {
	repo := repo.NewUserRepo(data)

	t.Run("exist", func(t *testing.T) {
		user := &model.User{Id: 2, CreatedAt: 1, UpdatedAt: 1, Username: "user002", Email: "email002"}

		err := repo.CreateUser(context.Background(), user)
		assert.NoError(t, err)

		u, exist, err := repo.GetUserByEmail(context.Background(), user.Email)
		assert.NoError(t, err)
		assert.True(t, exist)
		assert.Equal(t, user, u)
	})

	t.Run("no exist", func(t *testing.T) {
		_, exist, err := repo.GetUserByEmail(context.Background(), "no exist")
		assert.NoError(t, err)
		assert.False(t, exist)
	})
}

func TestGetUserByUsername(t *testing.T) {
	repo := repo.NewUserRepo(data)

	t.Run("exist", func(t *testing.T) {
		user := &model.User{Id: 3, CreatedAt: 1, UpdatedAt: 1, Username: "user003", Email: "email003"}
		err := repo.CreateUser(context.Background(), user)
		assert.NoError(t, err)

		u, exist, err := repo.GetUserByUsername(context.Background(), user.Username)
		assert.NoError(t, err)
		assert.True(t, exist)
		assert.Equal(t, user, u)
	})

	t.Run("no exist", func(t *testing.T) {
		_, exist, err := repo.GetUserByUsername(context.Background(), "no exist")
		assert.NoError(t, err)
		assert.False(t, exist)
	})
}

func TestGenerateUniqueUsername(t *testing.T) {
	repo := repo.NewUserRepo(data)

	t.Run("pinyin", func(t *testing.T) {
		username := "蓝"
		_username, err := repo.GenerateUniqueUsername(context.Background(), username)
		assert.NoError(t, err)
		assert.Equal(t, "lan", _username)
	})

	t.Run("too long", func(t *testing.T) {
		username := strings.Repeat("a", 31)
		_, err := repo.GenerateUniqueUsername(context.Background(), username)
		assert.Error(t, err)

		e, ok := err.(*errors.Error)
		assert.True(t, ok)
		assert.Equal(t, reason.UsernameInvalid, e.Reason)
	})

	t.Run("reserved", func(t *testing.T) {
		username := "admin"
		_, err := repo.GenerateUniqueUsername(context.Background(), username)
		assert.Error(t, err)

		e, ok := err.(*errors.Error)
		assert.True(t, ok)
		assert.Equal(t, reason.UsernameInvalid, e.Reason)
	})

	t.Run("repeat", func(t *testing.T) {
		repo.CreateUser(context.Background(), &model.User{Id: 4, Username: "user004"})

		username := "user004"
		_username, err := repo.GenerateUniqueUsername(context.Background(), username)
		assert.NoError(t, err)
		assert.Greater(t, len(_username), len(username))
		assert.Equal(t, username, _username[:len(username)])
	})
}

func TestIncrRank(t *testing.T) {
	repo := repo.NewUserRepo(data)

	user := &model.User{Id: 5, CreatedAt: 1, UpdatedAt: 1, Username: "user005", Email: "email005"}
	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	err = repo.IncrRank(context.Background(), user.Id, 0, 10)
	assert.NoError(t, err)

	u, exist, err := repo.GetUserByUsername(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, int64(10), u.Rank)
}

func TestUpdateEmailStatus(t *testing.T) {
	repo := repo.NewUserRepo(data)

	user := &model.User{Id: 6, CreatedAt: 1, UpdatedAt: 1, Username: "user006", Email: "email006"}
	err := repo.CreateUser(context.Background(), user)
	assert.NoError(t, err)

	err = repo.UpdateEmailStatus(context.Background(), user.Id, "verified")
	assert.NoError(t, err)

	u, exist, err := repo.GetUserByUsername(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.True(t, exist)
	assert.Equal(t, "verified", u.MailStatus)
}
