package main

import (
	"fmt"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"
)

const listCommentURL = `https://api.github.com/repos/%s/%s/pulls/%d/comments`

func newReviewCommentService(cli *utils.HttpClient, token string) *reviewCommentService {
	return &reviewCommentService{
		cli:   cli,
		token: token,
	}
}

type reviewCommentService struct {
	cli   *utils.HttpClient
	token string
}

func (h *reviewCommentService) listComment(org, repo string, prNum int) ([]commentListItem, error) {
	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf(listCommentURL, org, repo, prNum), nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+h.token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	var data []commentListItem
	_, err = h.cli.ForwardTo(req, &data)

	return data, err
}
