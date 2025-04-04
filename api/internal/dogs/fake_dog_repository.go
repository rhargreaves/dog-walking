package dogs

import (
	"github.com/google/uuid"
	"github.com/rhargreaves/dog-walking/api/internal/dogs/domain"
)

type FakeDogRepository struct {
	domain.DogRepository
	dogs []domain.Dog
}

func NewFakeDogRepository() *FakeDogRepository {
	return &FakeDogRepository{}
}

func (f *FakeDogRepository) Create(dog domain.Dog) (*domain.Dog, error) {
	dog.ID = uuid.New().String()
	f.dogs = append(f.dogs, dog)
	return &dog, nil
}

func (f *FakeDogRepository) List(limit int, name string, nextToken string) (*domain.DogList, error) {
	return &domain.DogList{
		Dogs:      f.dogs,
		NextToken: "",
	}, nil
}

func (f *FakeDogRepository) Get(id string) (*domain.Dog, error) {
	for _, dog := range f.dogs {
		if dog.ID == id {
			return &dog, nil
		}
	}
	return nil, domain.ErrDogNotFound
}

func (f *FakeDogRepository) Update(id string, dog *domain.Dog) error {
	for i, d := range f.dogs {
		if d.ID == id {
			f.dogs[i] = *dog
			return nil
		}
	}
	return domain.ErrDogNotFound
}

func (f *FakeDogRepository) Delete(id string) error {
	for i, d := range f.dogs {
		if d.ID == id {
			f.dogs = append(f.dogs[:i], f.dogs[i+1:]...)
			return nil
		}
	}
	return domain.ErrDogNotFound
}

func (f *FakeDogRepository) UpdatePhotoHash(id string, photoHash string) error {
	for i, d := range f.dogs {
		if d.ID == id {
			f.dogs[i].PhotoHash = photoHash
			return nil
		}
	}
	return domain.ErrDogNotFound
}

func (f *FakeDogRepository) UpdatePhotoStatus(id string, photoStatus string) error {
	for i, d := range f.dogs {
		if d.ID == id {
			f.dogs[i].PhotoStatus = photoStatus
			return nil
		}
	}
	return domain.ErrDogNotFound
}
