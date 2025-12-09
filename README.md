# libatrus-go
Go bindings for [libatrus](https://github.com/sinclairtarget/libatrus).

## Installing
This module depends on libatrus. You can depend on a system version of libatrus
(perhaps installed by a package manager) or on the vendored version of libatrus
included in this repository.

### Using Your System Version
If you already have libatrus installed, then Go will use `pkg-config` to locate
the necessary files at build time. You should be able to use this module like
any other Go module.

The system version of libatrus will be dynamically linked into your final
executable unless you specify otherwise with additional build flags. 

### Using the Vendored Version
As an alternative to using your system version of libatrus, you can build the
vendored copy of the library and statically link with it. You must have Zig
installed for this to work. You must also check out this repository locally.

Once you have checked out the repository, from the repository root, run:
```
$ git submodule update --init
```

Then navigate to `vendored/libatrus` and run `zig build`.

When building your downstream Go application, you must pass `-tags
bundled_libatrus` so that the locally built version of libatrus is used.

You must also add a `replace` directive to your `go.mod` file so that Go uses
the copy of the repository that you just checked out to satisfy your dependency
on this module. That directive should look something like:

```
replace github.com/sinclairtarget/libatrus-go => ../libatrus-go
```

Though the second part after the `=>` depends on where you've checked out this
repository.

## Running Tests
Build the vendored Zig library. Then run:
```
$ go test -tags bundled_libatrus
```
