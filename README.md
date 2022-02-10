# gotestgen
## description
generate Go test templates

## usage
```go install github.com/kimuson13/gotestgen/cmd/gotestgen@latest```  

after changing directory you want to generate test files and  

```gotestgen ./...```

## options
### -p
add `t.Parallel()` into some test functions
### -g
select the place that test file is wrote
#### example
if you type `tree .`
```
$ tree .
.
├── hoge
│   └── hoge.go
├ go.mod
├ go.sum
└ main.go
```
```
$ cat hoge/hoge.go

package hoge

func Hoge() {
  fmt.Println("Hoge")
}

func hoge() {
  fmt.Println("hoge")
}
```
```
$ gotestgen -g=[hoge:hoge]
```
```
$ tree .
.
├── hoge
│   ├── hoge.go
│   └── hoge_test.go
├ go.mod
├ go.sum
└ main.go
```
```
$ cat hoge/hoge_test.go

package hoge_test

import "testing"

func TestHoge(t *testing.T) {
  cases := map[string]struct {
  
  }{
    //write test cases below
  }
  
  for testName, tt := range cases {
		tt := tt
		t.Run(testName, func(t *testing.T) {
			
		})
	}
}
```
