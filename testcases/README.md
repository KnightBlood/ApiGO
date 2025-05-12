# 测试用例目录

本目录包含所有测试用例的YAML配置文件，按项目组织为子目录。

## 目录结构
```
testcases/
├── projectA/          # 项目A测试用例
│   ├── user_tests.yaml  # 用户管理测试
│   └── order_tests.yaml # 订单管理测试
├── projectB/          # 项目B测试用例
│   ├── product_tests.yaml # 产品管理测试
│   └── admin_tests.yaml   # 管理功能测试
└── README.md          # 本说明文件
```

## 测试用例编写指南

### 测试用例结构说明
```yaml
tests:
  - name: "测试用例名称"
    project: "所属项目"
    url: "接口路径"
    method: "HTTP方法"
    body_template:        # 请求体模板（可选）
      field1: "value1"
      field2: "value2"
    params_template:      # 查询参数模板（可选）
      param1: "value1"
    parameters:           # 参数组合及预期结果
      - values:
          param1: "value1"
          param2: "value2"
        expected_status: 200
        expected_body: '{"key":"value"}'
    expected_status: 200
    expected_body: '{"key":"value"}'
    pre_request:          # 前置请求（可选）
      name: "前置用例名称"
      extract:            # 从前置用例提取参数
        token: "$.token" # JSON路径提取
        id: "$.id"
```

### 关键字段说明
- `project`: 必须与env_config.yaml中定义的项目名称匹配
- `pre_request`: 指定前置请求，支持参数传递和循环依赖检测
- `parameters`: 支持参数化测试，每个参数组合会生成独立测试
- `extract`: 支持JSON路径参数提取（如`$.token`）
- `body_template`: 请求体模板，支持动态参数替换
- `params_template`: 查询参数模板，支持动态参数替换

### 编写建议
1. 每个功能模块使用单独的YAML文件
2. 保持测试用例文件大小适中（建议不超过50个用例/文件）
3. 使用清晰的测试用例名称，推荐格式：`操作_场景_预期结果`
4. 为每个测试用例编写明确的预期结果
5. 合理使用前置请求传递参数
6. 参数化测试应覆盖正向和负向场景
7. 避免循环依赖（如A依赖B，B又依赖A）
8. 保持YAML格式正确，注意缩进和特殊字符

### 测试用例示例
```yaml
- name: "Login Success"
  project: "projectA"
  url: "/login"
  method: "POST"
  body_template:
    username: "admin"
    password: "123456"
  parameters:
    - values:
        username: "admin"
        password: "wrongpass"
      expected_status: 401
      expected_body: '{"error":"Invalid credentials"}'
  expected_status: 200
  expected_body: '{"token":"abc123xyz"}'
```

### 运行测试
1. 修改env_config.yaml配置基础URL和项目信息
2. 在项目根目录运行测试：`go test -v`
3. 构建为二进制后运行：`dist/apigo.exe --config config --tests testcases --project projectA,projectB --format allure`
```

## 测试用例目录结构
```
testcases/
├── projectA/          # 项目A测试用例
│   ├── user_tests.yaml  # 用户管理测试（YAML格式）
│   ├── order_tests.yaml # 订单管理测试（YAML格式）
│   └── api_docs.toml    # API文档测试（TOML格式）
├── projectB/          # 项目B测试用例
│   ├── product_tests.yaml # 产品管理测试（YAML格式）
│   ├── admin_tests.yaml   # 管理功能测试（YAML格式）
│   └── inventory_tests.json # 库存管理测试（JSON格式）
└── README.md          # 本说明文件
```

## 测试用例编写指南

### 支持的配置格式
支持以下三种配置格式，可以混合使用：
1. YAML (.yaml, .yml)
2. JSON (.json)
3. TOML (.toml)

### 测试用例结构说明
不同格式的测试用例结构保持一致，以下是各种格式的等效示例：

#### YAML格式
```yaml
tests:
  - name: "Login Success"
    project: "projectA"
    url: "/login"
    method: "POST"
    body_template:
      username: "admin"
      password: "123456"
    parameters:
      - values:
          username: "admin"
          password: "wrongpass"
        expected_status: 401
        expected_body: '{"error":"Invalid credentials"}'
  expected_status: 200
  expected_body: '{"token":"abc123xyz"}'
```

#### JSON格式
```json
{
  "tests": [
    {
      "name": "Login Success",
      "project": "projectA",
      "url": "/login",
      "method": "POST",
      "body_template": {
        "username": "admin",
        "password": "123456"
      },
      "parameters": [
        {
          "values": {
            "username": "admin",
            "password": "wrongpass"
          },
          "expected_status": 401,
          "expected_body": '{"error":"Invalid credentials"}'
        }
      ],
      "expected_status": 200,
      "expected_body": '{"token":"abc123xyz"}'
    }
  ]
}
```

#### TOML格式
```toml
[[tests]]
name = "Login Success"
project = "projectA"
url = "/login"
method = "POST"
body_template = { username = "admin", password = "123456" }

