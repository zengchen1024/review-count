package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/opensourceways/server-common-lib/utils"
)

const searchPRURL = `https://api.github.com/search/issues?q=is:pr+created:>%s+commenter:%s&per_page=100&page=%d`

func newPRSearchService(cli *utils.HttpClient, token, reviewer string) *prSearchService {
	y, _, _ := time.Now().Date()
	return &prSearchService{
		cli:      cli,
		token:    token,
		reviewer: reviewer,
		firstDay: fmt.Sprintf("%d-01-01", y),
	}
}

type prSearchService struct {
	cli      *utils.HttpClient
	token    string
	reviewer string
	firstDay string
}

func (s *prSearchService) searchPR(page int) (*prSearchBody, error) {
	if page <= 0 {
		page = 1
	}

	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf(searchPRURL, s.firstDay, s.reviewer, page), nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	var data prSearchBody
	_, err = s.cli.ForwardTo(req, &data)

	return &data, err
}
