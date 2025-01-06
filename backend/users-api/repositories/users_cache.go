package users

import (
	"fmt"
	"time"
	userDAO "users-api/dao"

	"github.com/karlseguin/ccache"
	_ "github.com/karlseguin/ccache"
)

const (
	keyFormat = "user:%s"
)

type CacheConfig struct {
	MaxSize      int64
	ItemsToPrune uint32
	Duration     time.Duration
}

type Cache struct {
	client   *ccache.Cache
	duration time.Duration
}

func NewCache(config CacheConfig) Cache {
	client := ccache.New(ccache.Configure().
		MaxSize(config.MaxSize).
		ItemsToPrune(config.ItemsToPrune))
	return Cache{
		client:   client,
		duration: config.Duration,
	}
}

func (repo Cache) InsertUser(user userDAO.User) (userDAO.User, error) {
	key := fmt.Sprintf(keyFormat, user.Id)
	repo.client.Set(key, user, repo.duration)
	return userDAO.User{}, nil
}

func (repo Cache) GetUserById(id int64) (userDAO.User, error) {
	key := fmt.Sprintf(keyFormat, id)
	item := repo.client.Get(key)
	if item == nil {
		return userDAO.User{}, fmt.Errorf("not found item with key %s", key)
	}
	if item.Expired() {
		return userDAO.User{}, fmt.Errorf("item with key %s is expired", key)
	}
	userDao, ok := item.Value().(userDAO.User)
	if !ok {
		return userDAO.User{}, fmt.Errorf("error converting item with key %s", key)
	}
	return userDao, nil
}

func (repo Cache) GetUserByName(user userDAO.User) (userDAO.User, error) {
	key := fmt.Sprintf(keyFormat, user.User)
	item := repo.client.Get(key)
	if item == nil {
		return userDAO.User{}, fmt.Errorf("not found item with key %s", key)
	}
	if item.Expired() {
		return userDAO.User{}, fmt.Errorf("item with key %s is expired", key)
	}
	userDao, ok := item.Value().(userDAO.User)
	if !ok {
		return userDAO.User{}, fmt.Errorf("error converting item with key %s", key)
	}
	return userDao, nil
}
