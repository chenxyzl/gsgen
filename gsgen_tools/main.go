package main

import (
	"fmt"
	"github.com/chenxyzl/gsgen/gsgen_tools/internal"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gen",
		Short: "gen is a function generate, getter/setter/bson",
		Run: func(cmd *cobra.Command, args []string) {
			//parse dir
			dir, err := cmd.Flags().GetString("dir")
			if err != nil {
				panic(fmt.Sprintf("param parse err: dir, err:%v", err))
			}
			dir = filepath.Clean(dir)

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

			//parse bson
			genBson, err := cmd.Flags().GetBool("bson")
			if err != nil {
				panic(fmt.Sprintf("param parse err: bson, err:%v", err))
			}

			//parse head ext annotation
			headAnnotations, err := cmd.Flags().GetStringSlice("head_annotations")
			if err != nil {
				panic(fmt.Sprintf("param parse err: head_ext_annotation, err:%v", err))
			}

			//
			fmt.Printf("dir: %v\nfile suffix: %v\ngen getter: true[must]\ngen setter: %v\ngen bson:%v \n", dir, fileSuffix, genSetter, genBson)
			//
			internal.Gen(dir, fileSuffix, genSetter, genBson, headAnnotations)
		},
	}
	//增加默认命令
	rootCmd.Flags().StringP("dir", "d", "model", "target dir")
	rootCmd.Flags().StringSliceP("file_suffix", "f", []string{".model.go"}, "target file suffix")
	rootCmd.Flags().BoolP("setter", "s", false, "gen setter")
	rootCmd.Flags().BoolP("bson", "b", false, "gen bson")
	rootCmd.Flags().StringSliceP("head_annotations", "a", []string{}, "head annotations")

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
