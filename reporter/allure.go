package reporter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// GenerateAllureReport 生成Allure报告
func GenerateAllureReport(reportDir string) error {
	// 检查Allure命令是否存在
	allureCmd := "allure"
	if _, err := exec.LookPath(allureCmd); err != nil {
		return fmt.Errorf("Allure command not found. Please install Allure CLI first")
	}

	// 生成Allure报告
	cmd := exec.Command(allureCmd, "generate", reportDir, "--output", filepath.Join(reportDir, "html"), "--clean")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate Allure report: %w", err)
	}

	// 打印报告位置
	reportPath, _ := filepath.Abs(filepath.Join(reportDir, "html", "index.html"))
	fmt.Printf("Allure report generated at: file://%s\n", reportPath)

	return nil
}
