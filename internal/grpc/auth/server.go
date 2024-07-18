package auth

import (
	"cerberus/internal/services/auth"
	"cerberus/internal/storage"
	"context"
	cerberusv1 "github.com/alexbro4u/contracts/gen/go/cerberus"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

const (
	emailRegexp = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emptyValue  = 0
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (UserID int64, err error)
	IsAdmin(
		ctx context.Context,
		userID int64,
	) (isAdmin bool, err error)
}

type serverApi struct {
	cerberusv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	cerberusv1.RegisterAuthServer(gRPC, &serverApi{auth: auth})
}

func (s *serverApi) Login(
	ctx context.Context,
	req *cerberusv1.LoginRequest,
) (*cerberusv1.LoginResponse, error) {

	if err := validateLoginReq(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &cerberusv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverApi) Register(
	ctx context.Context,
	req *cerberusv1.RegisterRequest,
) (*cerberusv1.RegisterResponse, error) {

	if err := validateRegisterReq(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &cerberusv1.RegisterResponse{
		Uid: userID,
	}, nil
}

func (s *serverApi) IsAdmin(
	ctx context.Context,
	req *cerberusv1.IsAdminRequest,
) (*cerberusv1.IsAdminResponse, error) {

	if err := validateIsAdminReq(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUid())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &cerberusv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(emailRegexp)
	return re.MatchString(email)
}

func validateLoginReq(req *cerberusv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if !isValidEmail(req.GetEmail()) {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegisterReq(req *cerberusv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if !isValidEmail(req.GetEmail()) {
		return status.Error(codes.InvalidArgument, "invalid email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateIsAdminReq(req *cerberusv1.IsAdminRequest) error {
	if req.GetUid() == emptyValue {
		return status.Error(codes.InvalidArgument, "uid is required")
	}

	return nil
}
