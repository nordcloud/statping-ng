package notifiers

import (
	"strings"
	"time"

	"github.com/nordcloud/statping-ng/types/null"

	"github.com/nordcloud/statping-ng/types/failures"
	"github.com/nordcloud/statping-ng/types/notifications"
	"github.com/nordcloud/statping-ng/types/notifier"
	"github.com/nordcloud/statping-ng/types/services"
	"github.com/nordcloud/statping-ng/utils"
)

var _ notifier.Notifier = (*gotify)(nil)

type gotify struct {
	*notifications.Notification
}

func (g *gotify) Select() *notifications.Notification {
	return g.Notification
}

func (g *gotify) Valid(values notifications.Values) error {
	return nil
}

var Gotify = &gotify{&notifications.Notification{
	Method:      "gotify",
	Title:       "Gotify",
	Description: "Use Gotify to receive push notifications. Add your Gotify URL and App Token to receive notifications.",
	Author:      "Hugo van Rijswijk",
	AuthorUrl:   "https://github.com/hugo-vrijswijk",
	Icon:        "broadcast-tower",
	Delay:       time.Duration(5 * time.Second),
	Limits:      60,
	SuccessData: null.NewNullString(`{"title": "{{.Service.Name}}", "message": "Your service '{{.Service.Name}}' is currently online!", "priority": 2}`),
	FailureData: null.NewNullString(`{"title": "{{.Service.Name}}", "message": "Your service '{{.Service.Name}}' is currently failing! Reason: {{.Failure.Issue}}", "priority": 5}`),
	DataType:    "json",
	Form: []notifications.NotificationForm{{
		Type:        "text",
		Title:       "Gotify URL",
		SmallText:   "Gotify server URL, including http(s):// and port if needed",
		DbField:     "Host",
		Placeholder: "https://gotify.domain.com",
		Required:    true,
	}, {
		Type:        "text",
		Title:       "App Token",
		SmallText:   "The Application Token generated by Gotify",
		DbField:     "api_key",
		Placeholder: "TB5gatYYyR.FCD2",
		Required:    true,
	}}},
}

// Send will send a HTTP Post to the Gotify API. It accepts type: string
func (g *gotify) sendMessage(msg string) (string, error) {
	var url string
	if strings.HasSuffix(g.Host.String, "/") {
		url = g.Host.String + "message"
	} else {
		url = g.Host.String + "/message"
	}

	headers := []string{"X-Gotify-Key=" + g.ApiKey.String}

	content, _, err := utils.HttpRequest(url, "POST", "application/json", headers, strings.NewReader(msg), time.Duration(10*time.Second), true, nil)

	return string(content), err
}

// OnFailure will trigger failing service
func (g *gotify) OnFailure(s services.Service, f failures.Failure) (string, error) {
	out, err := g.sendMessage(ReplaceVars(g.FailureData.String, s, f))
	return out, err
}

// OnSuccess will trigger successful service
func (g *gotify) OnSuccess(s services.Service) (string, error) {
	out, err := g.sendMessage(ReplaceVars(g.SuccessData.String, s, failures.Failure{}))
	return out, err
}

// OnTest will test the Gotify notifier
func (g *gotify) OnTest() (string, error) {
	msg := `{"title:" "Test" "message": "Testing the Gotify Notifier", "priority": 0}`
	content, err := g.sendMessage(msg)

	return content, err
}

// OnSave will trigger when this notifier is saved
func (g *gotify) OnSave() (string, error) {
	return "", nil
}