[[tests.parameters]]
values = { username = "admin", password = "wrongpass" }
expected_status = 401
expected_body = '{"error":"Invalid credentials"}'

expected_status = 200
expected_body = '{"token":"abc123xyz"}'
```

### 编写建议
1. 每个功能模块使用单独的配置文件
2. 保持测试用例文件大小适中（建议不超过50个用例/文件）
3. 支持混合使用YAML/TOML/JSON格式
4. 不同格式的测试用例可以共存于同一目录
5. 根据团队偏好选择合适的配置格式
6. 保持配置文件格式一致性（建议单个项目使用统一格式）
7. 使用清晰的测试用例名称，推荐格式：`操作_场景_预期结果`
8. 为每个测试用例编写明确的预期结果
9. 合理使用前置请求传递参数
10. 参数化测试应覆盖正向和负向场景
11. 避免循环依赖（如A依赖B，B又依赖A）
12. 注意不同格式的特殊语法要求

## 测试用例结构对比
- **YAML格式**：
```yaml
tests:
  - name: "Login Success"
    project: "projectA"
    url: "/login"
    method: "POST"
    tags: ["smoke", "auth"]
    priority: 2
    body_template:
      username: "admin"
      password: "123456"
    parameters:
      - values:
          username: "admin"
          password: "wrongpass"
        expected_status: 401
        expected_body: '{"error":"Invalid credentials"}'
        tags: ["negative", "auth"]
        priority: 1
    expected_status: 200
    expected_body: '{"token":"abc123xyz"}'
```

- **JSON格式**：
```json
{
  "tests": [
    {
      "name": "Login Success",
      "project": "projectA",
      "url": "/login",
      "method": "POST",
      "tags": ["smoke", "auth"],
      "priority": 2,
      "body_template": {
        "username": "admin",
        "password": "123456"
      },
      "parameters": [
        {
          "values": {
            "username": "admin",
            "password": "wrongpass"
          },
          "expected_status": 401,
          "expected_body": '{"error":"Invalid credentials"}',
          "tags": ["negative", "auth"],
          "priority": 1
        }
      ],
      "expected_status": 200,
      "expected_body": '{"token":"abc123xyz"}'
    }
  ]
}
```

- **TOML格式**：
```toml
[[tests]]
name = "Login Success"
project = "projectA"
url = "/login"
method = "POST"
tags = ["smoke", "auth"]
priority = 2

[tests.body_template]
username = "admin"
password = "123456"

[[tests.parameters]]
[tests.parameters.values]
username = "admin"
password = "wrongpass"

expected_status = 401
expected_body = '{"error":"Invalid credentials"}'
tags = ["negative", "auth"]
priority = 1

expected_status = 200
```

### 格式选择建议
1. 推荐使用YAML作为主要配置格式
2. JSON适合程序生成或API交互
3. TOML适合简单配置
4. 所有格式的配置文件可以共存
5. 框架会自动识别配置文件格式
6. 不同格式的测试用例可以混合使用
7. 所有格式支持完全相同的测试特性

## 环境配置示例
```
config/
├── env_config.yaml   # 推荐：YAML格式的环境配置
├── env_config.json   # JSON格式的环境配置
└── env_config.toml   # TOML格式的环境配置
```

## 环境配置格式说明
支持以下三种环境配置格式：
1. YAML (.yaml, .yml) - 推荐格式
2. JSON (.json) - 适合程序生成
3. TOML (.toml) - 适合简单配置

### YAML格式（推荐）
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

### JSON格式
```json
{
  "login_endpoint": {
    "url": "/api/v1/login",
    "method": "POST",
    "token_field": "Authorization"
  },
  "projects": [
    {
      "name": "projectA",
      "base_url": "http://projectA/api/v1",
      "description": "用户管理系统"
    },
    {
      "name": "projectB",
      "base_url": "http://projectB/api/v2",
      "description": "订单管理系统"
    }
  ],
  "active_projects": ["projectA"],
  "default_project": "projectA",
  "global_base_url": "http://default/api"
}
```

### TOML格式
```toml
[login_endpoint]
url = "/api/v1/login"
method = "POST"
token_field = "Authorization"

[[projects]]
name = "projectA"
base_url = "http://projectA/api/v1"
description = "用户管理系统"

[[projects]]
name = "projectB"
base_url = "http://projectB/api/v2"
description = "订单管理系统"

active_projects = ["projectA"]
default_project = "projectA"
global_base_url = "http://default/api"
```

### 格式选择建议
1. 推荐使用YAML作为主要配置格式
2. JSON适合程序生成或API交互
3. TOML适合简单配置
4. 所有格式的配置文件可以共存
5. 框架会自动识别配置文件格式
