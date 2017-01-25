#include "binding.h"
#include "assert.h"
#include "stdio.h"
#include "_cgo_export.h"
#include "meson.h"

enum {
  API_MesonApiSetArgc,
  API_MesonApiAddArgv,
  API_MesonApiMain,
  API_MesonApiSetHandler,
};

MesonExportFunction MesonFrameworkFunctions[] = {
    {"MesonApiSetArgc", NULL},
    {"MesonApiAddArgv", NULL},
    {"MesonApiMain", NULL},
    {"MesonApiSetHandler", NULL},
    {NULL, NULL},
};

void MesonApiSetArgc(int argc) {
  typedef void (*tfn)(int argc);
  tfn pfn = (tfn)MesonFrameworkFunctions[API_MesonApiSetArgc].func;
  assert(pfn);
  pfn(argc);
}

void MesonApiAddArgv(int i, const char *argv) {
  typedef void (*tfn)(int i, const char *argv);
  tfn pfn = (tfn)MesonFrameworkFunctions[API_MesonApiAddArgv].func;
  assert(pfn);
  pfn(i, argv);
}

int MesonApiMain(void) {
  typedef int (*tfn)(void);
  tfn pfn = (tfn)MesonFrameworkFunctions[API_MesonApiMain].func;
  assert(pfn);
  return pfn();
}

static void mesonCallInitHandler(void)
{
	goCallInit();
}
static char* mesonWaitServerRequestHandler(void)
{
	return goWaitServerRequest();
}
static char* mesonPostServerResponseHandler(unsigned int id, const char* response, int needReply)
{
  return goPostServerResponse(id, (char* )response, needReply);

}
void mesonRegistHandler()
{
  typedef void (*tfn)(MesonInitHandler, MesonWaitServerRequestHandler,
                      MesonPostServerResponseHandler);
  tfn pfn = (tfn)MesonFrameworkFunctions[API_MesonApiSetHandler].func;
  assert(pfn);
  pfn(mesonCallInitHandler, mesonWaitServerRequestHandler, mesonPostServerResponseHandler);
}
