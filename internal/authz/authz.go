package authz

import (
	"context"
	"fmt"

	"github.com/hyloblog/hyloblog/internal/assert"
	"github.com/hyloblog/hyloblog/internal/authz/digitalsize"
	"github.com/hyloblog/hyloblog/internal/authz/internal/option"
	"github.com/hyloblog/hyloblog/internal/authz/internal/trafficsize"
	"github.com/hyloblog/hyloblog/internal/model"
)

func CanCreateSite(s *model.Store, userid string) (bool, error) {
	storageUsed, err := UserStorageUsed(s, userid)
	if err != nil {
		return false, fmt.Errorf("calculate user storage used: %w", err)
	}
	/* get user's site count */
	blogCount, err := s.CountLiveBlogsByUserID(context.TODO(), userid)
	if err != nil {
		return false, fmt.Errorf("get user project count: %w", err)
	}
	/* get user's tier features */
	plan, err := s.GetUserSubscriptionByID(context.TODO(), userid)
	if err != nil {
		return false, fmt.Errorf("get user subscription: %w", err)
	}
	tier := subscriptionTiers[plan]
	return tier.canCreateProject(storageUsed, int(blogCount)), nil
}

func HasAnalyticsCustomDomainsImagesEmails(
	s *model.Store, userid string,
) (bool, error) {
	plan, err := s.GetUserSubscriptionByID(context.TODO(), userid)
	if err != nil {
		return false, fmt.Errorf("get user subscription: %w", err)
	}
	tier := subscriptionTiers[plan]
	return tier.analyticsCustomDomainImagesEmails.Value(), nil
}

func CanUseTheme(s *model.Store, theme string, userid string) (bool, error) {
	plan, err := s.GetUserSubscriptionByID(context.TODO(), userid)
	if err != nil {
		return false, fmt.Errorf("get user subscription: %w", err)
	}
	tier := subscriptionTiers[plan]
	return tier.canUseTheme(theme), nil
}

type Tier interface {
	Name() string
}

func GetTiers() []Tier {
	return []Tier{
		subscriptionTiers[model.SubNameBasic],
		subscriptionTiers[model.SubNamePremium],
	}
}

type subscriptionTier struct {
	name                              string
	projects                          int
	storageSize                       digitalsize.Size
	visitorsPerMonth                  trafficsize.Size
	themes                            []string
	codeStyles                        []string
	analyticsCustomDomainImagesEmails option.Option
}

func (tier subscriptionTier) Name() string { return tier.name }

var subscriptionTiers = map[model.SubName]subscriptionTier{
	model.SubNameBasic: {
		name:                              "basic",
		projects:                          1,
		storageSize:                       32 * digitalsize.Megabyte,
		visitorsPerMonth:                  10000,
		themes:                            []string{"lit"},
		codeStyles:                        []string{"lit"},
		analyticsCustomDomainImagesEmails: option.New(false),
	},
	model.SubNamePremium: {
		name:                              "premium",
		projects:                          10,
		storageSize:                       digitalsize.Gigabyte,
		visitorsPerMonth:                  100000,
		themes:                            []string{"lit", "latex"},
		codeStyles:                        []string{"lit", "latex"},
		analyticsCustomDomainImagesEmails: option.New(true),
	},
}

type Feature interface {
	Name() string
	Value(Tier) string
}

func GetFeatures() []Feature {
	return []Feature{
		featureProjects,
		featureStorage,
		featureVisitorsPerMonth,
		featureCustomDomain,
		featureEmailSubscribers,
		featureAnalytics,
	}
}

type feature int

const (
	featureProjects feature = iota
	featureStorage
	featureVisitorsPerMonth
	featureCustomDomain
	featureEmailSubscribers
	featureAnalytics
)

func (f feature) Name() string {
	switch f {
	case featureProjects:
		return "Projects"
	case featureStorage:
		return "Storage"
	case featureVisitorsPerMonth:
		return "Visitors per month"
	case featureCustomDomain:
		return "Custom domain"
	case featureEmailSubscribers:
		return "Email subscribers"
	case featureAnalytics:
		return "Analytics"
	default:
		assert.Assert(false)
		return ""
	}
}

func (f feature) Value(rawtier Tier) string {
	tier, ok := rawtier.(subscriptionTier)
	assert.Assert(ok)
	switch f {
	case featureProjects:
		return fmt.Sprintf("%d", tier.projects)
	case featureStorage:
		return tier.storageSize.Abbrev(0)
	case featureVisitorsPerMonth:
		return tier.visitorsPerMonth.Abbrev(0)
	case featureCustomDomain,
		featureEmailSubscribers,
		featureAnalytics:
		return tier.analyticsCustomDomainImagesEmails.String()
	default:
		assert.Assert(false)
		return ""
	}
}

func (s *subscriptionTier) canCreateProject(
	sizeUsed digitalsize.Size, blogCount int,
) bool {
	return blogCount < s.projects && sizeUsed < s.storageSize
}

func (s *subscriptionTier) canUseTheme(theme string) bool {
	for _, t := range s.themes {
		if t == theme {
			return true
		}
	}
	return false
}
