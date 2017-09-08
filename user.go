package main

import (
	"math/rand"
	"time"

	"sync"

	"github.com/pkg/errors"
	"github.com/upamune/grpc-error-sample/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
)

var (
	ErrNoSuchUser    = errors.New("no such user")
	ErrDuplicateUser = errors.New("duplicate user")
)

type UserService struct {
	userMap   map[string]api.User
	userMapMu sync.RWMutex
	r         *rand.Rand
}

func NewUserService() *UserService {
	now := time.Now()

	return &UserService{
		userMap: make(map[string]api.User),
		r:       rand.New(rand.NewSource(now.UnixNano())),
	}
}
func (s *UserService) getUser(id string) (api.User, error) {
	s.userMapMu.RLock()
	user, ok := s.userMap[id]
	s.userMapMu.RUnlock()
	if !ok {
		return api.User{}, ErrNoSuchUser
	}

	return user, nil
}

func (s *UserService) saveUser(user api.User) error {
	if _, err := s.getUser(user.Id); err == nil {
		return ErrDuplicateUser
	} else {
		if err != ErrNoSuchUser {
			return err
		}
	}

	s.userMapMu.Lock()
	s.userMap[user.Id] = user
	s.userMapMu.Unlock()

	return nil
}

func (s *UserService) Get(ctx context.Context, req *api.GetUserRequest) (*api.GetUserResponse, error) {
	user, err := s.getUser(req.UserId)
	if err != nil {
		if err == ErrNoSuchUser {
			return nil, api.Errorf(codes.InvalidArgument, api.Error_NO_SUCH_USER, false, "no such user id: %s", req.UserId)
		}
		return nil, api.Errorf(codes.Internal, api.Error_UNKNOWN, false, "unknown error", req.UserId)
	}

	return &api.GetUserResponse{
		User: &user,
	}, nil
}
func (s *UserService) Create(ctx context.Context, req *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	user := api.User{
		Id:   req.Id,
		Name: req.Name,
	}
	if err := s.saveUser(user); err != nil {
		if err == ErrDuplicateUser {
			return nil, api.Errorf(codes.InvalidArgument, api.Error_DUPLICATE_USER_ID, false, "duplicate user id: %s", req.Id)
		}
		return nil, api.Errorf(codes.Internal, api.Error_UNKNOWN, false, "unknown error")
	}

	return &api.CreateUserResponse{}, nil
}

func (s *UserService) Delete(ctx context.Context, req *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	return nil, api.Errorf(codes.Internal, api.Error_DB_CONNECTION, false, "failed to connect to db")
}
