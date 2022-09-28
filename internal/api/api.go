package api

import (
	"context"
	"fmt"
	"time"

	"github.com/antgrutta/gh-discussions/internal/log"
	"github.com/diggs/go-backoff"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type queryVariables interface {
	githubv4.CreateDiscussionInput | githubv4.AddDiscussionCommentInput | githubv4.MarkDiscussionCommentAsAnswerInput
}

func newGHClient() *githubv4.Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: fmt.Sprint(viper.Get("IMPORT-TOKEN"))})
	httpClient := oauth2.NewClient(context.Background(), src)

	return githubv4.NewClient(httpClient)
}

func Query(query interface{}, variables map[string]interface{}) error {
	client := newGHClient()

	// Back off between 0 and exponentially, starting at 30 seconds, capping at 10 minutes
	exp := backoff.NewExponentialFullJitter(30*time.Second, 10*time.Minute)

	// Error to be returned
	var err error

	for {
		err := client.Query(context.Background(), query, variables)
		if err != nil {
			fmt.Println(err)
			exp.Backoff()
		} else {
			break
		}
	}
	return err
}

func Mutation[T queryVariables](mutation interface{}, variables T) error {
	client := newGHClient()
	logger := log.NewLogger()

	// Back off between 0 and exponentially, starting at 30 seconds, capping at 10 minutes
	exp := backoff.NewExponentialFullJitter(30*time.Second, 10*time.Minute)

	// Error to be returned
	var err error

	for {
		err = client.Mutate(context.Background(), mutation, variables, nil)
		if err != nil {
			logger.Println(err)
			exp.Backoff()
		} else {
			break
		}
	}
	return err
}
