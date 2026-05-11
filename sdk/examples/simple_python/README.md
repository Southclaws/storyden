# Simple Python Plugin Example

# 简单 Python 插件示例

## Run Locally (External Mode)

## 本地运行（外部模式）

```bash
uv run python main.py
```

在环境变量中设置 `STORYDEN_RPC_URL`。

Set `STORYDEN_RPC_URL` in your environment.

## Manifest（清单）

## Manifest

`manifest.yaml` 可用于：

Use `manifest.yaml` for:
- 外部插件设置（将 JSON/YAML 字段粘贴到 manifest input 中）
- 托管插件打包（在插件归档中包含 `manifest.yaml`）
- External plugin setup (paste JSON/YAML fields into manifest input)
- Supervised plugin packaging (include `manifest.yaml` in your plugin archive)
