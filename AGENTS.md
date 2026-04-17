# AGENTS.md - gocommit 项目上下文

## 项目概述

**gocommit** — AI 驱动的中文 Git 提交信息生成工具

GitHub: https://github.com/DuckInAShert/gocommit

### 核心理念
- 中文优先的 Conventional Commits 提交信息生成
- CLI + Agent Skill 双形态交付（命令行工具 + SKILL.md）
- Go 单二进制，无运行时依赖

### 项目愿景
从一个 commit 工具出发，逐步扩展为中文开发者的 Git 智能助手：
- v1.0: AI Commit（当前）
- v2.0: AI Git（+ branch命名、PR描述、changelog）
- v3.0: AI Code Review
- v4.0: Git Copilot（完整Git工作流AI伙伴）

## 项目结构

```
github-star/
├── cmd/gocommit/main.go       # CLI入口（cobra）
├── internal/
│   ├── git/diff.go            # git命令封装（diff、commit、stage等）
│   ├── ai/client.go           # OpenAI兼容API调用 + 中文commit prompt
│   ├── commit/commit.go       # 主流程编排（生成→确认→提交）
│   └── config/config.go       # 配置管理（~/.gocommit.json）
├── skill/gocommit-commit/
│   └── SKILL.md               # Agent Skill文件（opencode/claude兼容）
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

## 技术栈

- Go 1.25
- github.com/spf13/cobra — CLI框架
- OpenAI 兼容 API — 支持任何兼容端点（OpenAI、DeepSeek、Ollama、OpenCode Go等）

## 已实现功能（v0.1 MVP）

- git diff --cached 读取（排除lock文件等噪声）
- OpenAI兼容API调用，中文conventional commit生成
- 交互确认流程（y=确认 / e=编辑 / r=重试 / n=取消）
- config管理（api_key / base_url / model），支持环境变量覆盖
- cobra CLI：--all / --amend / --yes / --dry-run
- gocommit-commit Agent Skill文件
- 已编译安装到 /usr/local/bin/gocommit

## 当前配置

- API: OpenCode Go (https://opencode.ai/zen/go/v1)
- Model: kimi-k2.5
- 已验证可用，生成效果良好

## 待做 / 计划

### 协作者想练习的部分
- internal/git/diff.go — 想加 `GetNumStat()` 函数，解析 `git diff --cached --numstat`
- cmd/gocommit/main.go — 想加 `gocommit diff` 子命令，美观展示暂存变更概览
- cmd/gocommit/main.go — 想加 `gocommit setup` 交互式配置命令
- 对 TUI 感兴趣（charmbracelet/lipgloss 可能是下一步）

### 功能扩展方向
- 多候选消息生成（--generate N）
- gitmoji 支持
- git hook 集成
- 多语言支持（当前中文优先）

## 协作方式

- 项目所有者：DuckInAShert（xinranzhao）
- 协作模式：一起讨论设计，分工写代码
- 所有者想自己练手的模块：git diff、cobra CLI
- AI助手负责：AI调用模块、config管理、prompt调优、代码review

## 竞品参考

| 项目 | Stars | 语言 | 差异点 |
|------|-------|------|--------|
| Nutlope/aicommits | 8.9k | TypeScript | 最火，英文为主 |
| di-sukharev/opencommit | 7.2k | JavaScript | 功能最丰富 |
| coder/aicommit | 185 | Go | 简洁，Go生态最大 |
| dinoDanic/diny | 120 | Go | 免API Key，TUI |
| tfkhdyt/geminicommit | 214 | Go | 仅Gemini |

gocommit 的差异化：中文优先 + Go + Skill双形态 + OpenAI兼容

## 常用命令

```bash
# 编译安装
GOPROXY=https://goproxy.cn,direct go build -o /usr/local/bin/gocommit ./cmd/gocommit/

# 配置
gocommit config api_key=xxx base_url=https://opencode.ai/zen/go/v1 model=kimi-k2.5

# 测试
gocommit --dry-run
gocommit -a --dry-run   # 自动stage + 仅预览
```
