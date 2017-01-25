#import <Foundation/Foundation.h>
#import <AppKit/NSRunningApplication.h>
#import <Foundation/NSBundle.h>
#import <Foundation/NSString.h>
#import <CoreFoundation/CoreFoundation.h>
#include <stdio.h>
#include "binding.h"
#include "version.h"

static CFBundleRef mesonFramework = NULL;

static const unsigned char mesonAPIVersions[] __attribute__ ((section ("__MESON_VERSION,__meson_version"))) = {
  MESON_VERSION_MAJOR,
  MESON_VERSION_MINOR,
  MESON_VERSION_REVISON,
};

const unsigned char* mesonVersions() {
  return mesonAPIVersions;
}

int loadMesonFramework(const char* path)
{
  NSString* strPath;
  BOOL success = FALSE;
  strPath = [NSString stringWithUTF8String:path];
  if (strPath) {
    NSURL* fileURL = [NSURL fileURLWithPath:strPath];
    if (fileURL) {
      CFURLRef url = (CFURLRef )fileURL;
      if (url) {
        CFBundleRef bundle = CFBundleCreate(kCFAllocatorDefault, url);
        if (bundle) {
          MesonExportFunction* pExports = MesonFrameworkFunctions;
          success = TRUE;
          while (pExports->name) {
            CFStringRef name = CFStringCreateWithCString(kCFAllocatorDefault, pExports->name, kCFStringEncodingUTF8);
            if (!name) {
              success = FALSE;
              break;
            }
            pExports->func =  CFBundleGetFunctionPointerForName(bundle, name);
            if (!pExports->func) {
              success = FALSE;
              break;
            }
            CFRelease(name);
            ++pExports;
          }

          if (success) {
            // keep framework reference
            mesonFramework = (CFBundleRef )CFRetain(bundle);
          }

          CFRelease(bundle);
        }
      }
      [fileURL release];
    }
    [strPath release];
  }
  return success;
}

void freeMesonFramework()
{
  if (mesonFramework) {
    CFRelease(mesonFramework);
    mesonFramework = NULL;
  }
}
