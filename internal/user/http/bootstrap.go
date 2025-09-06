package user_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	storage_database_local "github.com/gabrielmrtt/taski/internal/storage/database/local"
	storage_database_postgres "github.com/gabrielmrtt/taski/internal/storage/database/postgres"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(g *gin.RouterGroup) {
	userRepository := user_database_postgres.NewUserPostgresRepository()
	userRegistrationRepository := user_database_postgres.NewUserRegistrationPostgresRepository()
	passwordRecoveryRepository := user_database_postgres.NewPasswordRecoveryPostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()
	uploadedFileRepository := storage_database_postgres.NewUploadedFilePostgresRepository()
	storageRepository := storage_database_local.NewLocalStorageRepository()

	userLoginService := user_services.NewUserLoginService(userRepository)
	registerUserService := user_services.NewRegisterUserService(userRepository, userRegistrationRepository, transactionRepository)
	verifyUserRegistrationService := user_services.NewVerifyUserRegistrationService(userRegistrationRepository, userRepository, transactionRepository)
	forgotUserPasswordService := user_services.NewForgotUserPasswordService(userRepository, passwordRecoveryRepository, transactionRepository)
	recoverUserPasswordService := user_services.NewRecoverUserPasswordService(userRepository, passwordRecoveryRepository, transactionRepository)
	getMeService := user_services.NewGetMeService(userRepository)
	changeUserPasswordService := user_services.NewChangeUserPasswordService(userRepository, transactionRepository)
	updateUserCredentialsService := user_services.NewUpdateUserCredentialsService(userRepository, transactionRepository)
	updateUserDataService := user_services.NewUpdateUserDataService(userRepository, transactionRepository, uploadedFileRepository, storageRepository)
	deleteUserService := user_services.NewDeleteUserService(userRepository, transactionRepository)

	userController := NewUserController(getMeService, changeUserPasswordService, updateUserCredentialsService, updateUserDataService, deleteUserService)
	userRegistrationController := NewUserRegistrationController(registerUserService, verifyUserRegistrationService, forgotUserPasswordService, recoverUserPasswordService)
	authController := NewAuthController(userLoginService)

	userController.ConfigureRoutes(g)
	userRegistrationController.ConfigureRoutes(g)
	authController.ConfigureRoutes(g)
}
