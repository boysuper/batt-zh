#import "native_internal.h"

static const NSTimeInterval BattMenuUpdateInterval = 60.0;

static NSMenuItem *ActionItem(BattMenuController *controller,
                              NSString *title,
                              NSString *key,
                              BattMenuItem identifier) {
    NSMenuItem *item = [[[NSMenuItem alloc] initWithTitle:title
                                                  action:@selector(menuAction:)
                                           keyEquivalent:key] autorelease];
    item.target = controller;
    item.tag = identifier;
    [controller rememberItem:item as:identifier];
    return item;
}

static NSMenuItem *DisplayItem(BattMenuController *controller,
                               NSString *title,
                               BattMenuItem identifier,
                               BOOL enabled) {
    NSMenuItem *item = [[[NSMenuItem alloc] initWithTitle:title
                                                  action:@selector(noop:)
                                           keyEquivalent:@""] autorelease];
    item.target = controller;
    item.tag = identifier;
    item.enabled = enabled;
    [controller rememberItem:item as:identifier];
    return item;
}

static NSMenu *AddSubmenu(BattMenuController *controller,
                          NSMenu *parent,
                          NSString *title,
                          BattMenuItem identifier) {
    NSMenu *submenu = [[[NSMenu alloc] initWithTitle:title] autorelease];
    submenu.autoenablesItems = NO;
    NSMenuItem *item = [[[NSMenuItem alloc] initWithTitle:title
                                                  action:nil
                                           keyEquivalent:@""] autorelease];
    item.submenu = submenu;
    [controller rememberItem:item as:identifier];
    [parent addItem:item];
    return submenu;
}

@implementation BattMenuController

- (instancetype)initWithHandle:(uintptr_t)handle version:(NSString *)version {
    self = [super init];
    if (self) {
        _handle = handle;
        _items = [[NSMutableDictionary alloc] init];
        _appVersion = [version copy];
        _menu = [[NSMenu alloc] initWithTitle:@"batt"];
        _menu.autoenablesItems = NO;
        _menu.delegate = self;
        _statusItem = [[[NSStatusBar systemStatusBar]
            statusItemWithLength:NSVariableStatusItemLength] retain];

        BattBuildMenu(self, version);
        BattApplyTooltips(self);
        _statusItem.menu = _menu;
        [self setStatusIconInstalled:NO capable:NO needsUpgrade:NO];
    }
    return self;
}

- (void)dealloc {
    [_timer invalidate];
    [_timer release];
    _menu.delegate = nil;
    [[NSStatusBar systemStatusBar] removeStatusItem:_statusItem];
    [_statusItem release];
    [_menu release];
    [_items release];
    [super dealloc];
}

- (NSMenuItem *)item:(BattMenuItem)item {
    return [_items objectForKey:[NSNumber numberWithInteger:item]];
}

- (void)rememberItem:(NSMenuItem *)item as:(BattMenuItem)identifier {
    [_items setObject:item forKey:[NSNumber numberWithInteger:identifier]];
}

- (void)menuAction:(NSMenuItem *)sender {
    BattMenuItem identifier = (BattMenuItem)sender.tag;
    switch (identifier) {
        case BattItemMagSafeEnabled:
        case BattItemMagSafeDisabled:
        case BattItemMagSafeAlwaysOff:
            [self item:BattItemMagSafeEnabled].state = identifier == BattItemMagSafeEnabled
                ? NSControlStateValueOn : NSControlStateValueOff;
            [self item:BattItemMagSafeDisabled].state = identifier == BattItemMagSafeDisabled
                ? NSControlStateValueOn : NSControlStateValueOff;
            [self item:BattItemMagSafeAlwaysOff].state = identifier == BattItemMagSafeAlwaysOff
                ? NSControlStateValueOn : NSControlStateValueOff;
            break;
        case BattItemPreventIdleSleep:
        case BattItemDisableChargingPreSleep:
        case BattItemPreventSystemSleep:
        case BattItemPreventSleepOnAdapterDisable:
        case BattItemNoDrain:
            sender.state = sender.state == NSControlStateValueOff
                ? NSControlStateValueOn : NSControlStateValueOff;
            break;
        default:
            break;
    }
    battMenuAction(_handle, identifier, sender.state == NSControlStateValueOn);
}

- (void)noop:(NSMenuItem *)sender {
    (void)sender;
}

