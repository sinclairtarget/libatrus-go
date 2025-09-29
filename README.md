# libatrus-go
Go bindings for [libatrus](https://github.com/sinclairtarget/libatrus).

## Running Tests
First, run:

```
$ git submodule update --init
```

Then build the vendored Zig library. Finally, you can run:

```
$ PKG_CONFIG_PATH=vendored/libatrus/zig-out/share/pkgconfig/ go test
```
