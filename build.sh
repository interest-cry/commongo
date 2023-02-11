test_httpnet_build_run(){
  go build example/cmd/testhttp.go
  cp testhttp cli
  ./testhttp -r=ser
}
test_httpnetgo_build_run(){
  go build example/cmd/testhttpnetgo.go
  cp testhttpnetgo cligo
  time ./testhttpnetgo -r=ser
}
test_httpnetgo_build_run