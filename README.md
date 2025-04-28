# Graph-Med-MCP 医疗知识图谱服务


Graph-Med-MCP 是一个基于 Go 语言开发的MCP服务，提供向[基于知识图谱的智能问答系统](https://github.com/poisunk/graph-med)，使 AI 模型能够查询和获取医疗知识图谱数据。


### 配置文件

```yaml
port: 9001

neo4j:
  addr: bolt://localhost:7687
  username: neo4j
  password: your_password
```

## 运行服务

默认情况下，服务将在 9001 端口启动

```bash
./graph-med-mcp -f path/to/config.yaml
```

## API 使用说明

### 获取疾病知识图谱

该服务提供了 `get_disease_subgraph` 工具，用于查询指定疾病的相关知识图谱。

请求参数：
- `disease_name`: 疾病名称（必填）

响应格式：
- 返回 JSON 格式的知识图谱数据，包含节点（nodes）和边（edges）信息
