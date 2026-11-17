#import <Cocoa/Cocoa.h>
#include <stdlib.h>
#import "native_internal.h"

static NSString *BattString(const char *value) {
    return value == NULL ? @"" : [NSString stringWithUTF8String:value];
}

static BattMenuController *BattController(BattMenuRef menu) {
    return (BattMenuController *)menu;
}

BattMenuRef batt_menu_create(uintptr_t handle, const char *version) {
    @autoreleasepool {
        [NSApplication sharedApplication];
        return [[BattMenuController alloc] initWithHandle:handle version:BattString(version)];
    }
}

void batt_menu_destroy(BattMenuRef menu) {
    @autoreleasepool {
        [BattController(menu) release];
    }
}

void batt_app_run(void) {
    @autoreleasepool {
        [[NSApplication sharedApplication] run];
    }
}

void batt_app_terminate(void) {
    @autoreleasepool {
        [[NSApplication sharedApplication] terminate:nil];
    }
}

void batt_menu_set_title(BattMenuRef menu, int item, const char *title) {
    @autoreleasepool {
        [BattController(menu) item:(BattMenuItem)item].title = BattString(title);
    }
}

void batt_menu_set_tooltip(BattMenuRef menu, int item, const char *tooltip) {
    @autoreleasepool {
        [BattController(menu) item:(BattMenuItem)item].toolTip = BattString(tooltip);
    }
}

void batt_menu_set_hidden(BattMenuRef menu, int item, bool hidden) {
    [BattController(menu) item:(BattMenuItem)item].hidden = hidden;
}

void batt_menu_set_enabled(BattMenuRef menu, int item, bool enabled) {
    [BattController(menu) item:(BattMenuItem)item].enabled = enabled;
}

void batt_menu_set_checked(BattMenuRef menu, int item, bool checked) {
    [BattController(menu) item:(BattMenuItem)item].state = checked
        ? NSControlStateValueOn : NSControlStateValueOff;
}

void batt_menu_set_status_icon(BattMenuRef menu,
                               bool installed,
                               bool capable,
                               bool needs_upgrade) {
    @autoreleasepool {
        [BattController(menu) setStatusIconInstalled:installed
                                            capable:capable
                                       needsUpgrade:needs_upgrade];
    }
}

void batt_menu_set_power(BattMenuRef menu, int item, const char *label, double value) {
    @autoreleasepool {
        [BattController(menu) setPowerItem:(BattMenuItem)item label:BattString(label) value:value];
    }
}


NSString *BattT(NSString *english) {
    if (english == nil) {
        return nil;
    }
    char *translated = battTranslate((char *)[english UTF8String]);
    if (translated == NULL) {
        return english;
    }
    return [NSString stringWithUTF8String:translated];
}

const char *batt_get_language(void) {
    NSString *mode = [[NSUserDefaults standardUserDefaults] stringForKey:@"BattLanguage"];
    if (mode.length == 0) {
        mode = @"auto";
    }
    return strdup([mode UTF8String]);
}

void batt_set_language(const char *mode) {
    NSString *value = mode == NULL ? @"auto" : [NSString stringWithUTF8String:mode];
    NSUserDefaults *defaults = [NSUserDefaults standardUserDefaults];
    [defaults setObject:value forKey:@"BattLanguage"];
    [defaults synchronize];
}

const char *batt_system_language(void) {
    NSString *first = [[NSLocale preferredLanguages] firstObject] ?: @"en";
    NSString *code = [first hasPrefix:@"zh"] ? @"zh" : @"en";
    return strdup([code UTF8String]);
}

void batt_relaunch(void) {
    NSString *bundlePath = [[NSBundle mainBundle] bundlePath];
    NSString *command = [NSString stringWithFormat:@"sleep 1; /usr/bin/open -a %@", bundlePath];
    NSTask *task = [[[NSTask alloc] init] autorelease];
    task.launchPath = @"/bin/sh";
    task.arguments = @[@"-c", command];
    [task launch];
    [NSApp terminate:nil];
}

void batt_menu_rebuild(BattMenuRef menu) {
    @autoreleasepool {
        BattMenuController *controller = BattController(menu);
        [controller.items removeAllObjects];
        [controller.menu removeAllItems];
        BattBuildMenu(controller, controller.appVersion);
        BattApplyTooltips(controller);
    }
}
