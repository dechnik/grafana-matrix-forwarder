package grafana

import (
	"bytes"
	"grafana-matrix-forwarder/cfg"
	"grafana-matrix-forwarder/matrix"
	"html/template"
	"log"
)

const (
	alertMessageStr     = `💔 ️<b>ALERT</b><p>Rule: <a href="{{ .RuleURL }}">{{ .RuleName }}</a> | {{ .Message }}</p>`
	resolvedMessageStr  = `💚 ️<b>RESOLVED</b><p>Rule: <a href="{{ .RuleURL }}">{{ .RuleName }}</a> | {{ .Message }}</p>`
	noDataMessageStr    = `❓️<b>NO DATA</b><ul><p>Rule: <a href="{{ .RuleURL }}">{{ .RuleName }}</a> | {{ .Message }}</p>`
	unknownMessageStr   = `❓️<b>UNKNOWN</b><ul><li>Rule: <a href="{{ .RuleURL }}">{{ .RuleName }}</a> | {{ .Message }}</li><li>State: <b>{{ .State }}</b></li></ul>`
	resolvedReactionStr = `✅`
)

var (
	alertMessageTemplate    = template.Must(template.New("alertMessage").Parse(alertMessageStr))
	resolvedMessageTemplate = template.Must(template.New("resolvedMessage").Parse(resolvedMessageStr))
	noDataMessageTemplate   = template.Must(template.New("noDataMessage").Parse(noDataMessageStr))
	unknownMessageTemplate  = template.Must(template.New("unknownMessage").Parse(unknownMessageStr))

	eventIDMap = map[string]string{}
)

// ForwardAlert sends the provided grafana.AlertPayload to the provided matrix.Writer using the provided roomID
func ForwardAlert(writer matrix.Writer, roomID string, alert AlertPayload, resolveMode cfg.ResolveMode) (err error) {
	resolveWithReaction := resolveMode == cfg.ResolveWithReaction

	alertID := alert.FullRuleID()
	if eventID, ok := eventIDMap[alertID]; ok {
		if alert.State == AlertStateResolved && resolveWithReaction {
			delete(eventIDMap, alertID)
			return sendReaction(writer, roomID, eventID)
		}
	}
	return sendRegularMessage(writer, roomID, alert, alertID)
}

func sendReaction(writer matrix.Writer, roomID string, eventID string) (err error) {
	_, err = writer.React(roomID, eventID, resolvedReactionStr)
	return
}

func sendRegularMessage(writer matrix.Writer, roomID string, alert AlertPayload, alertID string) (err error) {
	formattedMessageBody, err := buildFormattedMessageBodyFromAlert(alert)
	if err != nil {
		return
	}
	formattedMessage := matrix.NewSimpleFormattedMessage(formattedMessageBody)
	response, err := writer.Send(roomID, formattedMessage)
	if err == nil {
		eventIDMap[alertID] = response.EventID.String()
	}
	return
}

func buildFormattedMessageBodyFromAlert(alert AlertPayload) (message string, err error) {
	switch alert.State {
	case AlertStateAlerting:
		message, err = executeTemplate(alertMessageTemplate, alert)
	case AlertStateResolved:
		message, err = executeTemplate(resolvedMessageTemplate, alert)
	case AlertStateNoData:
		message, err = executeTemplate(noDataMessageTemplate, alert)
	default:
		log.Printf("alert received with unknown state: %s", alert.State)
		message, err = executeTemplate(unknownMessageTemplate, alert)
	}
	return message, err
}

func executeTemplate(template *template.Template, alert AlertPayload) (string, error) {
	buffer := new(bytes.Buffer)
	err := template.Execute(buffer, alert)
	return buffer.String(), err
}
