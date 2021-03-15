package sharesecret

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"

	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateSecret(content string, customPwd bool, expire time.Time) (Secret, error) {
	panic("implement me")
}

func (m *MockRepository) RemoveSecret(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) RemoveSecretsExpired() (int64, error) {
	panic("implement me")
}

func (m *MockRepository) HasSecretWithCustomPwd(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepository) GetSecret(id string) (Secret, error) {
	args := m.Called(id)
	return args.Get(0).(Secret), args.Error(1)
}

func TestGetContentSecretNoPassRequiredToSeeSecret(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"
	id := "727d7040-aac7-4dc3-ab44-938bfba92ebd"

	mockRepo := new(MockRepository)
	mockRepo.
		On("HasSecretWithCustomPwd", id).
		Return(false, nil)
	mockRepo.
		On("GetSecret", id).
		Return(Secret{
			ID: id,
			Content: "cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922",
			CustomPwd: false,
			CreatedAt: time.Now(),
			ExpiredAt: time.Now(),
		}, nil)
	mockRepo.
		On("RemoveSecret", id).
		Return(nil)

	sut := NewSecretService(mockRepo, key, pass)
	cs, err := sut.GetContentSecret(id, "")

	assert.Nil(t, err)
	assert.Equal(t, "My name is Bernie", cs)
}

func TestGetContentSecretWithPassRequiredToSeeSecret(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"
	id := "727d7040-aac7-4dc3-ab44-938bfba92ebd"

	mockRepo := new(MockRepository)
	mockRepo.
		On("HasSecretWithCustomPwd", id).
		Return(true, nil)
	mockRepo.
		On("GetSecret", id).
		Return(Secret{
			ID: id,
			Content: "cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922",
			CustomPwd: false,
			CreatedAt: time.Now(),
			ExpiredAt: time.Now(),
		}, nil)
	mockRepo.
		On("RemoveSecret", id).
		Return(nil)

	sut := NewSecretService(mockRepo, key, pass)
	cs, err := sut.GetContentSecret(id, pass)

	assert.Nil(t, err)
	assert.Equal(t, "My name is Bernie", cs)
}

func TestGetContentSecretWithPassButIsNotRequiredToSeeSecret(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"
	id := "727d7040-aac7-4dc3-ab44-938bfba92ebd"

	mockRepo := new(MockRepository)
	mockRepo.
		On("HasSecretWithCustomPwd", id).
		Return(false, nil)

	sut := NewSecretService(mockRepo, key, pass)
	cs, err := sut.GetContentSecret(id, pass)

	assert.NotNil(t, err)
	assert.Equal(t, "the password is not required", err.Error())
	assert.Empty(t, cs)
}

func TestGetContentSecretButNotExistOrWasViewed(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"
	id := "727d7040-aac7-4dc3-ab44-938bfba92ebd"

	mockRepo := new(MockRepository)
	mockRepo.
		On("HasSecretWithCustomPwd", id).
		Return(false, ErrSecretNotFound)

	sut := NewSecretService(mockRepo, key, pass)
	cs, err := sut.GetContentSecret(id, pass)

	assert.NotNil(t, err)
	assert.Equal(t, "it either never existed or has already been viewed", err.Error())
	assert.Empty(t, cs)
}

func TestGetContentSecretWithPassRequiredToSeeSecretButMissingPassword(t *testing.T) {

	key := "11111111111111111111111111111111"
	pass := "@myPassword"
	id := "727d7040-aac7-4dc3-ab44-938bfba92ebd"

	mockRepo := new(MockRepository)
	mockRepo.
		On("HasSecretWithCustomPwd", id).
		Return(true, nil)
	mockRepo.
		On("GetSecret", id).
		Return(Secret{
			ID: id,
			Content: "cb98267468c271c1a09bd6d03a919a2af89e9bde934b409258e9e462e2a7b312a9e6cb4d92582155f7a7c48922",
			CustomPwd: false,
			CreatedAt: time.Now(),
			ExpiredAt: time.Now(),
		}, nil)
	mockRepo.
		On("RemoveSecret", id).
		Return(nil)

	sut := NewSecretService(mockRepo, key, pass)
	cs, err := sut.GetContentSecret(id, "")

	assert.NotNil(t, err)
	assert.Empty(t, cs)
	assert.Equal(t, "you need a password to see the secret", err.Error())
}

func TestBadKeyInServiceDistinctToLen32(t *testing.T) {

	defer func() { recover() }()

	key := "111111"
	pass := "@myPassword"

	mockRepo := new(MockRepository)
	NewSecretService(mockRepo, key, pass)
	assert.Fail(t, "should have panicked")
}
