package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"

	allure "github.com/ozontech/allure-go/pkg/framework"
)

// 测试上下文
type TestContext struct {
	EnvConfig     *EnvConfig
	AuthHeader    string
	LoginEndpoint *LoginEndpointConfig
}

// 构造请求
func BuildRequest(ctx *TestContext, test TestCase) (*http.Request, error) {
	// 构造完整URL
	baseUrl := ctx.EnvConfig.GlobalBaseUrl
	if projectConfig := getProjectConfig(ctx.EnvConfig, test.Project); projectConfig != nil {
		baseUrl = projectConfig.BaseUrl
	}

	fullUrl := baseUrl + test.Url

	// 替换URL中的模板参数
	for k, v := range mergeMaps(test.BodyTemplate, test.ParamsTemplate) {
		fullUrl = strings.ReplaceAll(fullUrl, fmt.Sprintf("{{%s}}", k), v)
	}

	// 构造请求体
	body := make(map[string]interface{})
	// 添加模板参数
	for k, v := range test.BodyTemplate {
		body[k] = v
	}

	// 创建请求
	req, err := http.NewRequest(test.Method, fullUrl, nil)
	if err != nil {
		return nil, err
	}

	// 设置内容类型头
	if test.Method == "POST" || test.Method == "post" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 添加认证头
	if ctx.AuthHeader != "" && test.Name != "login" {
		req.Header.Set(ctx.LoginEndpoint.TokenField, ctx.AuthHeader)
	}

	// 返回请求
	return req, nil
}

// 发送请求
func SendRequest(req *http.Request) (*http.Response, error) {
	// 创建客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 验证响应并生成Allure结果
func ValidateResponse(resp *http.Response, test TestCase, resultsDir string) error {
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 创建测试结果文件
	result := allure.Result{
		Uuid:        generateUUID(),
		Name:        test.Name,
		Status:      "passed",
		Stage:       "finished",
		Description: fmt.Sprintf("Test case: %s", test.Name),
		Result: allure.ResultData{
			Start:    time.Now().UnixMilli(),
			Stop:     time.Now().UnixMilli() + 100,
			StatusMsg: "Test execution completed",
		},
	}

	// 验证状态码
	if resp.StatusCode != test.DefaultStatus {
		result.Status = "failed"
		result.Result.StatusMsg = fmt.Sprintf("Expected status %d but got %d", test.DefaultStatus, resp.StatusCode)
	}

	// 验证响应体
	if string(body) != test.DefaultBody {
		result.Status = "failed"
		result.Result.StatusMsg = fmt.Sprintf("Expected body '%s' but got '%s'", test.DefaultBody, string(body))
	}

	// 保存测试结果
	err = saveResultToFile(result, resultsDir)
	if err != nil {
		return fmt.Errorf("failed to save test result: %w", err)
	}

	return nil
}

// 生成唯一ID
func generateUUID() string {
	// 实际实现应使用更可靠的UUID生成方式
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// 保存测试结果到文件
func saveResultToFile(result allure.Result, resultsDir string) error {
	// 创建结果文件
	ofilePath := filepath.Join(resultsDir, fmt.Sprintf("%s-result.json", result.Uuid))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 序列化结果
	encoder := json.NewEncoder(file)
	err = encoder.Encode(result)
	if err != nil {
		return err
	}

	return nil
}

// Allure结果结构
type Result struct {
	Uuid  string   `json:"uuid"`
	Name  string   `json:"name"`
	Status string `json:"status"`
	Stage string  `json:"stage"`
	Description string `json:"description"`
	Result ResultData `json:"result"`
}

// ResultData 结构
type ResultData struct {
	Start   int64  `json:"start"`
	Stop    int64  `json:"stop"`
	StatusMsg string `json:"statusMessage"`
}

// 获取项目配置
func getProjectConfig(envConfig *EnvConfig, projectName string) *ProjectConfig {
	for _, project := range envConfig.Projects {
		if project.Name == projectName {
			return &project
	}
	// 如果未找到项目配置，返回默认配置
	return &ProjectConfig{
		Name:    projectName,
		BaseUrl: envConfig.GlobalBaseUrl,
	}
}