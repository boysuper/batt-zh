package gui

/*
#include <stdlib.h>
#include "native.h"
*/
import "C"

import "sync"

// zhTable 是英文原文到中文的对照表。
var zhTable = map[string]string{
	"batt icon": "batt 图标",
	"batt daemon not installed": "batt 后台服务未安装",
	"Your machine cannot run batt": "此设备无法运行 batt",
	"batt needs upgrade": "batt 需要升级",
	"Power Flow": "电源流向",
	"Upgrade Daemon...": "升级后台服务…",
	"Install Daemon...": "安装后台服务…",
	"Loading...": "加载中…",
	"Quick Limits": "快速设上限",
	"Set %ld%% Limit": "设上限 %ld%%",
	"One-Time Charge": "一次性充电",
	"Charge to Limit Now": "立即充到上限",
	"Charge to Full Once": "一次性充满",
	"Cancel One-Time Charge": "取消一次性充电",
	"Advanced": "高级设置",
	"Control MagSafe LED": "MagSafe 指示灯",
	"Enable": "启用",
	"Disable": "停用",
	"Always-off": "始终熄灭",
	"Prevent Idle Sleep when Charging": "充电时防止空闲睡眠",
	"Disable Charging before Sleep": "睡前停止充电",
	"Prevent System Sleep when Charging (Experimental)": "充电时阻止系统睡眠（实验性）",
	"Prevent Sleep when Adapter is Disabled": "断开电源时防止睡眠",
	"No-drain Above Native Limit (Experimental)": "高于上限不放电（实验性）",
	"Force Discharge...": "强制放电…",
	"Stop Force Discharge": "停止强制放电",
	"Indefinitely": "不限时",
	"1 Hour": "1 小时",
	"2 Hours": "2 小时",
	"4 Hours": "4 小时",
	"8 Hours": "8 小时",
	"12 Hours": "12 小时",
	"24 Hours": "24 小时",
	"2 Days": "2 天",
	"3 Days": "3 天",
	"7 Days": "7 天",
	"Auto Calibration (Experimental)...": "自动校准（实验性）…",
	"Status: Idle": "状态：空闲",
	"Start": "开始",
	"Pause": "暂停",
	"Resume": "恢复",
	"Cancel": "取消",
	"Version: ": "版本：",
	"Uninstall Daemon...": "卸载后台服务…",
	"Disable Charge Limit": "临时取消上限",
	"Quit Menubar App": "退出菜单栏",
	"Language": "语言",
	"Auto (System Language)": "自动（跟随系统）",
	"Chinese": "中文",
	"English": "English",
	"Language changed": "语言已切换",
	"Restarting batt to apply the new language.": "正在重启 batt 以应用新语言。",
	"Warning": "警告",
	"notes": "备注",
	"calibration": "校准",
	"Precautions": "注意事项",
	"1. Keep the lid open unless \"Prevent Sleep when Adapter is Disabled\" is enabled. When using that option with the lid closed, keep your Mac monitored and never place it in a bag.\n": "1. 除非已开启「断开电源时防止睡眠」，否则请保持开盖。开启该选项并合盖使用时，务必让 Mac 处于有人看管的状态，绝不能放进包里。\n",
	"2. Force Discharge cuts wall power, making your Mac run on battery power just as if it were unplugged. If the battery charge is exhausted, your Mac will shut down.": "2. 强制放电会切断墙电，让 Mac 像拔掉电源一样使用电池。电量耗尽后，Mac 会自动关机。",
	"\n3. An indefinite force discharge will continue until you stop it. Prefer a timed duration unless you specifically need indefinite discharge.": "\n3. 不限时的强制放电会一直持续，直到你手动停止。除确有需要，建议选择定时时长。",
	"Start Auto Calibration?": "开始自动校准？",
	"This will:\n": "流程如下：\n",
	"1. Discharge to 15% by default.\n": "1. 默认放电到 15%。\n",
	"2. Charge to 100%.\n": "2. 充到 100%。\n",
	"3. Hold at full charge (for 2 hours by default).\n": "3. 满电保持（默认 2 小时）。\n",
	"4. Discharge back to previous charge limit.\n\n": "4. 放回原来的充电上限。\n\n",
	"NOTES:\n": "注意：\n",
	"• batt will prevent idle sleep until calibration finishes, is cancelled, or fails.\n": "• 校准结束、被取消或失败之前，batt 会阻止空闲睡眠。\n",
	"• You can pause or cancel anytime from the menu.\n": "• 随时可以在菜单里暂停或取消。\n",
	"• Highly recommend keeping your Mac connected to power throughout the process to prevent the battery level from dropping below the threshold without timely charging.\n": "• 强烈建议全程接着电源，避免电量掉到阈值以下却没能及时充电。\n",
	"• Closing the lid or explicitly choosing Sleep can still force sleep unless \"Prevent Sleep when Adapter is Disabled\" is enabled. When using that option with the lid closed, keep your Mac monitored and never place it in a bag.": "• 除非开启「断开电源时防止睡眠」，合盖或手动选择睡眠仍会强制睡眠。开启该选项并合盖时，请保持有人看管，切勿放进包里。",
	"Prevent Sleep when Adapter is Disabled?": "断开电源时防止睡眠？",
	"This enables macOS’s global no-sleep setting while adapter input is disabled. Your Mac may remain awake with the lid closed, overheat in a bag, or drain its battery. Use this only when the Mac is monitored. Continue?": "这会在断开电源输入期间启用 macOS 的全局「不休眠」设置。你的 Mac 可能合盖后仍持续运行、放进包里过热，或耗尽电池。仅在有人看管时使用。是否继续？",
	"Enable No-drain?": "启用「不放电」模式？",
	"batt must restart powerd when this policy changes. batt waits until other applications stop preventing sleep before restarting, but a restart can still interrupt power-related activity. The regular native limit is used until the policy can be applied. Continue?": "更改此策略需要重启 powerd。batt 会等其他应用停止阻止睡眠后再重启，但重启仍可能打断与电源相关的活动。在策略生效前会继续使用常规原生上限。是否继续？",
	"Continue": "继续",
	"Current Limit: %d%%": "当前上限：%d%%",
	"State: ": "状态：",
	"State: Error": "状态：错误",
	"Auto Calibration (Experimental) Paused...": "自动校准（实验性）已暂停…",
	"Auto Calibration (Experimental) In Progress...": "自动校准（实验性）进行中…",
	"Status: Discharging (%d%% → %d%%)": "状态：放电中（%d%% → %d%%）",
	"Status: Charging (%d%% → 100%%)": "状态：充电中（%d%% → 100%%）",
	"Status: Holding (%02d:%02d:%02d left)": "状态：保持中（剩 %02d:%02d:%02d）",
	"Status: In Progress": "状态：进行中",
	"Failed to get executable path": "获取程序路径失败",
	"Installation failed": "安装失败",
	"Failed to start app at boot": "设置开机自启失败",
	"Failed to set limit": "设置充电上限失败",
	"Charge to %d%% Now": "立即充至 %d%%",
	"Charging to %d%%…": "正在充至 %d%%…",
	"Charging to %d%%, now %d%%": "正在充至 %d%%，当前 %d%%",
	"Failed to set MagSafe LED control": "设置 MagSafe 指示灯失败",
	"Failed to set prevent idle sleep": "设置「充电时防止空闲睡眠」失败",
	"Failed to set prevent system sleep": "设置「充电时阻止系统睡眠」失败",
	"Failed to set prevent sleep on adapter disable": "设置「断开电源时防止睡眠」失败",
	"Failed to set no-drain": "设置「不放电」模式失败",
	"No-drain": "不放电",
	"Failed to set force discharge": "设置强制放电失败",
	"Failed to stop force discharge": "停止强制放电失败",
	"Failed to start calibration": "开始校准失败",
	"Failed to pause calibration": "暂停校准失败",
	"Failed to resume calibration": "恢复校准失败",
	"Failed to cancel calibration": "取消校准失败",
	"Failed to uninstall daemon": "卸载后台服务失败",
	"Failed to unregister login item": "取消登录项失败",
	"Failed to disable charge limit": "取消充电上限失败",
	"Restoring %d%% limit…": "正在恢复 %d%% 上限…",
	"Restores to %d%% in %s": "%s 后恢复到 %d%%",
	"Failed to start charging to the limit": "启动「充到上限」失败",
	"Failed to start charging to 100%": "启动「充到 100%」失败",
	"Failed to cancel the one-time charge": "取消一次性充电失败",
	"Failed to get charging status": "获取充电状态失败",
	"macOS may charge to 100% for battery maintenance.": "macOS 为维护电池可能会充到 100%。",
	"A temporary charge-limit disable is in progress. To cancel the schedule, set a charge limit.": "临时取消上限进行中。要取消计划，重新设置一个充电上限即可。",
	"A temporary force discharge is in progress. To cancel the schedule and restore wall power, choose \"Stop Force Discharge\".": "临时强制放电进行中。要取消计划并恢复墙电，请选择「停止强制放电」。",
	"Quit the batt menubar app.": "退出 batt 菜单栏程序。",
	"Power flow data is updated every 60 seconds.": "电源数据每 60 秒刷新一次。",
	"Your batt daemon is not compatible with this client version and needs to be upgraded. This is usually caused by a new client version that requires a new daemon version. You can upgrade the batt daemon by running this command.": "你的 batt 后台服务与当前客户端版本不匹配，需要升级。通常是因为新版客户端要求新版后台服务。你可以在终端执行命令来升级后台服务。",
	"Install the batt daemon. batt daemon is a component that controls charging. You must enter your password to install it because controlling charging is a privileged action.": "安装 batt 后台服务。后台服务负责控制充电，属于特权操作，所以安装时需要输入你的密码。",
	"Charge the battery once without changing your charge limit.\n\n\"Charge to <limit>% Now\" starts charging right away, even when the charge sits between the lower and the upper limit, and stops at the limit.\n\n\"Charge to Full Once\" charges to 100% and then restores your limit automatically.\n\nA one-time charge survives a restart. Cancel it here, or set a charge limit to cancel it.": "在不改变充电上限的情况下，让电池再充一次。\n\n「立即充到上限」会马上开始充电——即使当前电量介于下限与上限之间——充到上限即停。\n\n「一次性充满」会充到 100%%，随后自动恢复原来的上限。\n\n一次性充电在重启后依然存在。在此取消即可，或重新设置一个上限也会取消它。",
	"Let batt control MagSafe LED to reflect the charging state of your MacBook (or force it off).\n\nNote that you must have a MagSafe LED on your MacBook to use this feature.": "让 batt 控制 MagSafe 指示灯，用灯光反映 MacBook 的充电状态（也可强制熄灭）。\n\n注意：MacBook 必须带有 MagSafe 指示灯才能使用此功能。",
	"Enable MagSafe LED control. The LED will reflect charging status:\n\n- Green: Charging is stopped.\n- Orange: Charging is in progress.\n- Off: Forced off, or temporarily after wake with active charging control.": "启用 MagSafe 指示灯控制。灯光反映充电状态：\n\n- 绿色：已停止充电。\n- 橙色：正在充电。\n- 熄灭：强制熄灭，或充电控制生效期间唤醒后短暂熄灭。",
	"Disable MagSafe LED control. The LED will stay in its default state (mostly orange).": "停用 MagSafe 指示灯控制。灯光会保持系统默认状态（多数时候是橙色）。",
	"Force the MagSafe LED to stay off regardless of charging state.": "无论充电状态如何，都强制 MagSafe 指示灯保持熄灭。",
	"Set whether to prevent idle sleep during a charging session.\n\nDue to macOS limitations, batt will be paused when your computer goes to sleep. As a result, when you are in a charging session and your computer goes to sleep, there is no way for batt to stop charging (since batt is paused by macOS) and the battery will charge to 100%. This option, together with \"Disable Charging before Sleep\", will prevent this from happening.\n\nThis option tells macOS NOT to go to sleep when the computer is in a charging session, so batt can continue to work until charging is finished. Note that it will only prevent **idle** sleep, when 1) charging is active 2) battery charge limit is enabled. So your computer can go to sleep as soon as a charging session is completed.\n\nHowever, this options does not prevent manual sleep (limitation of macOS). For example, if you manually put your computer to sleep (by choosing the Sleep option in the top-left Apple menu) or close the lid, batt will still be paused and the issue mentioned above will still happen. This is where \"Disable Charging before Sleep\" comes in.": "设置充电期间是否阻止空闲睡眠。\n\n受 macOS 限制，电脑一旦睡眠，batt 就会被暂停。因此在充电过程中如果电脑进入睡眠，batt 无法继续停止充电（已被系统暂停），电池会被充到 100%。本选项配合「睡前停止充电」可以避免这个问题。\n\n本选项让 macOS 在充电期间不要进入睡眠，这样 batt 能一直工作到充电结束。注意它只阻止**空闲**睡眠，且需同时满足：1）正在充电；2）已启用充电上限。所以充电一结束，电脑就可以照常睡眠。\n\n但它无法阻止手动睡眠（macOS 的限制）。例如通过左上角苹果菜单选择「睡眠」，或合上盖子，batt 仍会被暂停，上述问题照样发生。这时就需要「睡前停止充电」。",
	"Set whether to disable charging before sleep if charge limit is enabled.\n\nAs described in \"Prevent Idle Sleep when Charging\", batt will be paused by macOS when your computer goes to sleep, and there is no way for batt to continue controlling battery charging. This option will disable charging just before sleep, so your computer will not overcharge during sleep, even if the battery charge is below the limit.": "设置启用充电上限时，是否在睡前先停止充电。\n\n正如「充电时防止空闲睡眠」所说，电脑睡眠时 batt 会被 macOS 暂停，无法继续控制充电。本选项会在临睡前先关闭充电，这样即使电量低于上限，睡眠期间也不会被过度充电。",
	"Set whether to prevent system sleep during a charging session (experimental).\n\nThis option tells macOS to create power assertion, which prevents sleep, when all conditions are met:\n\n1) charging is active\n2) battery charge limit is enabled\n3) computer is connected to charger.\nSo your computer can go to sleep as soon as a charging session is completed / charger disconnected.\n\nDoes similar thing to prevent-idle-sleep, but works for manual sleep too.\n\nNote: please disable disable-charging-pre-sleep and prevent-idle-sleep, while this feature is in use": "设置充电期间是否阻止系统睡眠（实验性）。\n\n本选项让 macOS 在以下条件全部满足时创建电源断言（阻止睡眠）：\n\n1）正在充电\n2）已启用充电上限\n3）已连接电源适配器\n因此充电结束或拔掉电源后，电脑马上就能睡眠。\n\n作用与「充电时防止空闲睡眠」类似，但对手动睡眠同样有效。\n\n注意：使用本功能时，请关闭「睡前停止充电」和「充电时防止空闲睡眠」。",
	"Set whether to keep Clamshell mode alive while adapter input is disabled by temporarily disabling system sleep.\n\nThis prevents all sleep—not just idle sleep—while adapter input is disabled, allowing Force Discharge and Auto Calibration to continue with the lid closed. When adapter input returns, batt automatically restores the prior sleep setting.\n\nWARNING: Your Mac may remain awake with the lid closed, overheat in a bag, or drain its battery to 0% if unplugged and left unattended.": "设置断开电源输入期间，是否通过临时关闭系统睡眠来维持合盖模式运行。\n\n它会在断开电源时阻止所有睡眠（不只是空闲睡眠），让强制放电和自动校准可以在合盖状态下继续。电源恢复后，batt 会自动还原之前的睡眠设置。\n\n警告：合盖时 Mac 可能持续运行，放进包里会过热；若未接电源且无人看管，电池可能耗尽至 0%。",
	"Experimental macOS 27 native-mode option. Hold wall power above the limit instead of draining the battery.\n\nChanging this policy restarts powerd. batt waits for other applications' sleep assertions to end first and uses the regular native limit meanwhile.": "macOS 27 原生模式的实验性选项：电量高于上限时继续使用墙电，而不是放电。\n\n更改此策略会重启 powerd。batt 会先等其他应用的防睡眠断言结束，期间仍使用常规的原生上限。",
	"Cut power from the wall and run from the battery, either indefinitely or for a selected duration. After a temporary force discharge, batt automatically restores the power adapter.\n\nNOTE: if you are using Clamshell mode (using a Mac laptop with an external monitor and the lid closed), *cutting power will cause your Mac to go to sleep* unless \"Prevent Sleep when Adapter is Disabled\" is enabled. This is a limitation of macOS. There are ways to prevent this, but it is not recommended for most users.": "切断墙电、改用电池供电，可选不限时或指定时长。临时强制放电结束后，batt 会自动恢复电源适配器。\n\n注意：若使用合盖模式（笔记本接外接显示器并合上盖子），除非开启「断开电源时防止睡眠」，否则*断电会让 Mac 进入睡眠*。这是 macOS 的限制。有办法绕开，但不建议普通用户这么做。",
	"Calibration helps you calibrate your battery by automatically discharging and charging it according to best practices.\n\nbatt prevents idle sleep for the entire calibration session. Closing the lid or explicitly choosing Sleep can still force sleep.\n\nIt's recommended to run the calibration process every few months.": "校准会按最佳实践自动放电、充电，帮你把电池校准准确。\n\n整个校准期间 batt 会阻止空闲睡眠。合盖或手动选择睡眠仍可能强制睡眠。\n\n建议每隔几个月做一次校准。",
	"Uninstall the batt daemon. This will remove the batt daemon from your system. You must enter your password to uninstall it.\n\nAfter uninstalling the batt daemon, no charging control will be present on your system and your Mac will charge to 100% as normal. The menubar app will still be present, but all options will be disabled. You can remove the menubar app by moving it to the trash.": "卸载 batt 后台服务，会把它从系统中移除。卸载需要输入你的密码。\n\n卸载后系统中不再有任何充电控制，Mac 会像平常一样充到 100%。菜单栏程序仍然存在，但所有选项都会变灰，把它拖进废纸篓即可删除。",
	"Disable the battery charge limit and let your Mac charge to 100%, either indefinitely or for a selected duration. After a temporary disable, batt automatically restores your current limit.": "临时取消充电上限，让 Mac 充到 100%，可选不限时或指定时长。临时取消结束后，batt 会自动恢复原来的上限。",
	"Quit the batt menubar app, but keep the batt daemon running.\n\nSince the batt daemon is still running, batt can continue to control charging. This is useful if you don't want the menubar icon to show up, but still want to use batt. When the client is not running, you can change batt settings using the command line interface (batt). To prevent the menubar app from starting at login, you can remove it in System Settings -> General -> Login Items & Extensions -> remove batt.app from the list (do NOT remove the batt daemon).\n\nIf you want batt to stop controlling charging, you can use the \"Disable Charge Limit\" command. To uninstall, you can use the \"Uninstall Daemon\" command in the Advanced menu.": "退出 batt 菜单栏程序，但保留后台服务继续运行。\n\n后台服务还在，batt 依然能控制充电。如果不想让菜单栏出现图标、但仍想使用 batt，这个选项很有用。客户端不运行时，可以用命令行（batt）修改设置。若不想开机自启，可在「系统设置 → 通用 → 登录项与扩展」里移除 batt.app（不要移除后台服务）。\n\n若想让 batt 彻底停止控制充电，可使用「临时取消上限」命令；要卸载请用「高级设置」里的「卸载后台服务」。",
	"Not Charging": "未充电",
	"Charging": "充电中",
	"Discharging": "放电中",
	"Full": "满电",
	"Will Charge Soon": "即将充电",
	"Battery full; waiting for macOS": "电池已满，等待 macOS",
	"Limit paused by macOS": "上限被 macOS 暂停",
	"Waiting for macOS to apply limit": "等待 macOS 应用上限",
	"Status: Discharging to previous limit...": "状态：正在放电回原上限…",
	"Status: Restoring settings...": "状态：正在恢复设置…",
	"Status: Error - ": "状态：错误 - ",
	"Status: Error": "状态：错误",
	"Restoring power adapter…": "正在恢复电源适配器…",
	"Calibration": "校准",
	"%dd": "%d 天",
	"%dh": "%d 小时",
	"%dm": "%d 分钟",
	"System": "系统",
	"Adapter": "适配器",
	"Battery": "电池",
	"Battery full; charge limit paused until macOS clears the override.": "电池已满，macOS 清除覆盖前上限暂时不生效。",
	"macOS has not applied the limit yet and will do so once the battery level changes or charging stops, usually within a few minutes.": "macOS 尚未应用上限，通常会在电量变化或停止充电后的几分钟内生效。",
	"Cut power from the wall and run from the battery, either indefinitely or for a selected duration. After a temporary force discharge, batt automatically restores the power adapter.": "切断墙电、改用电池供电，可选不限时或指定时长。临时强制放电结束后，batt 会自动恢复电源适配器。",
	"1. Keep the lid open unless \"Prevent Sleep when Adapter is Disabled\" is enabled. When using that option with the lid closed, keep your Mac monitored and never place it in a bag.\n2. Force Discharge cuts wall power, making your Mac run on battery power just as if it were unplugged. If the battery charge is exhausted, your Mac will shut down.": "1. 除非已开启「断开电源时防止睡眠」，否则请保持开盖。开启该选项并合盖使用时，请随时留意 Mac，切勿放进包里。\n2. 强制放电会切断墙电，让 Mac 像拔掉电源一样靠电池运行。电量耗尽时 Mac 会关机。",
	"This will:\n1. Discharge to 15% by default.\n2. Charge to 100%.\n3. Hold at full charge (for 2 hours by default).\n4. Discharge back to previous charge limit.\n\nNOTES:\n• batt will prevent idle sleep until calibration finishes, is cancelled, or fails.\n• You can pause or cancel anytime from the menu.\n• Highly recommend keeping your Mac connected to power throughout the process to prevent the battery level from dropping below the threshold without timely charging.\n• Closing the lid or explicitly choosing Sleep can still force sleep unless \"Prevent Sleep when Adapter is Disabled\" is enabled. When using that option with the lid closed, keep your Mac monitored and never place it in a bag.": "将执行：\n1. 默认放电至 15%。\n2. 充至 100%。\n3. 保持满电（默认 2 小时）。\n4. 放电回原来的充电上限。\n\n注意：\n• 校准完成、取消或失败之前，batt 会阻止系统空闲睡眠。\n• 可随时在菜单里暂停或取消。\n• 强烈建议全程接着电源，避免电量掉到阈值以下却没能及时充电。\n• 合盖或手动睡眠仍会强制睡眠，除非已开启「断开电源时防止睡眠」。开启该选项并合盖使用时，请随时留意 Mac，切勿放进包里。",
}


