package canned

import (
	"github.com/litesoft-go/mockvoltdb/pkg/utils"
)

func PortResonses_http() *utils.Responder {
	r := &utils.Responder{}

	r.AddPath(`/api/2.0?Procedure=@PingPartitions&Parameters=[0]`).
		Add(`
{
	"status":1
}
`).
		NoMore()

	r.AddPath(`/api/2.0?admin=true&Procedure=@PrepareShutdown`).
		Add(`
{
	"status":1,	"results":{"0":[{"STATUS":3268496242920914946}]}
}
`).
		NoMore()
	r.AddPath(`/api/2.0?admin=true&Procedure=@Quiesce`).
		Add(`
{
	"status": 1, "results": {"0": [{"STATUS": 0}]}
}
`).
		NoMore()

	r.AddPath(`/api/2.0?admin=true&Procedure=@Statistics&Parameters=["SHUTDOWN_CHECK",0]`).
		Add(`
{
	"status": 1,
	"results": {"0": [
			{"ACTIVE": 0},
			{"ACTIVE": 3},
			{"ACTIVE": 0}]}
}
`).
		Add(`
{
	"status": 1,
	"results": {"0": [
			{"ACTIVE": 0},
			{"ACTIVE": 0}]}
}
`).
		NoMore()

	r.AddPath(`/api/2.0?admin=true&Procedure=@OpPseudoShutdown`).
		Add(`
{
	"status": 1,
	"results": {"0": [{"STATUS": 0}]}
}
`).
		NoMore()
	r.AddPath(`/api/2.0?admin=true&Procedure=@PrepareStopNode&Parameters=[0]`).
		Add(`
{
	"status":1,
	"results":{}
}
`).
		NoMore()

	r.AddPath(`/api/2.0?admin=true&Procedure=@Statistics&Parameters=["STOP_CHECK",0]`).
		Add(`
{
	"status": 1,
	"results": {
    "0": [
      {
        "HOST_ID": 0,
        "ACTIVE": 1
      },
      {
        "HOST_ID": 1,
        "ACTIVE": 1
      }]}
}
`).
		Add(`
{
	"status": 1,
	"results": {
    "0": [
      {
        "HOST_ID": 0,
        "ACTIVE": 0
      },
      {
        "HOST_ID": 1,
        "ACTIVE": 1
      }]}
}
`).
		NoMore()

	r.AddPath(`/api/2.0?admin=true&Procedure=@StopNode&Parameters=[0]`).
		Add(`
{
	"status":-5
}
`).
		NoMore()

	return r
}
