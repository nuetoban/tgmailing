package view

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type Prometheus struct {
	s StatisticsProvider
}

type labels map[string]interface{}

func (l labels) String() string {
	var (
		out    strings.Builder
		fields []string
	)

	out.WriteRune('{')

	// Sort labels
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fields = append(fields, fmt.Sprintf("%s=%v", k, l[k]))
	}

	out.WriteString(strings.Join(fields, ","))
	out.WriteRune('}')
	return out.String()
}

type metric struct {
	Value  interface{}
	Labels labels
}

func buildPromString(input map[string][]metric) string {
	var out strings.Builder

	for k, v := range input {
		for _, m := range v {
			out.WriteString(fmt.Sprintf("%s%s %v\n", k, m.Labels, m.Value))
		}
	}

	return out.String()
}

func NewPrometheus(s StatisticsProvider) *Prometheus {
	return &Prometheus{s: s}
}

func (p *Prometheus) StatisticsHTTP(w http.ResponseWriter, req *http.Request) {
	s := p.s.Statistics()

	label := make(labels)
	label["mailing_id"] = s.ID

	prom := make(map[string][]metric)
	prom["start_time"] = []metric{{Labels: label, Value: s.StartTime.Unix()}}

	for _, v := range s.Statistics {
		botSpecificLabels := make(labels)
		botSpecificLabels["mailing_id"] = s.ID
		botSpecificLabels["bot_id"] = v.BotID
		botSpecificLabels["bot_name"] = v.BotName

		if _, ok := prom["successful_send_attemts"]; !ok {
			prom["successful_send_attemts"] = []metric{}
		}

		prom["successful_send_attemts"] = append(prom["successful_send_attemts"], metric{Labels: botSpecificLabels, Value: v.SuccessfulSendAttempts})
		prom["failed_send_attemts"] = append(prom["failed_send_attemts"], metric{Labels: botSpecificLabels, Value: v.FailedSendAttempts})
	}

	w.Write([]byte(buildPromString(prom)))
}
