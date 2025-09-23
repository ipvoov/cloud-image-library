# 贡献指南

感谢您对智能云图库项目的关注！我们欢迎任何形式的贡献，包括但不限于代码、文档、问题报告和功能建议。

## 🤝 如何贡献

### 1. Fork 项目

首先，您需要 Fork 本项目到您的 GitHub 账户：

1. 访问 [项目主页](https://github.com/your-username/cloud-picture-library)
2. 点击右上角的 "Fork" 按钮
3. 选择您的账户作为 Fork 的目标

### 2. 克隆项目

```bash
# 克隆您 Fork 的项目
git clone https://github.com/your-username/cloud-picture-library.git
cd cloud-picture-library

# 添加原始项目作为上游仓库
git remote add upstream https://github.com/original-username/cloud-picture-library.git
```

### 3. 创建分支

```bash
# 确保您在主分支
git checkout main

# 拉取最新代码
git pull upstream main

# 创建新的功能分支
git checkout -b feature/your-feature-name
```

### 4. 开发环境设置

**后端环境:**
```bash
# 安装 Go 1.23+
go version

# 安装依赖
go mod tidy

# 配置数据库
# 参考 README.md 中的数据库配置部分

# 运行测试
go test ./...
```

**前端环境:**
```bash
cd picture-frontend

# 安装 Node.js 18+
node --version

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

### 5. 代码规范

#### Go 代码规范

- 遵循 [Go 官方代码规范](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 添加必要的注释和文档

```bash
# 格式化代码
gofmt -w .

# 检查代码质量
golint ./...

# 运行测试
go test -v ./...
```

#### Vue/TypeScript 代码规范

- 遵循 [Vue 3 风格指南](https://vuejs.org/style-guide/)
- 使用 ESLint 和 Prettier 格式化代码
- 使用 TypeScript 进行类型检查

```bash
# 检查代码质量
npm run lint

# 格式化代码
npm run format

# 类型检查
npm run type-check
```

### 6. 提交代码

```bash
# 添加修改的文件
git add .

# 提交代码（使用规范的提交信息）
git commit -m "feat: add new feature for image upload"

# 推送到您的 Fork
git push origin feature/your-feature-name
```

#### 提交信息规范

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**类型 (type):**
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

**示例:**
```
feat(auth): add OAuth2 login support

- Add Google OAuth2 provider
- Add GitHub OAuth2 provider
- Update login UI to support OAuth2

Closes #123
```

### 7. 创建 Pull Request

1. 访问您 Fork 的项目页面
2. 点击 "New Pull Request" 按钮
3. 选择您的功能分支
4. 填写 PR 描述，包括：
   - 功能描述
   - 修改内容
   - 测试情况
   - 相关 Issue（如果有）

## 📋 Pull Request 模板

```markdown
## 变更描述
简要描述本次 PR 的内容和目的。

## 变更类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 文档更新
- [ ] 代码重构
- [ ] 性能优化
- [ ] 其他

## 测试情况
- [ ] 已添加单元测试
- [ ] 已进行手动测试
- [ ] 所有测试通过

## 相关 Issue
关联的 Issue: #123

## 截图（如有 UI 变更）
请提供相关截图。

## 检查清单
- [ ] 代码遵循项目规范
- [ ] 已添加必要的注释
- [ ] 已更新相关文档
- [ ] 已通过所有测试
```

## 🐛 报告问题

### Bug 报告

如果您发现了 bug，请通过 [Issues](https://github.com/your-username/cloud-picture-library/issues) 报告：

**Bug 报告模板:**
```markdown
## Bug 描述
简要描述 bug 的情况。

## 重现步骤
1. 进入 '...'
2. 点击 '...'
3. 滚动到 '...'
4. 看到错误

## 预期行为
描述您期望的正确行为。

## 实际行为
描述实际发生的情况。

## 环境信息
- OS: [e.g. Windows 10, macOS 12, Ubuntu 20.04]
- 浏览器: [e.g. Chrome 95, Firefox 94]
- 版本: [e.g. v1.0.0]

## 截图
如果适用，请添加截图。

## 附加信息
添加任何其他相关信息。
```

### 功能建议

如果您有功能建议，请使用以下模板：

```markdown
## 功能描述
简要描述您希望添加的功能。

## 使用场景
描述这个功能的使用场景和价值。

## 实现建议
如果您有实现建议，请提供。

## 附加信息
添加任何其他相关信息。
```

## 🏗️ 开发指南

### 项目架构

**后端架构:**
```
internal/
├── cmd/           # 命令行入口
├── controller/    # 控制器层（处理HTTP请求）
├── dao/          # 数据访问层（数据库操作）
├── logic/        # 业务逻辑层（核心业务逻辑）
├── middleware/   # 中间件（认证、日志等）
├── model/        # 数据模型（实体和DO）
└── service/      # 服务层（业务服务接口）
```

**前端架构:**
```
src/
├── api/          # API 接口定义
├── components/   # 可复用组件
├── pages/        # 页面组件
├── stores/       # 状态管理
├── utils/        # 工具函数
└── router/       # 路由配置
```

### 添加新功能

1. **后端开发流程:**
   ```bash
   # 1. 定义 API 接口
   # 在 api/user/v1/ 中添加新的接口定义
   
   # 2. 实现控制器
   # 在 internal/controller/ 中添加控制器方法
   
   # 3. 编写业务逻辑
   # 在 internal/logic/ 中实现业务逻辑
   
   # 4. 更新数据模型
   # 在 internal/model/ 中添加或修改数据模型
   
   # 5. 添加测试
   # 编写单元测试和集成测试
   ```

2. **前端开发流程:**
   ```bash
   # 1. 定义 API 接口
   # 在 src/api/ 中添加 API 接口
   
   # 2. 创建页面组件
   # 在 src/pages/ 中创建页面组件
   
   # 3. 创建可复用组件
   # 在 src/components/ 中创建组件
   
   # 4. 更新路由
   # 在 src/router/ 中添加路由配置
   
   # 5. 更新状态管理
   # 在 src/stores/ 中添加状态管理
   ```

### 数据库变更

如果需要修改数据库结构：

1. **创建迁移文件:**
   ```bash
   # 创建新的迁移文件
   go run main.go migrate create add_new_table
   ```

2. **编写迁移脚本:**
   ```sql
   -- 在迁移文件中添加 SQL
   CREATE TABLE new_table (
     id BIGINT PRIMARY KEY AUTO_INCREMENT,
     name VARCHAR(128) NOT NULL,
     created_at DATETIME DEFAULT CURRENT_TIMESTAMP
   );
   ```

3. **执行迁移:**
   ```bash
   # 执行迁移
   go run main.go migrate up
   ```

### 测试指南

**后端测试:**
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/logic/picture

# 运行测试并显示覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**前端测试:**
```bash
# 运行单元测试
npm run test

# 运行 E2E 测试
npm run test:e2e

# 生成覆盖率报告
npm run test:coverage
```

## 📚 文档贡献

### 文档类型

- **API 文档**: 在 `docs/API.md` 中更新
- **部署文档**: 在 `docs/DEPLOYMENT.md` 中更新
- **开发文档**: 在 `docs/DEVELOPMENT.md` 中更新
- **用户手册**: 在 `docs/USER_GUIDE.md` 中更新

### 文档规范

- 使用 Markdown 格式
- 保持文档结构清晰
- 添加必要的代码示例
- 定期更新过时信息

## 🎯 贡献类型

### 代码贡献
- 新功能开发
- Bug 修复
- 性能优化
- 代码重构

### 文档贡献
- 完善 API 文档
- 添加使用示例
- 更新部署指南
- 编写用户手册

### 测试贡献
- 添加单元测试
- 编写集成测试
- 性能测试
- 安全测试

### 设计贡献
- UI/UX 设计
- 图标设计
- 用户体验优化

## 🏆 贡献者

感谢所有为项目做出贡献的开发者！

<!-- 这里会自动生成贡献者列表 -->

## 📞 联系我们

如果您有任何问题或建议，可以通过以下方式联系我们：

- **GitHub Issues**: [项目 Issues](https://github.com/your-username/cloud-picture-library/issues)
- **邮箱**: your-email@example.com
- **QQ群**: 123456789
- **微信群**: 扫描二维码加入

## 📄 许可证

本项目采用 MIT 许可证。详情请查看 [LICENSE](LICENSE) 文件。

---

再次感谢您的贡献！您的参与让这个项目变得更好。🎉
