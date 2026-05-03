# 配置加载系统修复说明

## 问题
之前系统存在的问题：
- 配置文件格式错误时，服务静默启动失败
- 无清晰的错误提示（缺少文件路径、行号、错误类型）
- 各服务硬编码配置值
- 不支持 YAML 严格模式检查（检测多余字段）

## 修复内容

### 1. Go 服务 (go/config.go)
- 使用 `gopkg.in/yaml.v3` 进行配置解析
- `decoder.KnownFields(true)` 启用严格模式，检测未知字段
- 详细的错误类型：
  - YAML syntax error
  - unknown configuration field
  - type mismatch error
  - validation error
- 错误包含文件路径和行号
- 使用 `os.Exit(1)` 退出程序（非静默失败）
- 支持 `CONFIG_PATH` 环境变量指定配置文件路径
- 完整的配置验证逻辑

### 2. Node.js 服务 (nodejs/config.js)
- 使用 `js-yaml` 库进行配置解析
- `try-catch` 捕获所有解析异常
- 从 YAML 错误中提取行号和错误原因
- 详细的错误类型检测和格式化
- 验证所有必需字段和值范围
- 失败时输出所有验证错误列表后退出

### 3. C++ 服务 (cpp/include/config.h, cpp/src/config.cpp)
- 自建 YAML 解析器（避免外部依赖）
- 严格的缩进检查和格式验证
- 类型转换异常捕获
- 未知字段检测
- 文件存在性检查
- 完整的配置验证
- 错误输出包含所有详细信息

## 新增文件
```
config/config.yaml              # 新的 YAML 配置文件
go/config.go                    # Go 配置加载
nodejs/config.js                # Node.js 配置加载
cpp/include/config.h            # C++ 配置头文件
cpp/src/config.cpp              # C++ 配置实现
config/config_bad_indent.yaml   # 测试：缩进错误
config/config_bad_type.yaml     # 测试：类型错误
config/config_extra_field.yaml  # 测试：多余字段
CONFIG_FIX_README.md            # 本文档
```

## 使用方式

### 设置配置文件路径
所有服务均支持 `CONFIG_PATH` 环境变量：

```bash
# Linux/Mac
export CONFIG_PATH=/path/to/config.yaml
./log_receiver

# Windows PowerShell
$env:CONFIG_PATH = "C:\path\to\config.yaml"
.\log_receiver.exe
```

### 错误输出示例

#### 1. 配置文件不存在
```
ERROR: configuration file not found in /path/to/missing.yaml: the file does not exist or is not a regular file
```

#### 2. YAML 语法错误
```
ERROR: YAML syntax error in /path/to/config.yaml (line 5): expected a key-value pair
```

#### 3. 类型不匹配
```
ERROR: type mismatch error in /path/to/config.yaml (line 3): invalid numeric value for key 'udp_port'
```

#### 4. 未知字段（Go/Node.js）
```
ERROR: unknown configuration field in /path/to/config.yaml (line 9): unexpected key: this_field_does_not_exist
```

#### 5. 验证错误（多个）
```
ERROR: configuration validation failed:
  - cpp_receiver.udp_port must be between 1 and 65535, got 0
  - cpp_receiver.redis_host is required
```

## 配置文件格式
参考 `config/config.yaml`：
```yaml
cpp_receiver:
  udp_port: 514
  tcp_port: 515
  queue_size: 100000
  redis_host: "localhost"
  redis_port: 6379
  redis_stream: "log_stream"

go_consumer:
  redis_host: "localhost"
  redis_port: 6379
  redis_stream: "log_stream"
  consumer_group: "log_consumers"
  consumer_name: "consumer_1"
  metrics_interval: 60
  grpc_port: 50051

node_api:
  http_port: 3000
  grpc_host: "localhost"
  grpc_port: 50051
  redis_host: "localhost"
  redis_port: 6379
```

## 测试方法
使用提供的测试配置文件验证错误检测：

```bash
# 测试缩进错误
CONFIG_PATH=config/config_bad_indent.yaml go run go/*.go

# 测试类型错误
CONFIG_PATH=config/config_bad_type.yaml node nodejs/main.js
```

## 退出码
所有服务在配置加载失败时均返回非零退出码（1），便于脚本检测失败。
