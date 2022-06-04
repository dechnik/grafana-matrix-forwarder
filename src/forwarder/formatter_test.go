package forwarder

import (
	"grafana-matrix-forwarder/cfg"
	"grafana-matrix-forwarder/model"
	"testing"
)

func Test_buildFormattedMessageBodyFromAlert(t *testing.T) {
	type args struct {
		alert model.Data
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "alertingStateTest",
			args: args{model.Data{
				State:    "alerting",
				RuleURL:  "http://example.com",
				RuleName: "sample",
				Message:  "sample message",
			}},
			want: "💔 <b>ALERT</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
		},
		{
			name: "alertingStateWithEvalMatchesTest",
			args: args{model.Data{
				State:    "alerting",
				RuleURL:  "http://example.com",
				RuleName: "sample",
				Message:  "sample message",
				EvalMatches: []struct {
					Value  float64
					Metric string
					Tags   map[string]string
				}{
					{
						Value:  10.65124,
						Metric: "sample",
						Tags:   map[string]string{},
					},
				},
			}},
			want: "💔 <b>ALERT</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p><ul><li><b>sample</b>: 10.65124</li></ul>",
		},
		{
			name: "alertingStateWithEvalMatchesAndTagsTest",
			args: args{model.Data{
				State:    "alerting",
				RuleURL:  "http://example.com",
				RuleName: "sample",
				Message:  "sample message",
				EvalMatches: []struct {
					Value  float64
					Metric string
					Tags   map[string]string
				}{
					{
						Value:  10.65124,
						Metric: "sample",
					},
				},
				Tags: map[string]string{"key1": "value1", "key2": "value2"},
			}},
			want: "💔 <b>ALERT</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p><ul><li><b>sample</b>: 10.65124</li></ul><p>Tags:</p><ul><li><b>key1</b>: value1</li><li><b>key2</b>: value2</li></ul>",
		},
		{
			name: "okStateTest",
			args: args{model.Data{
				State:    "ok",
				RuleURL:  "http://example.com",
				RuleName: "sample",
				Message:  "sample message",
			}},
			want: "💚 <b>RESOLVED</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
		},
		{
			name: "noDataStateTest",
			args: args{model.Data{
				State:    "no_data",
				RuleURL:  "http://example.com",
				RuleName: "sample",
				Message:  "sample message",
			}},
			want: "❓ <b>NO DATA</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
		},
		{
			name: "unknownStateTest",
			args: args{model.Data{
				State:    "invalid state",
				RuleURL:  "http://example.com",
				RuleName: "sample",
				Message:  "sample message",
			}},
			want: "❓ <b>UNKNOWN</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settings := cfg.AppSettings{MetricRounding: -1}
			got, err := buildFormattedMessageBodyFromAlert(tt.args.alert, settings)
			if err != nil {
				t.Errorf("buildFormattedMessageBodyFromAlert() error: %v", err)
			}
			if got != tt.want {
				t.Errorf("buildFormattedMessageBodyFromAlert() = %v, want %v", got, tt.want)
			}
		})
	}
}