- (void)menuWillOpen:(NSMenu *)menu {
    (void)menu;
    battMenuWillOpen(_handle);
    self.timer = [NSTimer timerWithTimeInterval:BattMenuUpdateInterval
                                        target:self
                                      selector:@selector(timerTick:)
                                      userInfo:nil
                                       repeats:YES];
    [[NSRunLoop mainRunLoop] addTimer:self.timer forMode:NSRunLoopCommonModes];
}

- (void)menuDidClose:(NSMenu *)menu {
    (void)menu;
    [self.timer invalidate];
    self.timer = nil;
}

- (void)timerTick:(NSTimer *)timer {
    (void)timer;
    battMenuTimerFired(_handle);
}

- (void)setStatusIconInstalled:(BOOL)installed
                       capable:(BOOL)capable
                  needsUpgrade:(BOOL)needsUpgrade {
    NSString *symbol = @"minus.plus.batteryblock";
    NSString *description = BattT(@"batt icon");
    if (!installed) {
        symbol = @"batteryblock.slash";
        description = BattT(@"batt daemon not installed");
    } else if (!capable) {
        symbol = @"minus.plus.batteryblock.exclamationmark";
        description = BattT(@"Your machine cannot run batt");
    } else if (needsUpgrade) {
        symbol = @"fluid.batteryblock";
        description = BattT(@"batt needs upgrade");
    }
    self.statusItem.button.image = [NSImage imageWithSystemSymbolName:symbol
                                            accessibilityDescription:description];
}

- (void)setPowerItem:(BattMenuItem)item label:(NSString *)rawLabel value:(double)value {
    BOOL isSystem = [rawLabel isEqualToString:@"System"];
    NSString *label = BattT(rawLabel);

    NSColor *color = NSColor.labelColor;
    char sign = ' ';
    if (!isSystem) {
        if (value > 0) {
            color = NSColor.systemGreenColor;
            sign = '+';
        } else if (value < 0) {
            color = NSColor.systemRedColor;
            sign = '-';
        }
    }

    NSString *labelWithColon = [label stringByAppendingString:@":"];
    NSMutableString *padded = [[labelWithColon mutableCopy] autorelease];
    // 中文按两个字符宽计算，保持在等宽字体下与英文版对齐。
    NSInteger displayWidth = 0;
    for (NSUInteger i = 0; i < padded.length; i++) {
        unichar ch = [padded characterAtIndex:i];
        displayWidth += (ch > 0x2E80) ? 2 : 1;
    }
    while (displayWidth < 9) {
        [padded appendString:@" "];
        displayWidth++;
    }

    NSString *text = [NSString stringWithFormat:@"%@%c%7.2fW", padded, sign, fabs(value)];
    NSMutableAttributedString *attributed = [[[NSMutableAttributedString alloc]
        initWithString:text] autorelease];
    NSRange labelRange = NSMakeRange(0, padded.length);
    NSRange valueRange = NSMakeRange(padded.length, text.length - padded.length);
    [attributed addAttribute:NSForegroundColorAttributeName
                       value:NSColor.secondaryLabelColor
                       range:labelRange];
    [attributed addAttribute:NSForegroundColorAttributeName value:color range:valueRange];
    [attributed addAttribute:NSFontAttributeName
                       value:[NSFont monospacedSystemFontOfSize:12 weight:NSFontWeightRegular]
                       range:NSMakeRange(0, text.length)];
    [self item:item].attributedTitle = attributed;
}

@end

