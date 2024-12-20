package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opensourceways/server-common-lib/utils"
)

type options struct {
	token          string
	reviewer       string
	maxReviewCount int
}

func (o *options) validate() error {
	if o.token == "" {
		return fmt.Errorf("miss token parameter")
	}

	if o.reviewer == "" {
		return fmt.Errorf("miss reviewer parameter")
	}

	if o.maxReviewCount <= 0 {
		return fmt.Errorf("miss max review count")
	}

	return nil
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	fs.StringVar(&o.token, "token", "", "github token")

	fs.StringVar(&o.reviewer, "reviewer", "", "github id of reviewer")

	fs.IntVar(&o.maxReviewCount, "mx_review_count", 500, "max review number that will be counted up to")

	fs.Parse(args)
	return o
}

func main() {
	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.validate(); err != nil {
		fmt.Printf("Invalid options, err:%s\n", err.Error())

		return
	}

	cli := utils.NewHttpClient(3)

	h := &reviewCount{
		prSearchService: newPRSearchService(&cli, o.token, o.reviewer),
		commentService:  newReviewCommentService(&cli, o.token),
		maxReviewCount:  500,
		reviewer:        o.reviewer,
	}

	h.countComment()
}
