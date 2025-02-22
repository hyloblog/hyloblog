package blog

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/hyloblog/hyloblog/internal/app/handler/request"
	"github.com/hyloblog/hyloblog/internal/app/handler/response"
	"github.com/hyloblog/hyloblog/internal/model"
	"github.com/hyloblog/hyloblog/internal/session"
	"github.com/hyloblog/hyloblog/internal/util"
)

type SubscriberData struct {
	Count            int
	CumulativeCounts template.JS
	Subscribers      []Subscriber
}

type CumulativeCount struct {
	Timestamp string `json:"timestamp"` /* Date as a string in "YYYY-MM-DD" format */
	Count     int    `json:"count"`     /* Cumulative subscriber count on that date */
}

type Subscriber struct {
	Email        string
	SubscribedOn string
	Status       string
}

func (b *BlogService) SubscriberMetrics(
	r request.Request,
) (response.Response, error) {
	sesh := r.Session()
	sesh.Println("SubscriberMetrics handler...")

	r.MixpanelTrack("SubscriberMetrics")

	blogID, ok := r.GetRouteVar("blogID")
	if !ok {
		return nil, createCustomError("", http.StatusNotFound)
	}

	data, err := b.subscriberMetrics(blogID)
	if err != nil {
		return nil, fmt.Errorf("subscriber metrics: %w", err)
	}
	return response.NewTemplate(
		[]string{"subscriber_metrics.html", "subscribers.html"},
		util.PageInfo{
			Data: struct {
				Title          string
				UserInfo       *session.UserInfo
				SubscriberData SubscriberData
			}{
				Title:          "Dashboard",
				UserInfo:       session.ConvertSessionToUserInfo(sesh),
				SubscriberData: data,
			},
		},
	), nil
}

func (b *BlogService) subscriberMetrics(blogID string) (SubscriberData, error) {
	subs, err := b.store.ListActiveSubscribersByBlogID(
		context.TODO(), blogID,
	)
	if err != nil {
		return SubscriberData{}, fmt.Errorf(
			"list active subscriber: %w", err,
		)
	}

	jsonSubscriberData, err := json.Marshal(
		buildSubscriberCumulativeCounts(subs),
	)
	if err != nil {
		return SubscriberData{}, fmt.Errorf(
			"json marshall: %w", err,
		)
	}
	return SubscriberData{
		Count:            len(subs),
		CumulativeCounts: template.JS(jsonSubscriberData),
		Subscribers:      convertSubscribers(subs),
	}, nil
}

func convertSubscribers(subs []model.Subscriber) []Subscriber {
	var res []Subscriber
	for _, s := range subs {
		res = append(res, Subscriber{
			Email:        s.Email,
			SubscribedOn: s.CreatedAt.Format("January 2, 2006"),
			Status:       string(s.Status),
		})
	}
	return res
}

func buildSubscriberCumulativeCounts(subs []model.Subscriber) []CumulativeCount {
	/* hold cumulative counts per hour */
	cumulativeCounts := make(map[time.Time]int)

	for _, sub := range subs {
		hour := sub.CreatedAt.Truncate(time.Hour)
		cumulativeCounts[hour]++
	}
	var result []CumulativeCount
	cumulativeSum := 0
	for hour := range cumulativeCounts {
		cumulativeSum += cumulativeCounts[hour]
		result = append(result, CumulativeCount{
			Timestamp: hour.Format(time.RFC1123Z), /* ISO format for browser compatibility */
			Count:     cumulativeSum,
		})
	}
	return result
}
