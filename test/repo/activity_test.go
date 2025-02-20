package repo

import (
	"context"
	"testing"

	"github.com/lantonster/askme/internal/model"
	"github.com/lantonster/askme/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestCreateActivity(t *testing.T) {
	repo := repo.NewActivityRepo(data)

	activity := &model.Activity{Id: 1}
	err := repo.CreateActivity(context.Background(), activity)
	assert.NoError(t, err)

	var activities []*model.Activity
	db.Model(&model.Activity{}).Find(&activities)
	assert.Len(t, activities, 1)
	assert.Equal(t, activity, activities[0])
}

func TestGetActivityByType(t *testing.T) {
	repo := repo.NewActivityRepo(data)

	t.Run("exist", func(t *testing.T) {
		activity := &model.Activity{Id: 2, Type: "test"}
		db.Create(activity)

		_activity, ok, err := repo.GetActivityByType(context.Background(), activity.Type)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, activity, _activity)
	})

	t.Run("not exist", func(t *testing.T) {
		_, ok, err := repo.GetActivityByType(context.Background(), "not_exist")
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}
