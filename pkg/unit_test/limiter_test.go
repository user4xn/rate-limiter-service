package unittest

import (
	"context"
	"net/http"
	"rate-limiter/pkg/app/limiter"
	"rate-limiter/pkg/dto"
	"rate-limiter/pkg/factory"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTest() limiter.Service {
	f := factory.NewFactory()
	svc := limiter.NewService(f)
	return svc
}

func TestFixedWindow(t *testing.T) {
	tests := []struct {
		name           string
		request        int
		payload        dto.PayloadLimiter
		expectedError  error
		expectedResult *dto.ResponseLimiter
	}{
		{
			name:    "First request - success",
			request: 1,
			payload: dto.PayloadLimiter{
				ClientID: "wildan1",
				Route:    "/api/v1/transactions",
			},
			expectedError: nil,
			expectedResult: &dto.ResponseLimiter{
				Status: http.StatusOK,
				Limit:  100,
				Remain: 99,
			},
		},
		{
			name:    "Limit exceeded - success",
			request: 100,
			payload: dto.PayloadLimiter{
				ClientID: "wildan1",
				Route:    "/api/v1/users",
			},
			expectedError: nil,
			expectedResult: &dto.ResponseLimiter{
				Status: http.StatusTooManyRequests,
				Limit:  100,
				Remain: 0,
			},
		},
		{
			name:    "Invalid payload - error",
			request: 1,
			payload: dto.PayloadLimiter{
				Route: "/api/v1/track",
			},
			expectedError:  assert.AnError,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := setupTest()

			if tt.request > 1 {
				for i := 0; i < tt.request; i++ {
					_, _ = svc.FixedWindow(context.Background(), tt.payload)
				}
			}

			resp, err := svc.FixedWindow(context.Background(), tt.payload)

			if tt.expectedError == nil {
				t.Logf("Request: %d", tt.request)
				t.Logf("ClientID: %s", tt.payload.ClientID)
				t.Logf("Route: %s", tt.payload.Route)
				t.Logf("Expected Remain: %d", tt.expectedResult.Remain)
				if resp != nil {
					t.Logf("Actual Remain: %d", resp.Remain)
				}
			}

			if tt.expectedError != nil {
				assert.Error(t, err)
				t.Logf("Excpected Error: %v", err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResult.Status, resp.Status)
				assert.Equal(t, tt.expectedResult.Remain, resp.Remain)
			}
		})
	}
}

func TestUpdateConfig(t *testing.T) {
	svc := setupTest()
	config := dto.PayloadConfigClient{
		ClientID: "wildan1",
		Route:    "/api/v1/orders",
		Limit:    50,
		Window:   30,
	}
	err := svc.SetClientConfigFixedWindow(context.Background(), config)
	assert.NoError(t, err)
	t.Logf("Config updated: %+v", config)

	payload := dto.PayloadLimiter{
		ClientID: "wildan1",
		Route:    "/api/v1/orders",
	}

	for i := 0; i < 50; i++ {
		_, _ = svc.FixedWindow(context.Background(), payload)
	}

	resp, err := svc.FixedWindow(context.Background(), payload)

	t.Logf("Expected Remain: 0")
	if resp != nil {
		t.Logf("Actual Remain: %d", resp.Remain)
	}

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Remain, 0)
}

func TestFixedWindow_Burst(t *testing.T) {
	svc := setupTest()

	config := dto.PayloadConfigClient{
		ClientID: "wildan-burst",
		Route:    "/api/v1/burst",
		Limit:    10,
		Window:   60,
	}
	err := svc.SetClientConfigFixedWindow(context.Background(), config)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	failCount := 0

	burstCount := 20 //simulate 20 burst request concurrently
	for i := 0; i < burstCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, _ := svc.FixedWindow(context.Background(), dto.PayloadLimiter{
				ClientID: config.ClientID,
				Route:    config.Route,
			})
			mu.Lock()
			defer mu.Unlock()
			if resp != nil && resp.Status == http.StatusOK {
				successCount++
			} else if resp != nil && resp.Status == http.StatusTooManyRequests {
				failCount++

			}
		}()
	}

	wg.Wait()

	t.Logf("Success: %d, Failed: %d", successCount, failCount)

	assert.Equal(t, config.Limit, successCount)
	assert.Equal(t, burstCount-config.Limit, failCount)
}
