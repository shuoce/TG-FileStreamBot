package main

import (
	"EverythingSuckz/fsb/config"
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
)

const versionString = "3.2.0"

var rootCmd = &cobra.Command{
	Use:               "fsb [command]",
	Short:             "Telegram 文件直链机器人",
	Long:              "Telegram 机器人，用于为 Telegram 媒体文件生成可直接播放/下载的链接。",
	Example:           "fsb run --port 8080",
	Version:           versionString,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	config.SetFlagsFromConfig(runCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(sessionCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintf(`Telegram 文件直链机器人 版本 %s`, versionString))
}

func printDiskInfo() {
	var stat syscall.Statfs_t

	if err := syscall.Statfs("/", &stat); err != nil {
		fmt.Println("无法获取磁盘信息:", err)
		return
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bavail * uint64(stat.Bsize)
	used := total - stat.Bfree*uint64(stat.Bsize)

	fmt.Printf("========== 磁盘信息 ==========\n")
	fmt.Printf("总容量: %.2f GB\n", float64(total)/1024/1024/1024)
	fmt.Printf("已使用: %.2f GB\n", float64(used)/1024/1024/1024)
	fmt.Printf("可用空间: %.2f GB\n", float64(free)/1024/1024/1024)
	fmt.Printf("==============================\n")
}

func main() {
	printDiskInfo()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
