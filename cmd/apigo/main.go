package main

import (
	"flag"
	"fmt"
	"os"

	"ApiGO/internal"
	"ApiGO/runner"
)

func main() {
	// 命令行参数
	configDir := flag.String("config", "config", "配置文件目录")
	testsDir := flag.String("tests", "testcases", "测试用例目录")
	reportsDir := flag.String("reports", "reports", "报告输出目录")
	project := flag.String("project", "", "指定要运行的项目（多个项目用逗号分隔）")
	format := flag.String("format", "console", "报告格式（console|allure）")
	includeTags := flag.String("include-tags", "", "包含的测试标签（用逗号分隔）")
	excludeTags := flag.String("exclude-tags", "", "排除的测试标签（用逗号分隔）")
	minPriority := flag.Int("min-priority", 0, "最小优先级（0=Low, 1=Normal, 2=High）")

	flag.Parse()

	// 解析项目列表
	activeProjects := []string{}
	if *project != "" {
		activeProjects = strings.Split(*project, ",")
	}

	// 解析标签
	includeTagList := []string{}
	if *includeTags != "" {
		includeTagList = strings.Split(*includeTags, ",")
	}
	excludeTagList := []string{}
	if *excludeTags != "" {
		excludeTagList = strings.Split(*excludeTags, ",")
	}

	// 加载环境配置
	envConfig, err := internal.LoadEnvConfig(filepath.Join(*configDir, "env_config.yaml"))
	if err != nil {
		log.Fatalf("Failed to load environment config: %v", err)
	}

	// 加载测试用例
	tests, err := internal.LoadTestCases(*testsDir, envConfig)
	if err != nil {
		log.Fatalf("Failed to load test cases: %v", err)
	}

	// 创建测试运行器
	runner, err := runner.NewRunner(*reportsDir, *format)
	if err != nil {
		log.Fatalf("Failed to create test runner: %v", err)
	}

	// 运行测试
	err = runner.Run(tests, includeTagList, excludeTagList)
	if err != nil {
		log.Fatalf("Test run failed: %v", err)
	}
}

// 分割项目字符串
func splitProjects(projectStr string) []string {
	return strings.Split(projectStr, ",")
}
