//go:build darwin

package services

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices

#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>

// moveWindowToBottomRight 使用 Accessibility API 将匹配标题的窗口移至屏幕右下角
// 注意：首次运行需在「系统设置 → 隐私与安全性 → 辅助功能」中授权本应用
static void moveWindowToBottomRight(const char* title, int margin) {
	@autoreleasepool {
		NSString *targetTitle = [NSString stringWithUTF8String:title];
		NSRect visibleFrame = [[NSScreen mainScreen] visibleFrame];

		// 通过进程名查找 PID（Wails 应用进程名即为二进制名）
		pid_t targetPid = 0;
		for (NSRunningApplication *app in [[NSWorkspace sharedWorkspace] runningApplications]) {
			if ([app.localizedName containsString:@"ItdTime"]) {
				targetPid = app.processIdentifier;
				break;
			}
		}
		if (targetPid == 0) return;

		// 通过 Accessibility API 枚举窗口
		AXUIElementRef appElem = AXUIElementCreateApplication(targetPid);
		CFArrayRef windowList = NULL;
		AXError err = AXUIElementCopyAttributeValues(appElem, kAXWindowsAttribute, 0, 100, &windowList);
		if (err != kAXErrorSuccess || windowList == NULL) {
			CFRelease(appElem);
			return;
		}

		NSArray *windows = (__bridge NSArray *)windowList;
		for (id w in windows) {
			AXUIElementRef winElem = (__bridge AXUIElementRef)w;

			CFStringRef winTitle = NULL;
			AXUIElementCopyAttributeValue(winElem, kAXTitleAttribute, (CFTypeRef *)&winTitle);
			if (winTitle == NULL) continue;

			BOOL match = [[(__bridge NSString *)winTitle stringByTrimmingCharactersInSet:
				[NSCharacterSet whitespaceCharacterSet]] isEqualToString:targetTitle];
			CFRelease(winTitle);
			if (!match) continue;

			// 获取窗口尺寸
			CFTypeRef sizeRef = NULL;
			AXUIElementCopyAttributeValue(winElem, kAXSizeAttribute, &sizeRef);
			CGSize winSize = {0, 0};
			AXValueGetValue((AXValueRef)sizeRef, kAXValueTypeCGSize, &winSize);
			if (sizeRef) CFRelease(sizeRef);

			// macOS 坐标原点在左下角
			CGPoint newPos = CGPointMake(
				visibleFrame.origin.x + visibleFrame.size.width  - winSize.width  - margin,
				visibleFrame.origin.y + margin
			);
			AXValueRef posValue = AXValueCreate(kAXValueTypeCGPoint, &newPos);
			AXUIElementSetAttributeValue(winElem, kAXPositionAttribute, posValue);
			CFRelease(posValue);
			break;
		}

		CFRelease(windowList);
		CFRelease(appElem);
	}
}
*/
import "C"
import "unsafe"

// MoveWindowToBottomRight 将指定标题的窗口移动到主屏幕工作区（不含菜单栏/Dock）右下角
func MoveWindowToBottomRight(title string, margin int32) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	C.moveWindowToBottomRight(cTitle, C.int(margin))
}
