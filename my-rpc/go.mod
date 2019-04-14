module my-rpc

go 1.12

replace golang.org/x/sys@v0.0.0-20180905080454-ebe1bf3edb33 => github.com/golang/sys v0.0.0-20190405154228-4b34438f7a67

require (
	github.com/kr/pretty v0.1.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
