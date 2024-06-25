package services

import (
	"fmt"

	"github.com/Elfsilon/car_booking/internal/bookings/models"
)

type CarsService struct {
	cars map[string]models.CarInfo
}

func NewMockCarsService() *CarsService {
	return &CarsService{
		cars: map[string]models.CarInfo{
			"c84fde82-6679-4384-af11-7406de3d3e14": {Sign: "Н314ХО123", Name: "Lada Granta", Color: "White"},
			"e78b2415-b47c-435a-91e2-655ec5a08023": {Sign: "М265ДЫ123", Name: "Lada Vesta", Color: "Blue"},
			"b46cbaa7-d02f-4571-ab6e-1883813715bf": {Sign: "К159ЕК93", Name: "Kia Rio", Color: "White"},
			"150bd48b-1671-469a-822c-cc236d670a45": {Sign: "Е358ВА93", Name: "Mitsubishi Lancer", Color: "Red"},
			"542af38b-c4b9-4d74-a7fa-e21c2a50a8cf": {Sign: "Х777ХХ123", Name: "Mercedes-Benz C-class", Color: "Black"},
		},
	}
}

func (s *CarsService) GetCarInfo(carID string) (models.CarInfo, error) {
	info, ok := s.cars[carID]
	if !ok {
		return models.CarInfo{}, fmt.Errorf("car with provided id=%v not found", carID)
	}
	return info, nil
}
