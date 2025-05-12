package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"ApiGO/internal"
	"ApiGO/reporter"
)

// Runner 测试运行器
type Runner struct {
	ReportDir  string
	Format     string
	TestCount  int
	FailCount  int
	StartTime  time.Time
	FinishTime time.Time
}

// 创建测试运行器
func NewRunner(reportDir string, format string) (*Runner, error) {
	// 创建报告目录
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create report directory: %w", err)
	}

	return &Runner{
		ReportDir:  reportDir,
		Format:     format,
		TestCount:  0,
		FailCount:  0,
		StartTime:  time.Now(),
	}, nil
}

// Run方法新增参数
func (r *Runner) Run(tests []internal.TestCase, includeTags []string, excludeTags []string) error {
	// 按标签过滤
	sortedTests := sortTestsByPriority(tests)
	sortedTests = filterTestsByTags(sortedTests, includeTags, excludeTags)

	// 按项目分组
	projectGroups := make(map[string][]internal.TestCase)
	for _, test := range sortedTests {
		project := test.Project
		if project == "" {
			project = "default"
		}
		projectGroups[project] = append(projectGroups[project], test)
	}

	// 生成报告
	switch r.Format {
	case "console", "":
		return r.runConsoleReport(projectGroups)
	case "allure":
		return r.runAllureReport(projectGroups)
	default:
		return fmt.Errorf("unsupported report format: %s", r.Format)
	}
}

// 控制台报告
func (r *Runner) runConsoleReport(projectGroups map[string][]internal.TestCase) error {
	// 实现控制台报告逻辑
	for projectName, projectTests := range projectGroups {
		err := runProjectTests(projectName, projectTests, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// Allure报告
func (r *Runner) runAllureReport(projectGroups map[string][]internal.TestCase) error {
	// 创建Allure结果目录
	resultsDir := filepath.Join(r.ReportDir, "allure", "results")
	if err := os.RemoveAll(resultsDir); err != nil {
		return fmt.Errorf("failed to clean allure results directory: %w", err)
	}
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create allure results directory: %w", err)
	}

	// 执行测试并生成Allure结果文件
	for projectName, projectTests := range projectGroups {
		err := runProjectTests(projectName, projectTests, r, resultsDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// 执行项目测试用例
func runProjectTests(projectName string, tests []internal.TestCase, r *Runner, resultsDir string) error {
	// 获取项目配置
	projectConfig := getProjectConfig(r, projectName)

	// 创建项目专用的测试上下文
	projectCtx := &internal.TestContext{
		EnvConfig: &internal.EnvConfig{
			LoginEndpoint:  r.EnvConfig.LoginEndpoint,
			Projects:       r.EnvConfig.Projects,
			DefaultProject: r.EnvConfig.DefaultProject,
			GlobalBaseUrl:  projectConfig.BaseUrl,
		},
		AuthHeader:    r.AuthHeader,
		LoginEndpoint: r.LoginEndpoint,
	}

	// 执行项目测试用例
	for _, test := range tests {
		// 执行测试用例
		err := executeTest(projectCtx, test, resultsDir)
		if err != nil {
			return fmt.Errorf("test failed: %w", err)
		}
	}

	return nil
}

// 设置测试上下文
func setupTestContext(envConfig *internal.EnvConfig) (*internal.TestContext, error) {
	// 初始化登录
	loginResult, err := internal.Login(envConfig)
	if err != nil {
		return nil, err
	}

	return &internal.TestContext{
		EnvConfig:     envConfig,
		AuthHeader:    loginResult.Token,
		LoginEndpoint: loginResult.LoginEndpoint,
	}, nil
}

// 执行单个测试用例
func executeTest(ctx *internal.TestContext, test internal.TestCase, resultsDir string) error {
	// 构造请求
	req, err := internal.BuildRequest(ctx, test)
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	// 发送请求
	resp, err := internal.SendRequest(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	// 验证响应
	err = internal.ValidateResponse(resp, test)
	if err != nil {
		return fmt.Errorf("response validation failed: %w", err)
	}

	// 生成Allure结果文件
	if resultsDir != "" {
		err = reporter.GenerateAllureResultFile(resultsDir, test, resp)
		if err != nil {
			return fmt.Errorf("failed to generate allure result file: %w", err)
		}
	}

	return nil
}

// 获取项目配置
func getProjectConfig(r *Runner, projectName string) *internal.ProjectConfig {
	for _, project := range r.EnvConfig.Projects {
		if project.Name == projectName {
			return &project
		}
	}
	// 如果未找到项目配置，返回默认配置
	return &internal.ProjectConfig{
		Name:    projectName,
		BaseUrl: r.EnvConfig.GlobalBaseUrl,
	}
}

// 按优先级排序测试用例
func sortTestsByPriority(tests []internal.TestCase) []internal.TestCase {
	sorted := make([]internal.TestCase, len(tests))
	copy(sorted, tests)

	sort.Slice(sorted, func(i, j int) bool {
		// 高优先级先执行，同优先级按名称排序
		if sorted[i].Priority != sorted[j].Priority {
			return sorted[i].Priority > sorted[j].Priority
		}
		return sorted[i].Name < sorted[j].Name
	})

	return sorted
}

// 按标签过滤测试用例
func filterTestsByTags(tests []internal.TestCase, includeTags []string, excludeTags []string) []internal.TestCase {
	if len(includeTags) == 0 && len(excludeTags) == 0 {
		return tests
	}

	filtered := make([]internal.TestCase, 0)
	for _, test := range tests {
		// 检查是否包含排除标签
		if hasAnyTag(test.Tags, excludeTags) {
			continue
		}

		// 如果包含包含标签，或者没有指定包含标签
		if len(includeTags) == 0 || hasAnyTag(test.Tags, includeTags) {
			filtered = append(filtered, test)
		}
	}
	return filtered
}

// 检查标签是否存在交集
func hasAnyTag(testTags []string, filterTags []string) bool {
	if len(testTags) == 0 || len(filterTags) == 0 {
		return false
	}

	for _, testTag := range testTags {
		for _, filterTag := range filterTags {
			if testTag == filterTag {
				return true
			}
		}
	}
	return false
}
