package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type debugRoundTripper struct {
	roundTripper http.RoundTripper
}

func (m *debugRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	dump, err := httputil.DumpRequestOut(request, false)
	if err == nil {
		log.Println("Request --->")
		log.Println(dump)
	}

	resp, err := m.roundTripper.RoundTrip(request)
	if err == nil {
		dump, err := httputil.DumpResponse(resp, false)
		if err == nil {
			log.Println(" <--- Response")
			log.Printf("%q", dump)
		}
	}

	return resp, err
}

type Client struct {
	client *http.Client
	url    string
	debug  bool
}

func NewClient(url string, debug bool) *Client {
	c := http.DefaultClient
	if debug {
		c = &http.Client{
			Transport: &debugRoundTripper{roundTripper: http.DefaultTransport},
		}
	}

	return &Client{
		client: c,
		url:    url,
		debug:  debug,
	}
}

func (c *Client) SaveFile(filename string, data []byte) error {
	queryURL, err := url.JoinPath(c.url, "upload", filename)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPut, queryURL, bytes.NewReader(data))
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)

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

func (c *Client) Get(filename string) ([]byte, error) {
	queryURL, err := url.JoinPath(c.url, "download", filename)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, queryURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func testUploadAndDownload(url string, filepath string, debug bool) error {
	c := NewClient(url, debug)

	body, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("can't read file error: %w", err)
	}

	fileName := "test_file"

	err = c.SaveFile(fileName, body)
	if err != nil {
		return fmt.Errorf("can't save file error: %w", err)
	}

	storageFile, err := c.Get(fileName)
	if err != nil {
		return fmt.Errorf("can't download file from storage error: %w", err)
	}

	if !bytes.Equal(body, storageFile) {
		return errors.New("uploaded data isn't equal to data from download")
	}

	return nil
}

func main() {
	url := flag.String("url", "http://localhost:8080", "server url")
	filePath := flag.String("file-path", "1mb_file", "path to file for upload")
	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	err := testUploadAndDownload(*url, *filePath, *debug)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
