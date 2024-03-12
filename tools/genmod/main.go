package main

import (
	"fmt"
	"gotest/tools/genmod/internal"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "genmod",
		Short: "genmod is a model function generate",
		Run: func(cmd *cobra.Command, args []string) {
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				panic("need param: dir")
			}
			fmt.Printf("gg:%s\n", Version)

			currentDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			fullDir := filepath.Clean(filepath.Join(currentDir, dir))
			fmt.Printf("dir:%v|currentDir:%v|fullDir:%v\n", dir, currentDir, fullDir)

			internal.Gen(dir)
		},
	}
	rootCmd.Flags().StringP("dir", "d", "model", "target dir")

	// 添加命令
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// greetCmd 创建 greet 命令
func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "exec file version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version:%s\n", Version)
		},
	}

	cmd.Flags().String("name", "World", "Name of the person to greet")

	return cmd
}
