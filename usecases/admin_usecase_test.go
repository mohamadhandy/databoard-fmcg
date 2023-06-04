package usecases

// func Test_AdminUseCase_CreateAdmin_ShouldReturnAdmin(t *testing.T) {
// 	// create mock function CreateAdmin
// 	mockAdminRepository := &mocksRepositories.AdminRepositoryInterface{}

// 	mockResult := make(chan repositories.RepositoryResult[models.Admin])
// 	mockResult <- repositories.RepositoryResult[models.Admin]{
// 		Data: models.Admin{
// 			Name: "John Cena",
// 		},
// 	}

// 	mockAdminRepository.On("CreateAdmin", mock.Anything).Return(mockResult)

// 	useCase := InitAdminUsecase(mockAdminRepository)
// 	res := useCase.CreateAdmin(models.AdminRequest{
// 		Name: "John Cena",
// 	})
// 	assert.NotNil(t, res)
// }

// create unit test for usecase
