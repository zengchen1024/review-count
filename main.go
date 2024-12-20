package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/opensourceways/server-common-lib/utils"
)

type options struct {
	help           bool
	token          string
	reviewer       string
	fileName       string
	maxReviewCount int
}

func (o *options) validate() error {
	if o.token == "" {
		return fmt.Errorf("miss token parameter")
	}

	if o.reviewer == "" {
		return fmt.Errorf("miss reviewer parameter")
	}

	if o.fileName == "" {
		return fmt.Errorf("miss file-name parameter")
	}

	if o.maxReviewCount <= 0 {
		return fmt.Errorf("miss max review count")
	}

	return nil
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	fs.BoolVar(&o.help, "help", false, "show usage")

	fs.StringVar(&o.token, "token", "", "github token")

	fs.StringVar(&o.reviewer, "reviewer", "", "github id of reviewer")

	fs.StringVar(&o.fileName, "file-name", "", "file path where the review details will be writen to")

	fs.IntVar(&o.maxReviewCount, "mx-review-count", 500, "max review number that will be counted up to")

	fs.Parse(args)

	return o
}

func usage() {
	fmt.Println("\nusage\n    ./review-count --token=xx --reviewer=xx --file-name=xx [--max-review-count=xx]\n")
}

func main() {
	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)

	if o.help {
		usage()

		return
	}

	if err := o.validate(); err != nil {
		fmt.Println(err)

		usage()

		return
	}

	cli := utils.NewHttpClient(3)

	excel, err := newExcel(o.fileName)
	if err != nil {
		fmt.Printf("init excel failed, err:%v\n", err)

		return
	}

	h := &reviewCount{
		excel:           excel,
		reviewer:        o.reviewer,
		maxReviewCount:  o.maxReviewCount,
		commentService:  newReviewCommentService(&cli, o.token),
		prSearchService: newPRSearchService(&cli, o.token, o.reviewer),
	}

	h.countComment()
}
