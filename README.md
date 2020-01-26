<!-- markdownlint-disable MD041 -->
<h1 align="center">:alarm_clock: Clock</h1>
<p align="center">
<!--  -->
<a href="https://github.com/jujili/clock/releases"> <img src="https://img.shields.io/github/v/tag/jujili/clock?include_prereleases&sort=semver" alt="Release" title="Release"></a>
<!--  -->
<a href="https://www.travis-ci.org/jujili/clock"><img src="https://www.travis-ci.org/jujili/clock.svg?branch=master"/></a>
<!--  -->
<a href="https://codecov.io/gh/jujili/clock"><img src="https://codecov.io/gh/jujili/clock/branch/master/graph/badge.svg"/></a>
<!--  -->
<a href="https://goreportcard.com/report/github.com/jujili/clock"><img src="https://goreportcard.com/badge/github.com/jujili/clock" alt="Go Report Card" title="Go Report Card"/></a>
<!--  -->
<a href="http://godoc.org/github.com/jujili/clock"><img src="https://img.shields.io/badge/godoc-clock-blue.svg" alt="Go Doc" title="Go Doc"/></a>
<!--  -->
<br/>
<!--  -->
<a href="https://github.com/jujili/clock/blob/master/CHANGELOG.md"><img src="https://img.shields.io/badge/Change-Log-blueviolet.svg" alt="Change Log" title="Change Log"/></a>
<!--  -->
<a href="https://golang.google.cn"><img src="https://img.shields.io/github/go-mod/go-version/jujili/clock" alt="Go Version" title="Go Version"/></a>
<!--  -->
<a href="https://github.com/aQuaYi/jili/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="MIT License" title="MIT License"/></a>
<!--  -->
<br/>
<!--  -->
<a target="_blank" href="//shang.qq.com/wpa/qunwpa?idkey=7f61280435c41608fb8cb96cf8af7d31ef0007c44b223c9e3596ce84dec329bc"><img border="0" src="https://img.shields.io/badge/QQ%20群-23%2053%2000%2093-blue.svg" alt="jili交流QQ群:23530093" title="jili交流QQ群:23530093"></a>
<!--  -->
<a href="https://mp.weixin.qq.com/s?__biz=MzA4MDU4NDI5Mw==&mid=2455230332&idx=1&sn=8086c43e259b0012596ed63d6ecd7d10&chksm=88017c76bf76f5604f2f3280ffd96029b5ccaf99db48d18066d3e3bc9bc8a2e1a05de1a3225f&mpshare=1&scene=1&srcid=&sharer_sharetime=1578553397373&sharer_shareid=5ce52651949258759d82d1bf31b455b5#rd"><img src="https://img.shields.io/badge/微信公众号-jujili-success.svg" alt="微信公众号：jujili" title="微信公众号：jujili"/></a>
<!--  -->
<a href="https://zhuanlan.zhihu.com/jujili"><img src="https://img.shields.io/badge/知乎专栏-jili-blue.svg" alt="知乎专栏：jili" title="知乎专栏：jili"/></a>
<!--  -->
</p>

clock 中的 `realClock` 和 `simulator` 结构体，都实现了 `clock.Clock` 接口。需要实际时间时，使用前者；需要人为地操纵时间的场合（比如：测试）中，使用后者。

- [总体思路](#%e6%80%bb%e4%bd%93%e6%80%9d%e8%b7%af)
- [安装与更新](#%e5%ae%89%e8%a3%85%e4%b8%8e%e6%9b%b4%e6%96%b0)
- [真实的 Clock](#%e7%9c%9f%e5%ae%9e%e7%9a%84-clock)
- [模拟的 Clock](#%e6%a8%a1%e6%8b%9f%e7%9a%84-clock)

## 总体思路

在 `time` 和 `context` 标准库中，有一些函数和方法依赖于系统的当前时间。[`Clock`](https://github.com/jujili/clock/blob/master/interface.go#L13) 接口就是由这些内容组成。

clock 模块分为真实与虚拟的两个部分。真实的部分是以 time 和 context 标准库为基础，封装成了 Clock 接口。
虚拟部分也实现了 Clock 接口，只是这一部分可以人为的操控时间的改变。

## 安装与更新

在命令行中输入以下内容，可以获取到最新版

```shell
go get -u github.com/jujili/clock
```

## 真实的 Clock

```go
c := clock.NewRealClock()

// 输出操作系统的当前时间
fmt.Println(c.Now())
```

真实的 Clock 就是把 `time` 和 `context` 标准库中的相关函数，封装成了 `realClock` 的方法。

## 模拟的 Clock

```go
now,_ := time.Parse("06-01-02", "20-01-26")
c := clock.NewSimulator(now)

// 输出: 2020-01-26 00:00:00 +0000 UTC
fmt.Println(c.Now())
```

`*Simulator` 和 `contextSim` 虽然是模拟的。但是，实现了与 `time` 和 `context` 标准库中同名函数**一样的行为**。
