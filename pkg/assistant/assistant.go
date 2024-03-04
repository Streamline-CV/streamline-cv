package assistant

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Streamline-CV/streamline-cv/api"
	"github.com/lithammer/dedent"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"strings"
)

type AiSuggestion struct {
	SuggestedText string       `json:"suggestedText"`
	Reasoning     string       `json:"reasoning"`
	Severity      api.Severity `json:"severity"`
}

type AiAssistant struct {
	client *openai.Client
}

func NewAiAssistant(openaiKey string) (*AiAssistant, error) {
	client := openai.NewClient(openaiKey)
	return &AiAssistant{
		client: client,
	}, nil
}

func (a *AiAssistant) Refactor(changeReport *api.ChangeReport) (*api.SuggestionReporting, error) {

	var suggestions []api.Suggestion

	for key, change := range changeReport.Changes {
		completion, err := a.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: generatePrompt(change),
				},
			},
		})
		if err != nil {
			log.Info().Msgf("Error getting completion for key %s: %v", key, err)
			continue
		}
		if len(completion.Choices) == 0 {
			log.Error().Msgf("Error while trying to get completions.")
			continue
		}
		suggestion := completion.Choices[0].Message.Content
		var aiSuggestion AiSuggestion
		err = json.Unmarshal([]byte(suggestion), &aiSuggestion)
		if err != nil {
			log.Error().Msgf("Failed parsing completion response %s", err)
		}
		log.Info().Msgf("Suggested text: %s with reasoning %s", aiSuggestion.SuggestedText, aiSuggestion.Reasoning)
		if aiSuggestion.SuggestedText != "" {
			suggestions = append(suggestions, api.Suggestion{
				Path:        change.Path,
				Line:        change.Target.Line,
				ColumnStart: change.Target.ColumnStart,
				ColumnEnd:   change.Target.ColumnEnd,
				Comment:     aiSuggestion.Reasoning,
				Value:       aiSuggestion.SuggestedText,
				Severity:    aiSuggestion.Severity,
			})
		}
	}

	return &api.SuggestionReporting{
		Suggestions: suggestions,
	}, nil
}

func generatePrompt(change api.Change) string {
	prompt := dedent.Dedent(`
	As a professional HR and hiring manager help to improve CV.\n
	CV section is "%s" and text is "%s".\n
	Please adapt this text for a resume, while also following these rules:\n
	– Avoid advanced and rarely used words.\n
	– At the same time, don't use informal language, it must look professional.\n
	– Avoid overly complex sentences but don't make them too short.\n
	- If the text is already short and clear fix just grammar issues.\n
	– Rephrase and reformulate but fully keep the original meaning and don’t leave anything out.\n
	Use the following strict json response (just json, no code markup) represented by this example:\n
	{
		"suggestedText": "Suggested text value",
		"reasoning": "I suggested updating text value because it sounds more professional",
        "severity": "Info" // severity level of the suggestion, can be Info,Warning,Error
	}\n
	If the is nothing to improve, leave suggestedText empty.\n
	Always fill "reasoning" field with the explanation for "suggestedText" or why no improvement is necessary.\n
	`)
	return fmt.Sprintf(prompt, strings.Join(change.Path, "->"), change.Target)
}
