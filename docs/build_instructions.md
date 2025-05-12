# 构建指南

本指南描述如何将ApiGO项目构建为可在Windows系统运行的独立二进制程序。

## 📋 准备工作

1. 确保已安装以下工具：
   - Go 1.20+（已添加到系统PATH）
   - Java 8+（用于运行Allure）
   - Allure CLI for Java（通过Java方式安装）

2. 确认项目结构完整：
   ```bash
   ├── config/                # 配置文件目录
   ├── testcases/             # 测试用例目录（支持多文件）
   ├── cmd/apigo/main.go      # 主程序入口
   ├── internal/              # 内部逻辑实现
   ├── runner/                # 测试运行器
   ├── build.bat              # Windows构建脚本
   └── README.md              # 项目说明文件
   ```

## 📁 环境配置格式支持
支持以下三种环境配置格式，可以混合使用：
1. YAML (.yaml, .yml) - 推荐格式
2. JSON (.json) - 适合程序生成
3. TOML (.toml) - 适合简单配置

## 🧱 构建注意事项
1. 所有格式的环境配置和测试用例都会被复制到构建目录
2. 修改任何格式的配置文件后直接运行run.bat即可
3. 不同格式的配置文件可以共存于同一目录
4. 构建脚本会自动处理所有支持的配置格式
5. 推荐使用YAML作为主要配置格式
6. 注意不同格式的特殊语法要求

## 🧱 构建步骤（Windows）

1. 打开命令行工具（CMD或PowerShell）
2. 确保在项目根目录下
3. 运行构建脚本：
   ```bash
   $ build.bat
   ```

4. 构建完成后，dist/目录结构：
   ```bash
   ├── apigo.exe                # 可执行文件
   ├── config/                  # 配置文件目录
   ├── testcases/               # 测试用例目录
   ├── reports/                 # 测试报告输出目录
   │   └── allure/              # Allure结果和报告
   │       └── results/         # Allure结果文件
   ├── run.bat                  # Windows运行脚本
   └── README.txt               # 运行说明文档
   ```

## 🚀 运行指南

### 使用构建后的程序
1. 运行所有测试
```bash
$ apigo.exe --config config --tests testcases --reports reports --format allure
```

2. 运行指定项目
```bash
$ apigo.exe --config config --tests testcases --reports reports --project projectA,projectB
```

3. 按标签运行测试
```bash
# 运行包含'smoke'标签的测试
$ apigo.exe --config config --tests testcases --reports reports --include-tags smoke

# 排除'wip'标签的测试
$ apigo.exe --config config --tests testcases --reports reports --exclude-tags wip

# 同时使用包含和排除标签
$ apigo.exe --config config --tests testcases --reports reports --include-tags smoke --exclude-tags wip
```

4. 按优先级运行测试
```bash
# 运行优先级大于等于1的测试
$ apigo.exe --config config --tests testcases --reports reports --min-priority 1
```

## 📦 分发包结构

构建完成后，dist/目录将包含：
```bash
├── apigo.exe                # 可执行文件
├── config/                  # 配置文件目录
├── testcases/               # 测试用例目录
├── reports/                 # 测试报告输出目录
├── run.bat                  # Windows运行脚本
└── README.txt               # 运行说明文档
```

## 📝 运行注意事项
1. 测试框架自动识别配置文件格式
2. 不需要修改代码即可使用不同格式的测试用例
3. 支持在同一项目中混合使用不同格式的测试用例
4. 参数提取和前置请求功能在所有格式中保持一致

## 📝 注意事项
1. 需要预安装Java 8+环境
2. Allure CLI需要Java运行时（下载Java版本的Allure CLI）
3. 所有测试用例和配置文件都在dist/目录内运行，无需修改代码
4. 修改测试用例后直接运行run.bat即可
5. 生成Allure HTML报告需要手动运行Java命令：
```bash
$ java -jar allure-commandline/bin/allure.jar generate reports/allure/results -o reports/allure/html --clean
```
6. 支持以下测试用例特性：
   - 参数化测试
   - 前置请求
   - 循环依赖检测
   - 多文件用例管理
   - 多项目支持
   - 测试用例标签
   - 测试用例优先级
   - 多格式支持（YAML/JSON/TOML）
