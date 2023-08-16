package protobuf

//go:generate protoc --go_out=.. --go_opt=module=github.com/iyarkov/kit iyarkov/kit/hlc.proto
//go:generate protoc --go_out=.. --go_opt=module=github.com/iyarkov/kit iyarkov/kit/object.proto
//go:generate protoc --go_out=.. --go_opt=module=github.com/iyarkov/kit iyarkov/kit/page_request.proto