var (
	langMu    sync.Mutex
	langReady bool
	langCode  string
	cstrMu    sync.Mutex
	cstrMap  = map[string]*C.char{}
)

func languageMode() string {
	mode := C.GoString(C.batt_get_language())
	if mode == "" {
		return "auto"
	}
	return mode
}

func resolveLanguage() string {
	langMu.Lock()
	defer langMu.Unlock()
	if langReady {
		return langCode
	}
	mode := languageMode()
	langReady = true
	if mode == "zh" || mode == "en" {
		langCode = mode
		return langCode
	}
	if C.GoString(C.batt_system_language()) == "zh" {
		langCode = "zh"
	} else {
		langCode = "en"
	}
	return langCode
}

// resetLanguage 让下一次取词重新读取语言设置。
func resetLanguage() {
	langMu.Lock()
	langReady = false
	langMu.Unlock()
}

func translate(en string) string {
	if resolveLanguage() != "zh" {
		return en
	}
	if zh, ok := zhTable[en]; ok {
		return zh
	}
	return en
}

// tr 用于 Go 侧界面文案。
func tr(en string) string { return translate(en) }

//export battTranslate
func battTranslate(english *C.char) *C.char {
	if english == nil {
		return nil
	}
	en := C.GoString(english)
	out := translate(en)
	cstrMu.Lock()
	defer cstrMu.Unlock()
	if c, ok := cstrMap[out]; ok {
		return c
	}
	c := C.CString(out)
	cstrMap[out] = c
	return c
}

func (c *menuController) syncLanguageChecks() {
	mode := languageMode()
	c.menu.setChecked(itemLanguageAuto, mode == "auto")
	c.menu.setChecked(itemLanguageChinese, mode == "zh")
	c.menu.setChecked(itemLanguageEnglish, mode == "en")
}

// setLanguage 切换语言并就地重建菜单，不退出进程。
func (c *menuController) setLanguage(mode string) {
	C.batt_set_language(C.CString(mode))
	resetLanguage()
	if c.menu != nil {
		c.menu.rebuild()
		c.syncLanguageChecks()
		c.refreshOnOpen()
	}
}
