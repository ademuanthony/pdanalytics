// Copyright (c) 2018-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package commstats

import (
	"context"
	"net/http"
	"time"

	"github.com/planetdecred/pdanalytics/app"
	"github.com/planetdecred/pdanalytics/web"
)

const (
	dateMiliTemplate = "2006-01-02 15:04:05.99"
	retryLimit       = 3
)

func Activate(ctx context.Context, store DataStore, server *web.Server, options *CommunityStatOptions) error {
	c := &Collector{
		client:    http.Client{Timeout: 10 * time.Second},
		dataStore: store,
		options:   options,
		server:    server,
	}

	if options.CommunityStat {
		go func() {
			c.Run(ctx)
		}()
	}

	if options.CommunityStatHttp {
		if err := c.setupServer(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Collector) Run(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	// continually check the state of the app until its free to run this module
	app.MarkBusyIfFree()

	log.Info("Fetching community stats...")

	app.ReleaseForNewModule()

	go c.startTwitterCollector(ctx)

	go c.startYoutubeCollector(ctx)

	// github
	go c.startGithubCollector(ctx)

	go c.startRedditCollector(ctx)
}

func SetAccounts(options CommunityStatOptions) {
	subreddits = options.Subreddit
	twitterHandles = options.TwitterHandles
	repositories = options.GithubRepositories
	youtubeChannels = options.YoutubeChannelName
}

func (c *Collector) setupServer() error {
	return nil
}
