[![Go Test](https://github.com/taylor-s-dean/errgo/actions/workflows/go.yml/badge.svg)](https://github.com/taylor-s-dean/errgo/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/taylor-s-dean/errgo.svg)](https://pkg.go.dev/github.com/taylor-s-dean/errgo)
# errgo
A verbose error wrapper in Go.

## Example
```go
code := codes.DeadlineExceeded
err := errors.New(code.String())
wrapper := Wrap(err, "wrapped error", codes.AlreadyExists)
str := fmt.Sprint(Wrap(wrapper, "another wrapped error", codes.NotFound).JSON())
fmt.Println(str)
```
```json
{
    "code": "NotFound (5)",
    "error": {
        "code": "AlreadyExists (6)",
        "error": "DeadlineExceeded",
        "file": "error_test.go",
        "function": "(*ErrorTestSuite).TestContainsCode()",
        "line": 176,
        "message": "wrapped error",
        "stack_trace": "goroutine 5 [running]:\nruntime/debug.Stack()\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/runtime/debug/stack.go:24 +0x65\ngithub.com/taylor-s-dean/errgo.Wrap({0x96bf40?, 0xc0001c4760}, {0x8cf452, 0xd}, 0x6)\n\t/home/runner/work/errgo/errgo/error.go:185 +0x66\ngithub.com/taylor-s-dean/errgo.(*ErrorTestSuite).TestContainsCode(0x0?)\n\t/home/runner/work/errgo/errgo/error_test.go:176 +0x18f\nreflect.Value.call({0xc0001c6080?, 0xc000014fc0?, 0x13?}, {0x8cb8f7, 0x4}, {0xc000061e70, 0x1, 0x1?})\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/reflect/value.go:586 +0xb07\nreflect.Value.Call({0xc0001c6080?, 0xc000014fc0?, 0xc00013ed40?}, {0xc00004f670?, 0xb58690?, 0x99ae55?})\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/reflect/value.go:370 +0xbc\ngithub.com/stretchr/testify/suite.Run.func1(0xc00016eea0)\n\t/home/runner/go/pkg/mod/github.com/stretchr/testify@v1.8.4/suite/suite.go:197 +0x4b6\ntesting.tRunner(0xc00016eea0, 0xc000196090)\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/testing/testing.go:1576 +0x10b\ncreated by testing.(*T).Run\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/testing/testing.go:1629 +0x3ea\n"
    },
    "file": "error_test.go",
    "function": "(*ErrorTestSuite).TestContainsCode()",
    "line": 177,
    "message": "another wrapped error",
    "stack_trace": "goroutine 5 [running]:\nruntime/debug.Stack()\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/runtime/debug/stack.go:24 +0x65\ngithub.com/taylor-s-dean/errgo.Wrap({0x96bfc0?, 0xc0001b2b40}, {0x8d2e6d, 0x15}, 0x5)\n\t/home/runner/work/errgo/errgo/error.go:185 +0x66\ngithub.com/taylor-s-dean/errgo.(*ErrorTestSuite).TestContainsCode(0x0?)\n\t/home/runner/work/errgo/errgo/error_test.go:177 +0x1b4\nreflect.Value.call({0xc0001c6080?, 0xc000014fc0?, 0x13?}, {0x8cb8f7, 0x4}, {0xc000061e70, 0x1, 0x1?})\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/reflect/value.go:586 +0xb07\nreflect.Value.Call({0xc0001c6080?, 0xc000014fc0?, 0xc00013ed40?}, {0xc00004f670?, 0xb58690?, 0x99ae55?})\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/reflect/value.go:370 +0xbc\ngithub.com/stretchr/testify/suite.Run.func1(0xc00016eea0)\n\t/home/runner/go/pkg/mod/github.com/stretchr/testify@v1.8.4/suite/suite.go:197 +0x4b6\ntesting.tRunner(0xc00016eea0, 0xc000196090)\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/testing/testing.go:1576 +0x10b\ncreated by testing.(*T).Run\n\t/opt/hostedtoolcache/go/1.20.10/x64/src/testing/testing.go:1629 +0x3ea\n"
}
```
