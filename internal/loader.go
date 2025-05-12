package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// 测试用例结构体
type TestCase struct {
	Name           string                 `mapstructure:"name"`
	Url            string                 `mapstructure:"url"`
	Method         string                 `mapstructure:"method"`
	BodyTemplate   map[string]interface{} `mapstructure:"body_template"`  // 请求体模板
	ParamsTemplate map[string]string      `mapstructure:"params_template"` // 查询参数模板
	Parameters     []TestParameter        `mapstructure:"parameters"`      // 参数组合及预期结果
	DefaultStatus  int                    `mapstructure:"expected_status"`// 默认状态码
	DefaultBody    string                 `mapstructure:"expected_body"`  // 默认响应体
	PreRequest     *PreRequestConfig      `mapstructure:"pre_request"`    // 前置请求配置
	SourceFile     string                 `mapstructure:"-"`              // 源文件路径（不从配置文件解析）
	Project        string                 `mapstructure:"project"`        // 所属项目
}

// 测试参数定义
type TestParameter struct {
	Values         map[string]string `mapstructure:"values"`              // 请求参数值
	ExpectedStatus int               `mapstructure:"expected_status"`     // 关联状态码
	ExpectedBody   string            `mapstructure:"expected_body"`       // 关联响应体
}

// 前置请求配置
type PreRequestConfig struct {
	Name       string            `mapstructure:"name"`          // 前置请求名称
	PathParam  map[string]string `mapstructure:"path_param"`    // 路径参数
	QueryParam map[string]string `mapstructure:"query_param"`   // 查询参数
	BodyParam  map[string]string `mapstructure:"body_param"`    // 请求体参数
	Extract    map[string]string `mapstructure:"extract"`       // 响应参数提取规则
	LoopCount  int               `mapstructure:"loop_count"`    // 循环次数
}

// 环境配置中的接口定义
type EndpointConfig struct {
	Name       string            `mapstructure:"name"`
	Url        string            `mapstructure:"url"`
	Method     string            `mapstructure:"method"`
	Body       map[string]string `mapstructure:"body"` // 支持多参数body
	TokenField string            `mapstructure:"token_field"`
}

// 加载测试用例（支持多文件和项目过滤）
func LoadTestCases(configDir string, envConfig *EnvConfig) ([]TestCase, error) {
	var allTests []TestCase

	// 获取要运行的项目列表
	activeProjects := getActiveProjects(envConfig)

	// 遍历配置目录下的所有YAML文件
	err := filepath.WalkDir(configDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 仅处理YAML文件
		if !d.IsDir() && (filepath.Ext(d.Name()) == ".yaml" || filepath.Ext(d.Name()) == ".yml") {
			tests, err := loadTestCasesFromFile(path, activeProjects)
			if err != nil {
				return fmt.Errorf("failed to load test cases from %s: %w", path, err)
			}
			allTests = append(allTests, tests...)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return allTests, nil
}

// 获取要运行的项目列表
func getActiveProjects(envConfig *EnvConfig) []string {
	if len(envConfig.ActiveProjects) > 0 {
		return envConfig.ActiveProjects
	}
	if envConfig.DefaultProject != "" {
		return []string{envConfig.DefaultProject}
	}
	// 如果未指定活动项目和默认项目，返回空列表表示运行所有项目
	return []string{}
}

// 从单个文件加载测试用例（支持项目过滤）
func loadTestCasesFromFile(filePath string, activeProjects []string) ([]TestCase, error) {
	// 设置配置文件和类型
	viper.SetConfigFile(filePath)
	
	// 根据文件扩展名设置配置类型
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".yaml", ".yml":
		viper.SetConfigType("yaml")
	case ".json":
		viper.SetConfigType("json")
	case ".toml":
		viper.SetConfigType("toml")
	default:
		return nil, fmt.Errorf("unsupported config file type: %s", filepath.Ext(filePath))
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var tests []TestCase
	// 根据文件类型解码
	switch viper.GetConfigType() {
	case "json":
		if err := viper.UnmarshalKey("tests", &tests); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON tests: %w", err)
		}
	case "toml":
		// TOML需要使用Unmarshal而是UnmarshalKey
		if err := viper.Unmarshal(&tests); err != nil {
			return nil, fmt.Errorf("failed to unmarshal TOML tests: %w", err)
		}
	default:
		// 默认使用YAML格式的UnmarshalKey方式
		if err := viper.UnmarshalKey("tests", &tests); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tests: %w", err)
		}
	}

	// 过滤测试用例
	filteredTests := make([]TestCase, 0, len(tests))
	for i := range tests {
		test := tests[i]
		test.SourceFile = filePath

		// 如果未指定活动项目，或者测试用例的项目在活动项目列表中
		if len(activeProjects) == 0 || contains(activeProjects, test.Project) {
			filteredTests = append(filteredTests, test)
		}
	}

	return filteredTests, nil
}

// 检查字符串是否存在于切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 加载环境配置
func LoadEnvConfig(configFile string) (*EnvConfig, error) {
	viper.SetConfigFile(configFile)

	// 根据文件扩展名设置配置类型
	switch strings.ToLower(filepath.Ext(configFile)) {
	case ".yaml", ".yml":
		viper.SetConfigType("yaml")
	case ".json":
		viper.SetConfigType("json")
	case ".toml":
		viper.SetConfigType("toml")
	default:
		return nil, fmt.Errorf("unsupported config file type: %s", filepath.Ext(configFile))
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	var envConfig EnvConfig

	// 根据配置类型使用不同的解码方式
	switch viper.GetConfigType() {
	case "json":
		if err := viper.Unmarshal(&envConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON config: %w", err)
		}
	case "toml":
		if err := viper.Unmarshal(&envConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal TOML config: %w", err)
		}
	default:
		// 默认使用YAML格式的解码方式
		if err := viper.Unmarshal(&envConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return &envConfig, nil
}