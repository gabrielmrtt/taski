package user_http

import (
	core_database_postgres "github.com/gabrielmrtt/taski/internal/core/database/postgres"
	user_database_postgres "github.com/gabrielmrtt/taski/internal/user/database/postgres"
	user_services "github.com/gabrielmrtt/taski/internal/user/services"
	"github.com/gin-gonic/gin"
)

func BootstrapControllers(engine *gin.Engine) {
	userRepository := user_database_postgres.NewUserPostgresRepository()
	userRegistrationRepository := user_database_postgres.NewUserRegistrationPostgresRepository()
	passwordRecoveryRepository := user_database_postgres.NewPasswordRecoveryPostgresRepository()
	transactionRepository := core_database_postgres.NewTransactionPostgresRepository()

	userLoginService := user_services.NewUserLoginService(userRepository)
	registerUserService := user_services.NewRegisterUserService(userRepository, userRegistrationRepository, transactionRepository)
	verifyUserRegistrationService := user_services.NewVerifyUserRegistrationService(userRegistrationRepository, userRepository, transactionRepository)
	forgotUserPasswordService := user_services.NewForgotUserPasswordService(userRepository, passwordRecoveryRepository, transactionRepository)
	recoverUserPasswordService := user_services.NewRecoverUserPasswordService(userRepository, passwordRecoveryRepository, transactionRepository)
	getMeService := user_services.NewGetMeService(userRepository)
	changeUserPasswordService := user_services.NewChangeUserPasswordService(userRepository, transactionRepository)
	updateUserCredentialsService := user_services.NewUpdateUserCredentialsService(userRepository, transactionRepository)
	updateUserDataService := user_services.NewUpdateUserDataService(userRepository, transactionRepository)
	deleteUserService := user_services.NewDeleteUserService(userRepository, transactionRepository)

	userController := NewUserController(getMeService, changeUserPasswordService, updateUserCredentialsService, updateUserDataService, deleteUserService)
	userRegistrationController := NewUserRegistrationController(registerUserService, verifyUserRegistrationService, forgotUserPasswordService, recoverUserPasswordService)
	authController := NewAuthController(userLoginService)

	userController.ConfigureRoutes(engine)
	userRegistrationController.ConfigureRoutes(engine)
	authController.ConfigureRoutes(engine)
}
