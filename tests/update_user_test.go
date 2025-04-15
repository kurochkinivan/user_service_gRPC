package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	userv1 "github.com/kurochkinivan/user_proto/gen/go/users"
	"github.com/kurochkinivan/user_service/tests/suite"
	"github.com/stretchr/testify/require"
)

const (
	existingUUID       = "44d34c46-d165-4585-a98b-001eca1366c1"
	existingInterestID = 1
)

func TestUpdate_User_DoesNot_Exist_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	req := &userv1.UpdateProfileRequest{
		UserId:      gofakeit.UUID(),
		Name:        gofakeit.Name(),
		Age:         gofakeit.Int32(),
		Gender:      gofakeit.Gender(),
		About:       gofakeit.LoremIpsumSentence(30),
		InterestsId: []int64{existingInterestID},
	}

	resp, err := st.UserClient.UpdateProfile(ctx, req)
	require.NoError(t, err)

	require.NotNil(t, resp)
	require.NotNil(t, resp.Profile)
	require.Equal(t, req.UserId, resp.Profile.UserId)
	require.Equal(t, req.Name, resp.Profile.Name)
	require.Equal(t, req.Age, resp.Profile.Age)
	require.Equal(t, req.Gender, resp.Profile.Gender)
	require.Equal(t, req.About, resp.Profile.Description)
	require.Len(t, resp.Profile.Interests, len(req.InterestsId))
}

func TestUpdate_User_Exists_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	req := &userv1.UpdateProfileRequest{
		UserId:      existingUUID,
		Name:        gofakeit.Name(),
		Age:         gofakeit.Int32(),
		Gender:      gofakeit.Gender(),
		About:       gofakeit.LoremIpsumSentence(30),
		InterestsId: []int64{existingInterestID},
	}

	resp, err := st.UserClient.UpdateProfile(ctx, req)
	require.NoError(t, err)

	require.NotNil(t, resp)
	require.NotNil(t, resp.Profile)
	require.Equal(t, existingUUID, resp.Profile.UserId)
	require.Equal(t, req.Name, resp.Profile.Name)
	require.Equal(t, req.Age, resp.Profile.Age)
	require.Equal(t, req.Gender, resp.Profile.Gender)
	require.Equal(t, req.About, resp.Profile.Description)
}
