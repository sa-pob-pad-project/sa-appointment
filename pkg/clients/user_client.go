package clients

import (
	"appointment-service/pkg/clients/dto"
	contextUtils "appointment-service/pkg/context"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserClient struct {
	baseUrl string
	hc      *http.Client
}

func New(baseUrl string) *UserClient {
	return &UserClient{
		baseUrl: baseUrl,
		hc: &http.Client{
			Timeout: http.DefaultClient.Timeout,
		},
	}
}

func (c *UserClient) GetDoctorById(ctx context.Context, doctorID string) (*dto.GetDoctorProfileResponseDto, error) {
	var doctorProfile dto.GetDoctorProfileResponseDto
	accessToken := contextUtils.GetAccessToken(ctx)
	if accessToken == "" {
		return nil, fmt.Errorf("access token is empty")
	}
	url := fmt.Sprintf("%s/v1/doctor/%s", c.baseUrl, doctorID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: accessToken,
	})
	fmt.Println("Requesting URL:", url)
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("status", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&doctorProfile); err != nil {
		return nil, err
	}

	return &doctorProfile, nil
}
