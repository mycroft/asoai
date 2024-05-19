package commands

import (
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

func NewModelsCommand() *cobra.Command {
	modelsCommand := cobra.Command{
		Use:   "models",
		Short: "list models",
		Long:  "list all available models exposed by the API",
		Run: func(cmd *cobra.Command, args []string) {
			listModels()
		},
	}

	return &modelsCommand
}

func listModels() {
	envVar := os.Getenv("OPENAI_API_KEY")
	if envVar == "" {
		fmt.Printf("could not find OPENAI_API_KEY")
		os.Exit(1)
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	models, err := client.ListModels(context.Background())
	if err != nil {
		fmt.Printf("could not list models: %v\n", err)
		os.Exit(1)
	}

	modelsList := []string{}

	for _, model := range models.Models {
		modelsList = append(modelsList, model.ID)
	}

	sort.Strings(modelsList)

	for _, model := range modelsList {
		fmt.Println(model)
	}
}
