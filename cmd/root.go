package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	server "github.com/hyloblog/hyloblog/internal/app"
	"github.com/hyloblog/hyloblog/internal/assert"
	"github.com/hyloblog/hyloblog/internal/config"
	"github.com/hyloblog/hyloblog/internal/dns"
	"github.com/hyloblog/hyloblog/internal/email/emailqueue"
	"github.com/hyloblog/hyloblog/internal/httpclient"
	"github.com/hyloblog/hyloblog/internal/model"
	"github.com/spf13/cobra"
)

const clientTimeout = 30 * time.Second

var rootCmd = &cobra.Command{
	Use:   "hyloblog",
	Short: "Run hyloblog",
	RunE: func(cmd *cobra.Command, args []string) error {
		rand.Seed(time.Now().UnixNano())
		db, err := config.Config.Db.Connect()
		if err != nil {
			return fmt.Errorf("could not connect to db: %w", err)
		}
		c := httpclient.NewHttpClient(clientTimeout)
		store := model.NewStore(db)
		if err := reserveSubdomains(
			getReservedSubdomains(),
			store,
		); err != nil {
			return fmt.Errorf("reserve subdomains: %w", err)
		}
		go func() {
			if err := emailqueue.Run(c, store); err != nil {
				log.Fatal("email queue error", err)
			}
		}()
		return server.Serve(c, store)
	},
}

func getReservedSubdomains() []string {
	var domains []string
	for _, sub := range config.Config.ReservedSubdomains {
		domains = append(domains, sub)
	}
	for _, rule := range config.Config.RedirectRules {
		parts := strings.Split(
			rule.From,
			fmt.Sprintf(".%s", config.Config.Hyloblog.RootDomain),
		)
		switch len(parts) {
		case 1: /* not a subdomain */
			continue
		case 2: /* subdomain */
			domains = append(domains, parts[0])
		default: /* double occurrence of root */
			assert.Assert(false)
		}
	}
	return domains
}

func reserveSubdomains(domains []string, s *model.Store) error {
	if err := s.DeleteReservedSubdomains(context.TODO()); err != nil {
		return fmt.Errorf("delete reserved subdomains: %w", err)
	}
	for _, rawsub := range domains {
		sub, err := dns.ParseSubdomain(rawsub)
		if err != nil {
			return fmt.Errorf("parse subdomain: %w", err)
		}
		if err := s.ReserveSubdomain(context.TODO(), sub); err != nil {
			return fmt.Errorf("reserve subdomain: %w", err)
		}
	}
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
