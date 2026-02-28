//go:build !plan9 && !windows
// +build !plan9,!windows

package run

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/go-nunu/nunu/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fsnotify/fsnotify"
	"github.com/go-nunu/nunu/internal/pkg/helper"
	"github.com/spf13/cobra"
)

var quit = make(chan os.Signal, 1)

type Run struct {
}

// 过滤目录  监控后缀  构建参数
var excludeDir string
var includeExt string
var buildFlags string

func init() {
	CmdRun.Flags().StringVarP(&excludeDir, "excludeDir", "", excludeDir, `eg: nunu run --excludeDir="tmp,vendor,.git,.idea"`)
	CmdRun.Flags().StringVarP(&includeExt, "includeExt", "", includeExt, `eg: nunu run --includeExt="go,tpl,tmpl,html,yaml,yml,toml,ini,json"`)
	CmdRun.Flags().StringVarP(&buildFlags, "buildFlags", "", buildFlags, `eg: nunu run --buildFlags="-tags cse"`)
	if excludeDir == "" {
		excludeDir = config.RunExcludeDir
	}
	if includeExt == "" {
		includeExt = config.RunIncludeExt
	}
}

var CmdRun = &cobra.Command{
	Use:     "run",
	Short:   "nunu run [main.go path]",
	Long:    "nunu run [main.go path]",
	Example: "nunu run cmd/server",
	Run: func(cmd *cobra.Command, args []string) {
		cmdArgs, programArgs := helper.SplitArgs(cmd, args)
		var dir string
		if len(cmdArgs) > 0 {
			dir = cmdArgs[0]
		}
		// 获取当前工作目录
		base, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
			return
		}

		// 如果你没有指定目录
		if dir == "" {
			cmdPath, err := helper.FindMain(base, excludeDir)

			if err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
				return
			}
			switch len(cmdPath) {
			case 0:
				fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", "The cmd directory cannot be found in the current directory")
				return
			case 1:
				for _, v := range cmdPath {
					dir = v
				}
			default:
				var cmdPaths []string
				for k := range cmdPath {
					cmdPaths = append(cmdPaths, k)
				}
				sort.Strings(cmdPaths)

				prompt := &survey.Select{
					Message:  "Which directory do you want to run?",
					Options:  cmdPaths,
					PageSize: 10,
				}
				e := survey.AskOne(prompt, &dir)
				if e != nil || dir == "" {
					return
				}
				dir = cmdPath[dir]
			}
		}
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		fmt.Printf("\033[35mNunu run %s.\033[0m\n", dir)
		fmt.Printf("\033[35mWatch excludeDir %s\033[0m\n", excludeDir)
		fmt.Printf("\033[35mWatch includeExt %s\033[0m\n", includeExt)
		fmt.Printf("\033[35mWatch buildFlags %s\033[0m\n", buildFlags)
		watch(dir, programArgs)

	},
}

func watch(dir string, programArgs []string) {

	// Listening file path
	watchPath := "./"

	// 初始化 fsnotify 监控器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer watcher.Close()

	// 处理过滤目录
	excludeDirArr := strings.Split(excludeDir, ",")
	// 要监控的后缀
	includeExtArr := strings.Split(includeExt, ",")
	buildFlagsArr := make([]string, 0)
	if buildFlags = strings.TrimSpace(buildFlags); buildFlags != "" {
		buildFlagsArr = strings.Split(buildFlags, " ")
	}
	// 优化查找性能 , 使得比对该文件 是否是查找后缀的时间复杂度 变成O(1)
	includeExtMap := make(map[string]struct{})
	for _, s := range includeExtArr {
		includeExtMap[s] = struct{}{}
	}

	// 递归遍历当前目录下的文件
	err = filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, s := range excludeDirArr {
			// 如果过滤目录为空, 则跳过
			if s == "" {
				continue
			}
			// 如果路径包含过滤目录, 则跳过
			if strings.HasPrefix(path, s) {
				return nil
			}

		}
		// 如果是文件, 则判断是否是要监控的后缀
		if !info.IsDir() {
			ext := filepath.Ext(info.Name())
			// 去掉. 查找是否是自己要监控的后缀
			if _, ok := includeExtMap[strings.TrimPrefix(ext, ".")]; ok {
				err = watcher.Add(path)
				if err != nil {
					fmt.Println("Error:", err)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 启动程序
	cmd := start(dir, buildFlagsArr, programArgs)

	// 死循环监听
	for {
		select {
		case <-quit:
			// 收到退出信号, 发送SIGINT信号给进程组, 强制结束所有子进程
			err = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)

			if err != nil {
				fmt.Printf("\033[31mserver exiting...\033[0m\n")
				return
			}
			fmt.Printf("\033[31mserver exiting...\033[0m\n")
			os.Exit(0)

		case event := <-watcher.Events:
			// 收到文件修改事件, 发送SIGKILL信号给进程组, 强制结束所有子进程
			// 过滤操作类型：只有 创建(Create)、写入(Write)、删除(Remove) 才触发
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf("\033[36mfile modified: %s\033[0m\n", event.Name)
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

				cmd = start(dir, buildFlagsArr, programArgs)
			}
		case err := <-watcher.Errors:
			fmt.Println("Error:", err)
		}
	}
}

// 构造并执行go run 命令, 启动程序
func start(dir string, buildFlagsArgs []string, programArgs []string) *exec.Cmd {

	// 1. 构造命令参数切片
	// 最终形态类似: go run -race ./cmd/server -conf=config.yaml
	run := []string{"run"}
	// 如-race
	run = append(run, buildFlagsArgs...)

	// 如 ./cmd/server
	run = append(run, dir)

	// 创建命令
	cmd := exec.Command("go", append(run, programArgs...)...)

	// 	创建一个进程组, 往后go run命令创建的子进程(server)也会属于这个进程组
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatalf("\033[33;1mcmd run failed\u001B[0m")
	}
	time.Sleep(time.Second)
	fmt.Printf("\033[32;1mrunning...\033[0m\n")
	return cmd

}
