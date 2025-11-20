package converter

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	commonV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/common/v1"
	userV1 "github.com/crafty-ezhik/rocket-factory/shared/pkg/proto/user/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserRegInfoToModel - конвертирует регистрационные данные в сервисную модель
func UserRegInfoToModel(data *userV1.UserRegistrationInfo) serviceModel.UserRegistrationInfo {
	return serviceModel.UserRegistrationInfo{
		Info:     userInfoToModel(data.Info),
		Password: data.Password,
	}
}

// userInfoToModel - конвертирует commonV1.UserInfo в сервисную модель
func userInfoToModel(data *commonV1.UserInfo) serviceModel.UserInfo {
	return serviceModel.UserInfo{
		Login:               data.Login,
		Email:               data.Email,
		NotificationMethods: notificationMethodsToModel(data.NotificationMethod),
	}
}

// notificationMethodsToModel - конвертирует []*commonV1.NotificationMethod в сервисную модель
func notificationMethodsToModel(methods []*commonV1.NotificationMethod) []serviceModel.NotificationMethod {
	out := make([]serviceModel.NotificationMethod, len(methods))
	for i, method := range methods {
		out[i] = serviceModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}
	return out
}

// UserToProto - конвертирует сервисную модель пользователя в proto модель
func UserToProto(user serviceModel.User) *commonV1.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &commonV1.User{
		Uuid:      user.UUID.String(),
		Info:      userInfoToProto(user.Info),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func userInfoToProto(info serviceModel.UserInfo) *commonV1.UserInfo {
	return &commonV1.UserInfo{
		Login:              info.Login,
		Email:              info.Email,
		NotificationMethod: notificationMethodsToProto(info.NotificationMethods),
	}
}

func notificationMethodsToProto(methods []serviceModel.NotificationMethod) []*commonV1.NotificationMethod {
	out := make([]*commonV1.NotificationMethod, len(methods))
	for i, method := range methods {
		out[i] = &commonV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}
	return out
}
