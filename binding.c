#include "meson.h"
#include "_cgo_export.h"
#include <stdlib.h>

void mesonCallInitHandler(void)
{
	goCallInit();
}
char* mesonWaitServerRequestHandler(void)
{
	return goWaitServerRequest();
}
char* mesonPostServerResponseHandler(unsigned int id, const char* response, int needReply)
{
  return goPostServerResponse(id, (char* )response, needReply);

}
void mesonRegistHandler()
{
	MesonApiSetHandler(mesonCallInitHandler, mesonWaitServerRequestHandler, mesonPostServerResponseHandler);
}
