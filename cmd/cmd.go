package cmd

import (
	"fmt"
	"github.com/gvcgo/fcode/cnf"
	"github.com/gvcgo/fcode/models"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {
	rootCmd = &cobra.Command{}

	rootCmd.AddCommand(list)
	rootCmd.AddCommand(use)
	rootCmd.AddCommand(serve)
	rootCmd.AddCommand(showPath)
	rootCmd.AddCommand(stop)

}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

var list = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List available model names.",
	Long:    "List available model names.",
	Example: "fcode list",
	Run: func(cmd *cobra.Command, args []string) {
		aiModels := cnf.DefaultConf.AIModels
		if len(aiModels) == 0 {
			fmt.Println("No model is available!")
			return
		}
		for _, m := range aiModels {
			if m.Name == cnf.DefaultModel.Name {
				fmt.Println("**", m.Name)
			} else {
				fmt.Println(m.Name)
			}
		}
	},
}

var use = &cobra.Command{
	Use:     "use",
	Aliases: []string{"u"},
	Short:   "Use an available model.",
	Long:    "Use an available model.",
	Example: "fcode use qwen",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		modelName := args[0]
		url := fmt.Sprintf("http://localhost%s/v1/choose/model?name=%s", cnf.DefaultConf.GetPort(), modelName)
		resp, err := http.Post(url, "application/json", nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			content, _ := io.ReadAll(resp.Body)
			fmt.Println(string(content))
		}
	},
}

var serve = &cobra.Command{
	Use:     "serve",
	Aliases: []string{"s"},
	Short:   "Run server for lsp-ai.",
	Long:    "Run server for lsp-ai.",
	Example: "fcode serve",
	Run: func(cmd *cobra.Command, args []string) {
		models.Serve()
	},
}

var showPath = &cobra.Command{
	Use:     "show",
	Aliases: []string{"sh"},
	Short:   "Show config file path.",
	Long:    "Show config file path.",
	Example: "fcode show",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cnf.DefaultConf.GetPath())
	},
}

var stop = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"st"},
	Short:   "Stop fcode server.",
	Long:    "Stop fcode server.",
	Example: "fcode stop",
	Run: func(cmd *cobra.Command, args []string) {
		url := fmt.Sprintf("http://localhost%s/v1/stop", cnf.DefaultConf.GetPort())
		http.Post(url, "application/json", nil)
	},
}
