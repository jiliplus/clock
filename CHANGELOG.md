# 修改日志

这个项目中，所有需要值得一提的改变，都会在这里罗列出来。

本文的格式基于[如何维护更新日志](https://keepachangelog.com/zh-CN/1.0.0/)，并且这个项目的版本号基于[语义化版本2.0.0](https://semver.org/lang/zh-CN/).

## [最新更改]

## [0.9.0] - 2020-01-30

### 安全改进

- 使用 `GoMock` 替换了部分测试代码，避免了覆盖率污染。

## [0.1.0] - 2020-01-26

### 添加

- 定义了 `Clock` 接口
- `contextSim` 结构体实现了 `context.Context` 接口
- `Set` 函数，把 `Clock` 放入上下文
- `Get` 函数，从 `Clock` 取出上下文
- 封装了一系列从上下文中取出 `Clock` 的后续操作。
- `NewRealClock` 返回 `Clock` 接口变量，其方法是对 `time` 标准库的封装。
- `NewSimulator` 返回 `*Simulator` 变量，它实现了 `Clock` 接口，并能由 `Add`，`AddOrPanic`，`Set`，`SetOrPanic` 和 `Move` 方法驱动运行。
- `*Ticker` 实现了 `*time.Ticker` 一样的功能。由 `*Simulator` 生成的 `*Ticker` 可以由其驱动。
- `*Timer` 实现了 `*time.Timer` 一样的功能。由 `*Simulator` 生成的 `*Timer` 可以由其驱动。

[最新更改]: https://github.com/jujili/clock/compare/v0.9.0...HEAD
[0.9.0]: https://github.com/jujili/clock/compare/v0.1.0...0.9.0
[0.1.0]: https://github.com/jujili/clock/compare/v0.0.0...v0.1.0

<!-- ### 添加 -->
<!-- ### 变更 -->
<!-- ### 待删除 -->
<!-- ### 已删除 -->
<!-- ### 修复 -->
<!-- ### 安全改进 -->
<!--  -->
<!-- ### Added 新添加的功能。 -->
<!-- ### Changed 对现有功能的变更。 -->
<!-- ### Deprecated 已经不建议使用，准备很快移除的功能。 -->
<!-- ### Removed 已经移除的功能。 -->
<!-- ### Fixed 对bug的修复 -->
<!-- ### Security 对安全的改进 -->
