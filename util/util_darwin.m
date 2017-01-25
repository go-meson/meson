//-*-objc-*-
#import <AppKit/NSRunningApplication.h>
#import <Foundation/Foundation.h>
#import <Foundation/NSBundle.h>
#import <Foundation/NSString.h>

char *mesonGetBundlePath() {
  NSRunningApplication *app = [NSRunningApplication currentApplication];
  char *ret = NULL;
  if (app) {
    NSURL *url = app.bundleURL;
    if (pURL) {
      NSString *str = url.path;
      if (str) {
        const char *pstr = [str UTF8String];
        ret = strdup(pstr);
        [str release];
      }
      [url release];
    }
    [app release];
  }

  return ret;
}

char *mesonGetSystemDirectoryPath(int type) {
  NSSearchPathDirectory pathType;
  NSSearchPathDomainMask domainType = NSUserDomainMask;
  BOOL isCreate = FALSE;

  switch (type) {
  case 1:
    pathType = NSCachesDirectory;
    domainType = NSUserDomainMask;
    isCreate = TRUE;
    break;
  case 2:
    pathType = NSDocumentDirectory;
    break;
  case 3:
    pathType = NSDesktopDirectory;
    break;
  default:
    return NULL;
  }

  char *ret = NULL;
  NSFileManager *fm = (NSFileManager *)[[NSFileManager alloc] init];
  if (fm) {
    NSURL *url = [fm URLForDirectory:pathType
                            inDomain:domainType
                   appropriateForURL:nil
                              create:TRUE
                               error:nil];
    if (url) {
      NSString *str = url.path;
      if (str) {
        const char *pstr = [str UTF8String];
        ret = strdup(pstr);
        [str release];
      }
      [url release];
    }
    [fm release];
  }
  return ret;
}
