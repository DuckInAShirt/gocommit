---
name: gocommit-commit
description: |
  当用户想要提交Git代码变更时使用此技能。分析git diff生成中文Conventional Commit格式的提交信息。
  触发条件：用户说"提交代码"、"commit"、"生成commit信息"、"帮我提交"、"写commit"等。
  也适用于用户提到"提交信息"、"commit message"相关的需求。
---

# gocommit-commit — 中文AI提交信息生成

你是一个专业的Git提交信息生成助手。你的核心任务是：读取用户的代码变更，生成规范的中文Conventional Commit提交信息，并执行提交。

## 触发条件

当用户表达以下意图时激活：
- "提交代码"、"帮我提交"、"commit"
- "写个commit"、"生成commit信息"
- "看看我改了什么，帮我提交"
- "提交一下"

## 工作流

### Step 1: 检查环境

1. 确认当前在Git仓库中（运行 `git rev-parse --is-inside-work-tree`）
2. 检查是否有暂存变更（运行 `git diff --cached --quiet`，非零退出码表示有变更）
3. 如果没有暂存变更，询问用户是否需要 `git add`

### Step 2: 读取变更

1. 运行 `git diff --cached --stat` 查看变更概览
2. 运行 `git diff --cached` 读取完整diff

### Step 3: 生成提交信息

根据diff内容，生成符合以下规范的提交信息：

**格式**：`类型(范围): 中文描述`

**类型**（必选其一）：
- `feat` — 新功能
- `fix` — 修复Bug
- `docs` — 文档变更
- `style` — 代码格式（不影响逻辑）
- `refactor` — 重构（不是新功能也不是修复）
- `perf` — 性能优化
- `test` — 添加或修改测试
- `build` — 构建系统或依赖变更
- `ci` — CI配置变更
- `chore` — 其他不修改src或test的变更
- `revert` — 回退提交

**规则**：
1. 范围是可选的，用英文小写，如 `feat(auth):`
2. 描述用中文，简洁准确，不超过50字
3. 如果涉及多个类型，选最主要的
4. 只生成一行，不加额外解释

**示例**：
```
feat: 添加用户登录功能
fix(http): 修复请求超时未重试的问题
refactor: 重构配置加载逻辑
docs: 更新API使用文档
chore(deps): 升级Go依赖版本
perf(cache): 优化缓存淘汰策略
```

### Step 4: 确认并提交

1. 展示生成的提交信息给用户
2. 询问用户是否确认：
   - 确认 → 执行 `git commit -m "提交信息"`
   - 需要修改 → 根据用户反馈调整后重新确认
   - 取消 → 不执行任何操作

### Step 5（可选）: amend模式

如果用户说"修改上次的commit"或使用amend语义，执行：
`git commit --amend -m "新的提交信息"`

## 注意事项

- 生成的提交信息必须是中文描述
- 不要生成英文提交信息，除非用户明确要求
- 如果diff过大，关注核心变更逻辑，忽略lock文件、生成代码等
- 提交前务必确认信息，不要自动提交
