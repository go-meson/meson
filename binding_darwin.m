//-*-objc-*-
#import <Foundation/Foundation.h>
#import <AppKit/NSRunningApplication.h>
#import <Foundation/NSBundle.h>
#import <Foundation/NSString.h>

char* mesonGetBundlePath() {
  NSRunningApplication* app = [NSRunningApplication currentApplication];
  char* ret = NULL;
  if (app) {
    NSURL* url = app.bundleURL;
    if (pURL) {
      NSString* str = url.path;
      if (str) {
        const char* pstr = [str UTF8String];
        ret = strdup(pstr);
        [str release];
      }
      [url release];
    }
    [app release];
  }

  return ret;
}

void mesonGetFrameworkLocation(char* buff)
{
  NSFileManager* fm = [[NSFileManager alloc] init];

  [fm URLForDirectory:NSCachesDirectory inDomain:NSUserDomainMask appropriateForURL:nil create:true error:nil];

  [fm release];
}
