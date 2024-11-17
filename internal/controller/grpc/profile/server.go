package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/VmesteApp/auth-service/internal/entity"
	profilev1 "github.com/VmesteApp/protobuf/gen/go/profile"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Profile interface {
	VkProfile(context context.Context, userID uint64) (entity.VkProfile, error)
}

type serverApi struct {
	profilev1.UnimplementedProfileServiceServer
	profile Profile
}

func Register(gRPC *grpc.Server, profile Profile) {
	profilev1.RegisterProfileServiceServer(gRPC, &serverApi{profile: profile})
}

func (s *serverApi) GetVkID(context context.Context, req *profilev1.GetVkIDRequest) (*profilev1.GetVkIDResponse, error) {
	fmt.Println(uint64(req.UserID))
	fmt.Println(req)
	vkProfile, err := s.profile.VkProfile(context, uint64(req.UserID))

	if errors.Is(err, entity.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "failed get vk id")
	}

	return &profilev1.GetVkIDResponse{
		VkID: int64(vkProfile.VkID),
	}, nil
}
