package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run unit tests in ./pkg/unit_test",
	Run: func(cmd *cobra.Command, args []string) {
		runTest()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTest() {
	fmt.Println("Running tests in ./pkg/unit_test...")

	testCmd := exec.Command("go", "test", "-count=1", "./pkg/unit_test", "-v")
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr
	testCmd.Env = os.Environ() // inherit env vars

	if err := testCmd.Run(); err != nil {
		fmt.Println("Test failed:", err)
		os.Exit(1)
	}
}
