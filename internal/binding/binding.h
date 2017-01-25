#ifndef __INC_MESON_BINDING_BINDING_H__
#define __INC_MESON_BINDING_BINDING_H__

#ifdef __cplusplus
extern "C" {
#endif
#define MESON_EXPORT __attribute__((visibility("default")))

  typedef struct MesonExportFunction {
    const char* name;
    void* func;
  } MesonExportFunction;

  extern MesonExportFunction MesonFrameworkFunctions[];

  extern int loadMesonFramework(const char* path);
  extern void freeMesonFramework();
  extern void mesonRegistHandler();
  extern const unsigned char* mesonVersions();


#ifdef __cplusplus
};
#endif

#endif
