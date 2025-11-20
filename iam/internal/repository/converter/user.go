package converter

import (
	serviceModel "github.com/crafty-ezhik/rocket-factory/iam/internal/model"
	repoModel "github.com/crafty-ezhik/rocket-factory/iam/internal/repository/model"
)

func UserToServiceModel(user repoModel.User) serviceModel.User {
	return serviceModel.User{
		UUID: user.UUID,
		Info: serviceModel.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: notificationMethodsToModel(user.Info.NotificationMethods),
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserRegInfoToRepoModel(data serviceModel.UserRegistrationInfo) repoModel.UserRegistrationInfo {
	return repoModel.UserRegistrationInfo{
		Info: repoModel.UserInfo{
			Login:               data.Info.Login,
			Email:               data.Info.Email,
			NotificationMethods: notificationMethodsToRepo(data.Info.NotificationMethods),
		},
	}
}

func notificationMethodsToRepo(methods []serviceModel.NotificationMethod) []repoModel.NotificationMethod {
	out := make([]repoModel.NotificationMethod, len(methods))
	for i, method := range methods {
		out[i] = repoModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}
	return out
}

func notificationMethodsToModel(methods []repoModel.NotificationMethod) []serviceModel.NotificationMethod {
	out := make([]serviceModel.NotificationMethod, len(methods))
	for i, method := range methods {
		out[i] = serviceModel.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		}
	}
	return out
}
