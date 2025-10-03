package userinfra

import (
	coredatabase "github.com/gabrielmrtt/taski/internal/core/database"
	storagedatabase "github.com/gabrielmrtt/taski/internal/storage/infra/database"
	userdatabase "github.com/gabrielmrtt/taski/internal/user/infra/database"
	userhttp "github.com/gabrielmrtt/taski/internal/user/infra/http"
	userservice "github.com/gabrielmrtt/taski/internal/user/service"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type BootstrapInfraOptions struct {
	RouterGroup  *gin.RouterGroup
	DbConnection *bun.DB
}

func BootstrapInfra(options BootstrapInfraOptions) {
	userRepository := userdatabase.NewUserBunRepository(options.DbConnection)
	userRegistrationRepository := userdatabase.NewUserRegistrationBunRepository(options.DbConnection)
	passwordRecoveryRepository := userdatabase.NewPasswordRecoveryBunRepository(options.DbConnection)
	transactionRepository := coredatabase.NewTransactionBunRepository(options.DbConnection)
	uploadedFileRepository := storagedatabase.NewUploadedFileBunRepository(options.DbConnection)
	storageRepository := storagedatabase.NewLocalStorageRepository()

	userLoginService := userservice.NewUserLoginService(userRepository)
	registerUserService := userservice.NewRegisterUserService(userRepository, userRegistrationRepository, transactionRepository)
	verifyUserRegistrationService := userservice.NewVerifyUserRegistrationService(userRegistrationRepository, userRepository, transactionRepository)
	forgotUserPasswordService := userservice.NewForgotUserPasswordService(userRepository, passwordRecoveryRepository, transactionRepository)
	recoverUserPasswordService := userservice.NewRecoverUserPasswordService(userRepository, passwordRecoveryRepository, transactionRepository)
	getMeService := userservice.NewGetMeService(userRepository)
	changeUserPasswordService := userservice.NewChangeUserPasswordService(userRepository, transactionRepository)
	updateUserCredentialsService := userservice.NewUpdateUserCredentialsService(userRepository, transactionRepository)
	updateUserDataService := userservice.NewUpdateUserDataService(userRepository, transactionRepository, uploadedFileRepository, storageRepository)
	deleteUserService := userservice.NewDeleteUserService(userRepository, transactionRepository)

	userController := userhttp.NewUserHandler(getMeService, changeUserPasswordService, updateUserCredentialsService, updateUserDataService, deleteUserService)
	userRegistrationController := userhttp.NewUserRegistrationHandler(registerUserService, verifyUserRegistrationService, forgotUserPasswordService, recoverUserPasswordService)
	authController := userhttp.NewAuthHandler(userLoginService)

	userController.ConfigureRoutes(options.RouterGroup)
	userRegistrationController.ConfigureRoutes(options.RouterGroup)
	authController.ConfigureRoutes(options.RouterGroup)
}
