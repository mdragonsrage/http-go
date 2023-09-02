package house

import "errors"

type InMemoryStorage struct {
	data map[string]bool
}

func (db *InMemoryStorage) GetAll() ([]Lightbulb, error) {
	lbs := []Lightbulb{}

	for k, v := range db.data {
		lbs = append(lbs, Lightbulb{Name: k, On: v})
	}

	return lbs, nil
}

func (db *InMemoryStorage) Get(name string) (Lightbulb, error) {
	if val, exists := db.data[name]; exists {
		return Lightbulb{Name: name, On: val}, nil
	}

	return Lightbulb{}, errors.New("key not dound")
}

func (db *InMemoryStorage) Create(lb Lightbulb) error {
	db.data[lb.Name] = lb.On
	return nil
}

func (db *InMemoryStorage) Update(lb Lightbulb) error {
	db.data[lb.Name] = lb.On
	return nil
}

func (db *InMemoryStorage) Delete(name string) error {
	delete(db.data, name)
	return nil
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: map[string]bool{},
	}
}
