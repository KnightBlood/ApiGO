# ApiGO - 接口自动化测试框架

基于Go语言的多功能接口自动化测试框架，支持YAML/JSON/TOML配置，提供Allure报告和原生报告两种输出方式。

## 📌 功能特性
- ✅ 支持GET/POST等HTTP方法
- ✅ 参数化测试（Data-Driven Testing）
- ✅ 前置请求依赖管理（支持循环依赖检测）
- ✅ 多文件测试用例管理
- ✅ 多项目支持（可指定单个或多个被测项目）
- ✅ Allure报告和原生报告双输出
- ✅ 环境配置管理

## 🛠️ 安装指南
```bash
# 克隆项目
$ git clone https://github.com/yourname/apigo.git

# 进入项目目录
$ cd apigo

# 安装依赖
$ go mod download
```

## 📂 目录结构
```bash
├── config/                # 配置文件目录
│   ├── env_config.yaml    # 环境配置文件
│   └── test_config.yaml   # 测试配置文件
├── testcases/             # 测试用例目录（支持多文件）
├── internal/              # 内部逻辑实现
├── reports/               # 测试报告输出目录
│   └── allure/            # Allure原始数据
├── main_test.go           # 测试入口文件
├── go.mod                 # Go模块配置
└── README.md              # 项目说明文件
```

## ⚙️ 环境配置
`config/env_config.yaml` - 环境配置示例：
```yaml
login_endpoint:
  url: "/api/v1/login"
  method: "POST"
  token_field: "Authorization"

projects:
  - name: "projectA"
    base_url: "http://projectA/api/v1"
    description: "用户管理系统"
  - name: "projectB"
    base_url: "http://projectB/api/v2"
    description: "订单管理系统"

active_projects:
  - "projectA"
default_project: "projectA"
global_base_url: "http://default/api"
```

## 📋 测试用例格式
支持以下三种测试用例格式：
1. YAML (.yaml, .yml) - 推荐用于复杂配置
2. JSON (.json) - 推荐用于程序生成的配置
3. TOML (.toml) - 推荐用于简单配置

所有格式的测试用例可以共存于testcases目录，框架会自动识别并加载。

## 💾 构建项目
```bash
# 安装依赖
$ go mod download

# 构建二进制文件（自动生成dist目录包含所有必要文件）
$ build.bat
```

## ▶️ 运行测试
支持所有测试过滤选项：
```bash
# 运行所有测试
$ dist\apigo.exe ^
    --config config ^
    --tests testcases ^
    --reports reports ^
    --format allure

# 运行指定项目的测试
$ dist\apigo.exe ^
    --config config ^
    --tests testcases ^
    --reports reports ^
    --project projectA,projectB ^
    --format allure

# 按标签运行测试
$ dist\apigo.exe ^
    --config config ^
    --tests testcases ^
    --reports reports ^
    --include-tags smoke ^
    --exclude-tags wip

# 按优先级运行测试
$ dist\apigo.exe ^
    --config config ^
    --tests testcases ^
    --reports reports ^
    --min-priority 1
```

## 📦 分发包结构
构建完成后，dist/ 目录将包含：
```bash
├── apigo.exe                # 可执行文件
├── config/                  # 配置文件目录
├── testcases/               # 测试用例目录
├── reports/                 # 测试报告输出目录
├── run.bat                  # Windows运行脚本
└── README.txt               # 运行说明文档
```

## 📁 配置文件
支持以下三种环境配置格式（以YAML为主）：
1. YAML (.yaml, .yml) - 推荐格式
2. JSON (.json)
3. TOML (.toml)

## ▶️ 运行测试
支持所有格式的配置文件：
```bash
# 使用YAML配置运行
$ dist\apigo.exe --config config\env_config.yaml

# 使用JSON配置运行
$ dist\apigo.exe --config config\env_config.json

# 使用TOML配置运行
$ dist\apigo.exe --config config\env_config.toml
```

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

## 📊 测试报告
- **原生报告**：控制台实时输出
- **Allure报告**：
  ```bash
  # 安装Allure命令行工具
  $ npm install -g allure-commandline

  # 查看报告
  $ allure open reports/allure/
  ```

## 📌 注意事项
1. YAML格式必须严格符合规范
2. 前置请求配置需注意循环依赖
3. 参数提取使用JSON路径语法（如：$.token）
4. 多文件测试用例共享相同的环境配置
5. Allure报告需要Node.js环境支持

## 📚 示例项目
项目示例：[https://github.com/yourname/apigo](https://github.com/yourname/apigo)