void BattBuildMenu(BattMenuController *controller, NSString *version) {
    NSMenu *root = controller.menu;
    NSMenu *power = AddSubmenu(controller, root, BattT(@"Power Flow"), BattItemPowerFlow);
    [power addItem:DisplayItem(controller, @"", BattItemPowerSystem, YES)];
    [power addItem:DisplayItem(controller, @"", BattItemPowerAdapter, YES)];
    [power addItem:DisplayItem(controller, @"", BattItemPowerBattery, YES)];
    [controller setPowerItem:BattItemPowerSystem label:@"System" value:0];
    [controller setPowerItem:BattItemPowerAdapter label:@"Adapter" value:0];
    [controller setPowerItem:BattItemPowerBattery label:@"Battery" value:0];

    [root addItem:ActionItem(controller, BattT(@"Upgrade Daemon..."), @"u", BattItemUpgrade)];
    [root addItem:ActionItem(controller, BattT(@"Install Daemon..."), @"i", BattItemInstall)];
    [root addItem:DisplayItem(controller, BattT(@"Loading..."), BattItemState, NO)];
    [root addItem:DisplayItem(controller, BattT(@"Loading..."), BattItemCurrentLimit, NO)];
    NSMenuItem *nativeLimitNotice = DisplayItem(controller, @"", BattItemNativeLimitNotice, NO);
    nativeLimitNotice.hidden = YES;
    [root addItem:nativeLimitNotice];
    [root addItem:[NSMenuItem separatorItem]];
    [root addItem:DisplayItem(controller, BattT(@"Quick Limits"), BattItemQuickLimits, NO)];

    const NSInteger limits[] = {50, 60, 70, 80, 90};
    const BattMenuItem limitItems[] = {
        BattItemLimit50, BattItemLimit60, BattItemLimit70, BattItemLimit80, BattItemLimit90,
    };
    for (NSUInteger index = 0; index < 5; index++) {
        NSInteger limit = limits[index];
        NSString *title = [NSString stringWithFormat:BattT(@"Set %ld%% Limit"), (long)limit];
        NSString *key = [NSString stringWithFormat:@"%ld", (long)limit];
        [root addItem:ActionItem(controller, title, key, limitItems[index])];
    }

    [root addItem:[NSMenuItem separatorItem]];
    [root addItem:DisplayItem(controller, BattT(@"One-Time Charge"), BattItemChargeOnce, NO)];
    [root addItem:ActionItem(controller, BattT(@"Charge to Limit Now"), @"",
                              BattItemChargeOnceLimit)];
    [root addItem:ActionItem(controller, BattT(@"Charge to Full Once"), @"",
                              BattItemChargeOnceFull)];
    NSMenuItem *chargeOnceStatus = DisplayItem(controller, @"", BattItemChargeOnceStatus, NO);
    chargeOnceStatus.hidden = YES;
    [root addItem:chargeOnceStatus];
    NSMenuItem *chargeOnceCancel = ActionItem(controller, BattT(@"Cancel One-Time Charge"), @"",
                                              BattItemChargeOnceCancel);
    chargeOnceCancel.hidden = YES;
    [root addItem:chargeOnceCancel];

    [root addItem:[NSMenuItem separatorItem]];
    NSMenu *advanced = AddSubmenu(controller, root, BattT(@"Advanced"), BattItemAdvanced);
    NSMenu *magSafe = AddSubmenu(controller, advanced, BattT(@"Control MagSafe LED"), BattItemMagSafe);
    [magSafe addItem:ActionItem(controller, BattT(@"Enable"), @"", BattItemMagSafeEnabled)];
    [magSafe addItem:ActionItem(controller, BattT(@"Disable"), @"", BattItemMagSafeDisabled)];
    [magSafe addItem:ActionItem(controller, BattT(@"Always-off"), @"", BattItemMagSafeAlwaysOff)];

    [advanced addItem:ActionItem(controller, BattT(@"Prevent Idle Sleep when Charging"), @"",
                                  BattItemPreventIdleSleep)];
    [advanced addItem:ActionItem(controller, BattT(@"Disable Charging before Sleep"), @"",
                                  BattItemDisableChargingPreSleep)];
    [advanced addItem:ActionItem(controller,
                                  BattT(@"Prevent System Sleep when Charging (Experimental)"), @"",
                                  BattItemPreventSystemSleep)];
    [advanced addItem:ActionItem(controller,
                                  BattT(@"Prevent Sleep when Adapter is Disabled"), @"",
                                  BattItemPreventSleepOnAdapterDisable)];
    [advanced addItem:ActionItem(controller,
                                  BattT(@"No-drain Above Native Limit (Experimental)"), @"",
                                  BattItemNoDrain)];
    NSMenu *forceDischarge = AddSubmenu(controller, advanced, BattT(@"Force Discharge..."),
                                        BattItemForceDischarge);
    NSMenuItem *forceDischargeCountdown = DisplayItem(controller, @"",
                                                       BattItemForceDischargeCountdown, NO);
    forceDischargeCountdown.hidden = YES;
    [forceDischarge addItem:forceDischargeCountdown];
    [forceDischarge addItem:ActionItem(controller, BattT(@"Stop Force Discharge"), @"",
                                        BattItemForceDischargeStop)];
    [forceDischarge addItem:ActionItem(controller, BattT(@"Indefinitely"), @"",
                                        BattItemForceDischargeIndefinitely)];
    [forceDischarge addItem:[NSMenuItem separatorItem]];
    [forceDischarge addItem:ActionItem(controller, BattT(@"1 Hour"), @"",
                                        BattItemForceDischarge1Hour)];
    [forceDischarge addItem:ActionItem(controller, BattT(@"2 Hours"), @"",
                                        BattItemForceDischarge2Hours)];
    [forceDischarge addItem:ActionItem(controller, BattT(@"4 Hours"), @"",
                                        BattItemForceDischarge4Hours)];
    [forceDischarge addItem:ActionItem(controller, BattT(@"8 Hours"), @"",
                                        BattItemForceDischarge8Hours)];

    NSMenu *calibration = AddSubmenu(controller, advanced,
                                     BattT(@"Auto Calibration (Experimental)..."),
                                     BattItemAutoCalibration);
    [calibration addItem:DisplayItem(controller, BattT(@"Status: Idle"),
                                      BattItemCalibrationStatus, NO)];
    [calibration addItem:ActionItem(controller, BattT(@"Start"), @"", BattItemCalibrationStart)];
    [calibration addItem:ActionItem(controller, BattT(@"Pause"), @"", BattItemCalibrationPause)];
    [calibration addItem:ActionItem(controller, BattT(@"Resume"), @"", BattItemCalibrationResume)];
    [calibration addItem:ActionItem(controller, BattT(@"Cancel"), @"", BattItemCalibrationCancel)];

    [advanced addItem:[NSMenuItem separatorItem]];
    NSMenu *language = AddSubmenu(controller, advanced, BattT(@"Language"), BattItemLanguage);
    [language addItem:ActionItem(controller, BattT(@"Auto (System Language)"), @"",
                                 BattItemLanguageAuto)];
    [language addItem:ActionItem(controller, BattT(@"Chinese"), @"", BattItemLanguageChinese)];
    [language addItem:ActionItem(controller, BattT(@"English"), @"", BattItemLanguageEnglish)];

    [advanced addItem:[NSMenuItem separatorItem]];
    NSString *versionTitle = [BattT(@"Version: ") stringByAppendingString:version ?: @""];
    [advanced addItem:DisplayItem(controller, versionTitle, BattItemVersion, NO)];
    [advanced addItem:ActionItem(controller, BattT(@"Uninstall Daemon..."), @"", BattItemUninstall)];

    [root addItem:[NSMenuItem separatorItem]];
    NSMenu *disableLimit = AddSubmenu(controller, root, BattT(@"Disable Charge Limit"),
                                      BattItemDisableLimit);
    NSMenuItem *disableCountdown = DisplayItem(controller, @"",
                                                BattItemDisableLimitCountdown, NO);
    disableCountdown.hidden = YES;
    [disableLimit addItem:disableCountdown];
    [disableLimit addItem:ActionItem(controller, BattT(@"Indefinitely"), @"d",
                                      BattItemDisableLimitIndefinitely)];
    [disableLimit addItem:[NSMenuItem separatorItem]];
    [disableLimit addItem:ActionItem(controller, BattT(@"1 Hour"), @"",
                                      BattItemDisableLimit1Hour)];
    [disableLimit addItem:ActionItem(controller, BattT(@"2 Hours"), @"",
                                      BattItemDisableLimit2Hours)];
    [disableLimit addItem:ActionItem(controller, BattT(@"4 Hours"), @"",
                                      BattItemDisableLimit4Hours)];
    [disableLimit addItem:ActionItem(controller, BattT(@"8 Hours"), @"",
                                      BattItemDisableLimit8Hours)];
    [disableLimit addItem:ActionItem(controller, BattT(@"12 Hours"), @"",
                                      BattItemDisableLimit12Hours)];
    [disableLimit addItem:ActionItem(controller, BattT(@"24 Hours"), @"",
                                      BattItemDisableLimit24Hours)];
    [disableLimit addItem:ActionItem(controller, BattT(@"2 Days"), @"",
                                      BattItemDisableLimit2Days)];
    [disableLimit addItem:ActionItem(controller, BattT(@"3 Days"), @"",
                                      BattItemDisableLimit3Days)];
    [disableLimit addItem:ActionItem(controller, BattT(@"7 Days"), @"",
                                      BattItemDisableLimit7Days)];
    [root addItem:ActionItem(controller, BattT(@"Quit Menubar App"), @"q", BattItemQuit)];
}
