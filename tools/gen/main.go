package main

import (
	"fmt"
	"gen/internal"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gen",
		Short: "gen is a function generate, getter/setter/mongo",
		Run: func(cmd *cobra.Command, args []string) {
			//parse dir
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				panic(fmt.Sprintf("param parse err: dir, err:%v", err))
			}
			currentDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			fullDir := filepath.Clean(filepath.Join(currentDir, dir))

			//parse setter
			fileSuffix, err := cmd.Flags().GetStringSlice("file_suffix")
			if err != nil {
				panic(fmt.Sprintf("param parse err: file_suffix, err:%v", err))
			}

			//parse setter
			genSetter, err := cmd.Flags().GetBool("setter")
			if err != nil {
				panic(fmt.Sprintf("param parse err: setter, err:%v", err))
			}

			//parse mongo
			genMongo, err := cmd.Flags().GetBool("mongo")
			if err != nil {
				panic(fmt.Sprintf("param parse err: mongo, err:%v", err))
			}

			//
			fmt.Printf("dir: %v\nfile suffix: %v\ngen getter: true[must]\ngen setter: %v\ngen mongo:%v \n", fullDir, fileSuffix, genSetter, genMongo)
			//
			internal.Gen(fullDir, fileSuffix, genSetter, genMongo)
		},
	}
	//增加默认命令
	rootCmd.Flags().StringP("dir", "d", "model", "target dir")
	rootCmd.Flags().StringSliceP("file_suffix", "f", []string{".model.go"}, "target file suffix")
	rootCmd.Flags().BoolP("setter", "s", false, "gen setter")
	rootCmd.Flags().BoolP("mongo", "m", false, "gen mongo")

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
			fmt.Printf("version:%s\n", internal.Version)
		},
	}

	cmd.Flags().String("name", "World", "Name of the person to greet")

	return cmd
}
