package dreamkast

import (
	"context"
	"dreamkast-weaver/internal/dkui/domain"
	"dreamkast-weaver/internal/dkui/value"
	"dreamkast-weaver/internal/stacktrace"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DkApiClient interface {
	GetTracks(ctx context.Context, n value.ConfName) (*domain.ViewerCounts, error)
}

type DkApiClientImpl struct {
	dkApiAddr string
}

var _ DkApiClient = (*DkApiClientImpl)(nil)

func NewDkApiClientImpl() DkApiClient {
	addr := os.Getenv("DREAMKAST_ADDR")
	if addr == "" {
		addr = "https://staging.dev.cloudnativedays.jp"
	}
	return &DkApiClientImpl{
		dkApiAddr: addr,
	}
}

type Track struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ChannelArn string `json:"channelArn"`
	// VideoPlatform string    `json:"videoPlatform"`
	// VideoId       string    `json:"videoId"`
	// OnAirTalk     OnAirTalk `json:"onAirTalk"`
}

func (d *DkApiClientImpl) GetTracks(ctx context.Context, n value.ConfName) (*domain.ViewerCounts, error) {
	url := d.dkApiAddr + "/api/v1/tracks/?eventAbbr=" + n.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	rowResp, err := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, stacktrace.With(fmt.Errorf("GetTracks: %w", err))
	}
	defer rowResp.Body.Close()

	if rowResp.StatusCode != 200 {
		return nil, stacktrace.With(fmt.Errorf("Error: status code %w", err))
	}

	body, _ := io.ReadAll(rowResp.Body)

	var gtResp []Track
	if err := json.Unmarshal(body, &gtResp); err != nil {
		return nil, stacktrace.With(err)
	}

	return viewerCountConv.fromDk(gtResp)
}

var viewerCountConv _viewerCountConv

type _viewerCountConv struct{}

func (_viewerCountConv) fromDk(v []Track) (*domain.ViewerCounts, error) {
	conv := func(v *Track) (*domain.ViewerCount, error) {
		trackID, err := value.NewTrackID(int32(v.Id))
		if err != nil {
			return nil, err
		}
		ca, err := value.NewChannelArn(v.ChannelArn)
		if err != nil {
			return nil, err
		}
		tn, err := value.NewTrackName(v.Name)
		if err != nil {
			return nil, err
		}

		dvc := domain.NewViewerCount(trackID, ca, tn, 0)
		return dvc, nil
	}

	var items []domain.ViewerCount
	for _, p := range v {
		vc := p
		dvc, err := conv(&vc)
		if err != nil {
			return nil, stacktrace.With(fmt.Errorf("convert view count from dreamkast: %w", err))
		}
		items = append(items, *dvc)
	}

	return &domain.ViewerCounts{Items: items}, nil
}
