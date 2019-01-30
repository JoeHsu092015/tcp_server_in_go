
MAINPATH=./cmd/main.go
EXECNAME=server
EXECDIR=exec

go-run:
	GO111MODULE=on go run $(MAINPATH)

go-build:
	GO111MODULE=on go build -o $(EXECNAME) $(MAINPATH)

go-build_linux_x86:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build  -o $(EXECDIR)/linux/x86/$(EXECNAME) $(MAINPATH)
go-build_linux_x64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o $(EXECDIR)/linux/x64/$(EXECNAME) $(MAINPATH)
go-build_windows_x86:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build  -o $(EXECDIR)/windows/x86/$(EXECNAME) $(MAINPATH)
go-build_windows_x64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o $(EXECDIR)/windows/x64/$(EXECNAME) $(MAINPATH)

go-build-all: go-build_linux_x86 go-build_linux_x64 go-build_windows_x86 go-build_windows_x64

