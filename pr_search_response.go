package main

import (
	"fmt"
	"strings"
)

type prSearchItem struct {
	Number        int    `json:"number"`
	RepositoryURL string `json:"repository_url"`
}

func (i *prSearchItem) commentURL() string {
	org, repo := i.orgRepo()

	return fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls/%d/comments",
		org, repo, i.Number,
	)
}

func (i *prSearchItem) orgRepo() (string, string) {
	v := strings.Split(i.RepositoryURL, "/")
	n := len(v) - 1
	return v[n-1], v[n]
}

func (i *prSearchItem) desc() string {
	org, repo := i.orgRepo()

	return fmt.Sprintf("%s/%s:%d", org, repo, i.Number)
}

func (i *prSearchItem) pullRequestURL() string {
	org, repo := i.orgRepo()

	return fmt.Sprintf("https://github.com/%s/%s/pull/%d", org, repo, i.Number)
}

type prSearchBody struct {
	TotalCount int            `json:"total_count"`
	Incomplete bool           `json:"incomplete_results"`
	Items      []prSearchItem `json:"items"`
}

func (d *prSearchBody) validate() error {
	if len(d.Items) != d.TotalCount {
		return fmt.Errorf(
			"count does not match, expect:%v != actual:%v",
			d.TotalCount, len(d.Items),
		)
	}

	return nil
}

func (d *prSearchBody) showOne() {
	if len(d.Items) == 0 {
		fmt.Println("NO DATA")
	} else {
		fmt.Println(d.Items[0].desc())
	}
}

func (d *prSearchBody) complete() bool {
	return !d.Incomplete
}
