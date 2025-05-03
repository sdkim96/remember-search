package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

func generateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func Invoke[T any](
	systemPrompt string,
	userPrompt string,
	client openai.Client,
) (T, error) {

	var result T
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "result",
		Description: openai.String("Structured response"),
		Schema:      generateSchema[T](),
		Strict:      openai.Bool(true),
	}

	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemPrompt),
				openai.UserMessage(userPrompt),
			},
			ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
			},
			Model: openai.ChatModelGPT4oMini,
		},
	)
	if err != nil {
		var empty T
		return empty, fmt.Errorf("openai request error: %w", err)
	}

	// ✨ structured output 결과 꺼내기
	data := resp.Choices[0].Message.Content

	// JSON 파싱
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		var empty T
		return empty, fmt.Errorf("json unmarshal error: %w", err)
	}

	return result, nil
}
