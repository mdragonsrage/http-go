package house

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockStorage struct {
	err error
	lb  Lightbulb
	lbs []Lightbulb
}

func (m *MockStorage) GetAll() ([]Lightbulb, error) {
	return m.lbs, m.err
}

func (m *MockStorage) Get(name string) (Lightbulb, error) {
	return m.lb, m.err
}

func (m *MockStorage) Create(lightbulb Lightbulb) error {
	return m.err
}

func (m *MockStorage) Update(lightbulb Lightbulb) error {
	return m.err
}

func (m *MockStorage) Delete(name string) error {
	return m.err
}

func TestCreateLightbulb(t *testing.T) {
	type args struct {
		storage Storage
		r       func() *http.Request
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{
			name: "create_returns_200_when_all_good",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewReader([]byte(`{"name":"livingroom"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "create_returns_500_when_storage_misbehaves",
			args: args{
				storage: &MockStorage{
					err: errors.New("something's wrong"),
				},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/create", bytes.NewReader([]byte(`{"name":"livingroom"}`)))
					return req
				},
			},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name: "create_returns_400_when_request_body_is_invalid",
			args: args{
				storage: &MockStorage{},
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodPost, "/create", nil)
					return req
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateLightbulb(tt.args.storage)
			w := httptest.NewRecorder()
			handler(w, tt.args.r())
			result := w.Result()
			if result.StatusCode != tt.wantStatusCode {
				t.Errorf("Create() statusCode = %v, wantStatusCode = %v", result.StatusCode, tt.wantStatusCode)

			}
		})
	}
}
