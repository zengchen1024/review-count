package main

import (
	"fmt"
	"time"
)

type reviewCount struct {
	excel           *excel
	reviewer        string
	maxReviewCount  int
	commentService  *reviewCommentService
	prSearchService *prSearchService
}

func (r *reviewCount) countComment() {
	count := 0
	n := 0

	for i, done := 0, false; count < r.maxReviewCount && !done; i++ {
		n, done = r.countOnce(i)
		count += n
	}

	fmt.Printf("total review count: %d\n", count)

	if err := r.excel.save(); err != nil {
		fmt.Printf("save file failed, err:%v\n", err)
	}
}

func (r *reviewCount) countOnce(times int) (count int, done bool) {
	v, err := r.prSearchService.searchPR(times)
	if err != nil {
		fmt.Printf("search pr failed, err:%v\n", err)

		return
	}

	done = v.complete()

	for i := range v.Items {
		item := &v.Items[i]

		v, err := r.countCommentOfPR(item)
		if err != nil {
			fmt.Printf("count comment of pr:%s failed, err:%v\n", item.desc(), err)
			return
		}

		count += v

		time.Sleep(100 * time.Millisecond)
	}

	return count, done
}

func (r *reviewCount) countCommentOfPR(item *prSearchItem) (int, error) {
	org, repo := item.orgRepo()
	comments, err := r.commentService.listComment(org, repo, item.Number)
	if err != nil {
		return 0, err
	}

	record := reviewRecord{pr: item.pullRequestURL()}

	count := 0
	for i := range comments {
		if comment := &comments[i]; comment.isTarget(r.reviewer) {
			record.commentURL = comment.Location
			record.comment = comment.Body

			if err := r.excel.write(&record); err != nil {
				return 0, err
			}

			count++
		}
	}

	return count, nil
}
