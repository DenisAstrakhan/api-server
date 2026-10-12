package web_service

type WebService struct {
	webRepository WebRepository
}

type WebRepository interface {
	GetFile(filePath string) ([]byte, error)
}

func MewWebService(repository WebRepository) WebService {
	return WebService{
		webRepository: repository,
	}
}
