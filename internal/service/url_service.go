package service

import "url-shortner/internal/repository"

type URLService struct {
	urlRepo *repository.URLRepository
}

func NewService(urlRepo *repository.URLRepository) *URLService {
	return &URLService{urlRepo: urlRepo}
}

func (s *URLService) SaveURL(urlToSave string, alias string) error {
	return s.urlRepo.SaveURL(urlToSave, alias)
}
func (s *URLService) GetURL(alias string) (string, error) {
	return s.urlRepo.GetURL(alias)
}
func (s *URLService) UpdateURL(alias string, newURL string) error {
	return s.urlRepo.UpdateURL(alias, newURL)
}
func (s *URLService) DeleteURL(alias string) error {
	return s.urlRepo.DeleteURL(alias)
}
func (s *URLService) GetAllURLs() (map[string]string, error) {
	return s.urlRepo.GetAllURLs()
}