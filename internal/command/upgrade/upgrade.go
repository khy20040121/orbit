package upgrade

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/khy20040121/orbit/config"
	"github.com/spf13/cobra"
)

var CmdUpgrade = &cobra.Command{
	Use:     "upgrade",
	Short:   "Upgrade the orbit command.",
	Long:    "Upgrade the orbit command.",
	Example: "orbit upgrade",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("go install %s\n", config.OrbitCmd)
		cmd := exec.Command("go", "install", config.OrbitCmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("go install %s error\n", err)
		}
		fmt.Printf("\n🎉 orbit upgrade successfully!\n\n")
	},
}
