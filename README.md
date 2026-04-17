# gocommit

AI 驱动的中文 Git 提交信息生成工具。

分析你的代码变更，自动生成符合 Conventional Commits 规范的中文提交信息。

## 特性

- **中文优先** — 提交描述使用中文，符合中文开发者习惯
- **Conventional Commits** — 自动识别变更类型，生成规范格式
- **OpenAI 兼容** — 支持任何 OpenAI 兼容 API（OpenAI、DeepSeek、Ollama 等）
- **交互确认** — 生成后可确认、编辑、重试或取消
- **Agent Skill** — 可安装为 OpenCode/Claude Code 的 Skill，在 AI 对话中直接使用
- **Go 单二进制** — 无运行时依赖，一个文件搞定

## 安装

```bash
go install github.com/xinranzhao/gocommit/cmd/gocommit@latest
```

## 快速开始

```bash
# 1. 配置 API Key
gocommit config api_key=sk-xxx

# 或使用环境变量
export OPENAI_API_KEY=sk-xxx

# 2. 暂存你的变更
git add .

# 3. 生成提交信息
gocommit
```

## 使用

### 基本用法

```bash
git add <files...>
gocommit
```

会显示变更概览，生成中文提交信息，等待你确认：

```
 src/auth.go | 15 +++++++
 1 file changed, 15 insertions(+)

Generating commit message... done!

  feat: 添加JWT认证中间件

Commit with this message? [y/e/r/n] (y=yes, e=edit, r=retry, n=abort):
```

### 命令行参数

| 参数 | 短 | 说明 |
|------|----|------|
| `--all` | `-a` | 自动暂存所有变更 |
| `--amend` | | 修改上一次提交 |
| `--yes` | `-y` | 跳过确认，自动提交 |
| `--dry-run` | `-d` | 只显示消息，不提交 |

### 配置

```bash
# 查看当前配置
gocommit config

# 设置配置项
gocommit config api_key=sk-xxx
gocommit config base_url=https://api.deepseek.com/v1
gocommit config model=deepseek-chat
```

也支持环境变量（优先级高于配置文件）：

```bash
export OPENAI_API_KEY=sk-xxx
export OPENAI_BASE_URL=https://api.deepseek.com/v1
export OPENAI_MODEL=deepseek-chat
```

### 使用 Ollama 本地模型

```bash
gocommit config base_url=http://localhost:11434/v1
gocommit config model=qwen2.5:7b
gocommit config api_key=ollama
```

## 安装为 Agent Skill

### OpenCode

```bash
# 克隆仓库
git clone https://github.com/xinranzhao/gocommit.git /tmp/gocommit

# 安装 Skill
mkdir -p ~/.config/opencode/skills/gocommit-commit
cp /tmp/gocommit/skill/gocommit-commit/SKILL.md ~/.config/opencode/skills/gocommit-commit/
rm -rf /tmp/gocommit
```

安装后，在 OpenCode 对话中直接说 "帮我提交" 即可触发。

### Claude Code

```bash
git clone https://github.com/xinranzhao/gocommit.git /tmp/gocommit
mkdir -p ~/.claude/skills/gocommit-commit
cp /tmp/gocommit/skill/gocommit-commit/SKILL.md ~/.claude/skills/gocommit-commit/
rm -rf /tmp/gocommit
```

## 提交信息格式

生成的提交信息遵循 Conventional Commits 规范，描述使用中文：

```
类型(范围): 中文描述
```

**类型**：`feat` | `fix` | `docs` | `style` | `refactor` | `perf` | `test` | `build` | `ci` | `chore` | `revert`

**示例**：
```
feat: 添加用户登录功能
fix(http): 修复请求超时未重试的问题
refactor: 重构配置加载逻辑
docs: 更新API使用文档
chore(deps): 升级Go依赖版本
```

## License

MIT
