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

func (r *FakeDogRepository) Create(dog domain.Dog) (*domain.Dog, error) {
	dog.ID = uuid.New().String()
	r.dogs = append(r.dogs, dog)
	return &dog, nil
}

func (r *FakeDogRepository) List(limit int, name string, nextToken string) (*domain.DogList, error) {
	return &domain.DogList{
		Dogs:      r.dogs,
		NextToken: "",
	}, nil
}

func (r *FakeDogRepository) Get(id string) (*domain.Dog, error) {
	for _, dog := range r.dogs {
		if dog.ID == id {
			return &dog, nil
		}
	}
	return nil, domain.ErrDogNotFound
}

func (r *FakeDogRepository) Update(id string, dog *domain.Dog) error {
	for i, d := range r.dogs {
		if d.ID == id {
			r.dogs[i] = *dog
			return nil
		}
	}
	return domain.ErrDogNotFound
}

func (r *FakeDogRepository) Delete(id string) error {
	for i, d := range r.dogs {
		if d.ID == id {
			r.dogs = append(r.dogs[:i], r.dogs[i+1:]...)
			return nil
		}
	}
	return domain.ErrDogNotFound
}

func (r *FakeDogRepository) UpdatePhotoHash(id string, photoHash string) error {
	for i, d := range r.dogs {
		if d.ID == id {
			r.dogs[i].PhotoHash = photoHash
			return nil
		}
	}
	return domain.ErrDogNotFound
}

func (r *FakeDogRepository) UpdatePhotoStatus(id string, photoStatus string) error {
	for i, d := range r.dogs {
		if d.ID == id {
			r.dogs[i].PhotoStatus = photoStatus
			return nil
		}
	}
	return domain.ErrDogNotFound
}
