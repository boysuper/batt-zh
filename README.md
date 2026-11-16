# batt-zh · batt 社区汉化版

[batt](https://github.com/charlie0129/batt) 是一个 macOS 菜单栏电池充电管理工具（限制充电上限、一次性充满、强制放电、自动校准）。
本项目在 **batt v0.10.0 源码**基础上新增了 **中/英双语界面**，并为图标加上红色「中文版」横幅。

> 非官方汉化，与原作者无隶属关系。上游代码与项目版权归 [@charlie0129](https://github.com/charlie0129) 所有，遵循 **GPL-2.0**。

---

## 功能

- **运行时中英切换**：高级设置 → 语言 → 自动（跟随系统）/ 中文 / English
- **切换不重启**：就地重建菜单（`batt_menu_rebuild`），进程不退出，不会出现「闪退重启」的观感
- **全量汉化**：菜单项、状态行、确认弹窗、18 条悬停长提示、电源三行标签（系统 / 适配器 / 电池）
- **图标**：原图标右上角叠加红色「中文版」横幅（10 个尺寸全覆盖）
- **不改底层**：充电控制、SMC 逻辑、daemon 行为、版本号（0.10.0）、bundle id（`cc.chlc.batt`）均保持不变

## 安装

从 [Releases](https://github.com/boysuper/batt-zh/releases) 下载 `batt-0.10.0-双语版.dmg`，卷内有两个文件夹：

| 文件夹 | 内容 |
| --- | --- |
| `batt-中文版（汉化 i18n）` | 本项目的汉化构建（ad-hoc 签名） |
| `batt-原版（官方 0.10.0）` | 上游原版，用于对照 / 回滚 |

把需要的那一个 `batt.app` 拖进「应用程序」即可（两者同名，只能装一个）。

### Gatekeeper 提示（重要）

汉化版使用 **ad-hoc 临时签名**，没有 Apple Developer ID 签名与公证。在别人的 Mac 上首次打开会被拦下：

> 「无法打开 batt.app，因为无法验证开发者」

放行方式（任选其一）：

1. 右键 `batt.app` → **打开** → 再点「打开」
2. 系统设置 → 隐私与安全性 → 仍要打开
3. 终端：`xattr -dr com.apple.quarantine /Applications/batt.app`

要做到「双击直接打开、无任何提示」，只能购买 Apple Developer Program（99 美元/年）后用 Developer ID 证书签名并公证（notarization）。

## 从源码构建

需要 macOS 13+、Go 1.24、Xcode Command Line Tools。

```bash
make app VERSION=v0.10.0
# 产物：bin/batt.app
```

## 语言实现

- 词典：`pkg/gui/i18n.go`（英文原文为 key、中文为 value，150+ 条）
- 桥接：`pkg/gui/native_bridge.m`（`batt_get_language` / `batt_set_language` / `batt_menu_rebuild` / `battTranslate`）
- 持久化：`NSUserDefaults` 键 `BattLanguage`，值 `auto` / `zh` / `en`
- 菜单：`pkg/gui/native_menu.m` 中「高级设置 → 语言」子菜单

## 与上游升级的关系

本项目**没有改版本号**（仍是 0.10.0），bundle id 也与上游一致，因此：

- `brew upgrade batt` 或安装官方新版会**覆盖**掉汉化（二进制被替换，翻译层随之消失）
- 上游发新版后，需要把本仓库的改动重新应用到新版本源码再编译（改动集中在 `pkg/gui/`，冲突面很小）

## 许可与致谢

- 上游项目：[charlie0129/batt](https://github.com/charlie0129/batt)，作者 **@charlie0129**，许可 **GPL-2.0**
- 本仓库同样以 **GPL-2.0** 发布，详见 [LICENSE](./LICENSE)
- 汉化与图标改动由本仓库维护者完成；如有更好译法，欢迎 Issue / PR

---

## English

An unofficial Chinese localization of [batt](https://github.com/charlie0129/batt) v0.10.0 (macOS menubar charge limiter).

- Runtime language switch: **Advanced → Language → Auto (System) / Chinese / English**
- Switching rebuilds the menu **in place** — the process never quits, so there is no "crash & relaunch" feeling
- Full localization: menu items, status lines, alert dialogs, all 18 tooltips, power rows (System / Adapter / Battery)
- App icon gets a red "中文版" badge at the top-right
- Charging control, SMC logic, daemon behavior, version (0.10.0) and bundle id are untouched

Build: `make app VERSION=v0.10.0` → `bin/batt.app`

License: **GPL-2.0**, original work © [@charlie0129](https://github.com/charlie0129).
