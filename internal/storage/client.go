package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/dyus/filebalancer/internal/db"
)

type Status struct {
	UsedSpace int64 `json:"usedSpace"`
}

type Client struct {
	client *http.Client
}

func (s *Client) Write(ctx context.Context, chunk *db.FilePart, data io.Reader) error {
	queryURL, err := url.JoinPath(chunk.Storage, "upload", chunk.Path)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, queryURL, data)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *Client) Read(ctx context.Context, chunk *db.FilePart) (io.ReadCloser, error) {
	queryURL, err := url.JoinPath(chunk.Storage, "download", chunk.Path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		if resp != nil {
			defer resp.Body.Close()
		}

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (s *Client) Delete(ctx context.Context, chunk *db.FilePart) error {
	queryURL, err := url.JoinPath(chunk.Storage, chunk.Path)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, queryURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func NewClient() *Client {
	return &Client{
		client: http.DefaultClient,
	}
}
