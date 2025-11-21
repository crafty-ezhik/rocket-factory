package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	authV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/common/v1"
)

func WhoamiResponseToProto(data model.WhoamiResponse, sessionUUID string) *authV1.WhoamiResponse {
	return &authV1.WhoamiResponse{
		Session: SessionToProto(data.Session, sessionUUID),
		User:    UserToProto(data.User),
	}
}

func SessionToProto(data model.Session, sessionUUID string) *commonV1.Session {
	var updatedAt *timestamppb.Timestamp
	if data.UpdatedAt != nil {
		updatedAt = timestamppb.New(*data.UpdatedAt)
	}

	return &commonV1.Session{
		Uuid:      sessionUUID,
		CreatedAt: timestamppb.New(data.CreatedAt),
		UpdatedAt: updatedAt,
		ExpiresAt: timestamppb.New(data.ExpiresAt),
	}
}
