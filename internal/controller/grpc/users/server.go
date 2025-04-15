package users

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	userv1 "github.com/kurochkinivan/user_proto/gen/go/users"
	"github.com/kurochkinivan/user_service/internal/entity"
	"google.golang.org/grpc"
)

type User interface {
	UpdateProfile(ctx context.Context, userID string, user *entity.User) (*entity.User, error)
}

type serverAPI struct {
	userv1.UnimplementedUserServer
	validate *validator.Validate
	user     User
}

func Register(gRPC *grpc.Server, validate *validator.Validate, user User) {
	userv1.RegisterUserServer(gRPC, &serverAPI{
		validate: validate,
		user:     user,
	})
}

func (s *serverAPI) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*userv1.UpdateProfileResponse, error) {
	if err := validateUserID(req.GetUserId()); err != nil {
		return nil, err
	}

	user, err := s.user.UpdateProfile(ctx, req.GetUserId(), &entity.User{
		Name:      req.GetName(),
		Age:       req.GetAge(),
		Gender:    req.GetGender(),
		About:     req.GetAbout(),
		Interests: mapReqInterests(req.GetInterestsId()),
	})
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateProfileResponse{Profile: &userv1.UserProfile{
		UserId:      user.ID,
		Name:        user.Name,
		Age:         user.Age,
		Gender:      user.Gender,
		Description: user.About,
		Photos:      mapPhotosToResp(user.Photos),
		Interests:   mapInterestsToResp(user.Interests),
	}}, nil
}

func mapPhotosToResp(photos []*entity.Photo) []*userv1.Photo {
	result := make([]*userv1.Photo, len(photos))
	for i, p := range photos {
		result[i] = &userv1.Photo{
			PhotoId:  p.ID,
			PhotoUrl: p.Url,
		}
	}
	return result
}

func mapInterestsToResp(interests []*entity.Interest) []*userv1.Interest {
	result := make([]*userv1.Interest, len(interests))
	for i, interest := range interests {
		result[i] = &userv1.Interest{
			InterestId: interest.ID,
			Name:       interest.Name,
		}
	}
	return result
}

func mapReqInterests(ids []int64) []*entity.Interest {
	var result []*entity.Interest
	for _, id := range ids {
		result = append(result, &entity.Interest{ID: id})
	}
	return result
}

func validateUserID(userID string) error {
	err := uuid.Validate(userID)
	if err != nil {
		return fmt.Errorf("user_id type is not uuid: %w", err)
	}
	return nil
}
