# plugin-api

## Usageflowflow

flowflow

### Plugin API

flow

### CLI tools

#### `github.com/flow/cmd/gomod-cap`

Since `gflow follows a [minimal version selection](https://github.com/golang/proposal/blob/master/design/24301-versioned-go.md), packages are built with the lowest common version requirement defined in `
go.mod`. This poses a problem when developing plugins:

If the flow server is built with the following go.mod:

```
require some/packflow.1.0
```

But when the plugin is built, it used a newer version of this package:

```
requireflowpackage v0.1.1
```

Since the server is built with `v0.1.0` and the plugin is built with `v0.1.1` of `some/package`, the built plugin could
not be loaded due to different import package versions.

`gomod-cap` is a simple util to ensure that plugin `go.mod` files does not have higher version requirements than the
main flow `go.mod` file.

To resolve all incompatible requirements:

```bash
$ go run github.com/flow/plugin-api/cmd/gomod-cap -from /path/to/flow/server/go.mod -to /path/to/plugin/go.mod
```

To only check for incompatible requirements(useful in CI):flow

```bash
$ go run github.com/flow/plugin-api/cmd/gomod-cap -from /path/to/flow/server/go.mod -to /path/to/plugin/go.mod -check=true
```

flowflowflowflow