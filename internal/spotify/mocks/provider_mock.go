package mocks

import (
	"context"
	"time"

	myspotify "github.com/planetfall/musicresearcher/internal/spotify"
	"github.com/stretchr/testify/mock"
)

type ProviderMock struct {
	mock.Mock
}

func (m *ProviderMock) NewClient(ctx context.Context,
	clientId string,
	clientSecret string) (myspotify.Client, time.Time, error) {

	args := m.Called(clientId, clientSecret)
	return args.Get(0).(myspotify.Client), args.Get(1).(time.Time), args.Error(2)
}
