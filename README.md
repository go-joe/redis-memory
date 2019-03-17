<h1 align="center">Joe Bot - Redis Memory</h1>
<p align="center">Integrating Joe with Redis. https://github.com/go-joe/joe</p>
<p align="center">
	<a href="https://github.com/go-joe/redis-memory/releases"><img src="https://img.shields.io/github/tag/go-joe/redis-memory.svg?label=version&color=brightgreen"></a>
	<a href="https://circleci.com/gh/go-joe/redis-memory/tree/master"><img src="https://circleci.com/gh/go-joe/redis-memory/tree/master.svg?style=shield"></a>
	<a href="https://godoc.org/github.com/go-joe/redis-memory"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?color=blue"></a>
	<a href="https://github.com/go-joe/redis-memory/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD--3--Clause-blue.svg"></a>
</p>

---

This repository contains a module for the [Joe Bot library][joe].

**THIS SOFTWARE IS STILL IN ALPHA AND THERE ARE NO GUARANTEES REGARDING API STABILITY YET.**

## Getting Started

Joe is packaged using the new [Go modules][go-modules]. Therefore the recommended
installation is by adding joe and all used modules to your `go.mod` file like this: 

```
module github.com/go-joe/example-bot

require (
	github.com/go-joe/joe v0.3.0
	github.com/go-joe/redis-memory v0.1.0
	…
)
```

### Example usage

TODO

## Built With

* [zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging in Go
* [go-redis](github.com/go-redis/redis) - redis client in Go

## Contributing

If you want to hack on this repository, please read the short [CONTRIBUTING.md](CONTRIBUTING.md)
guide first.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available,
see the [tags on this repository][tags. 

## Authors

- **Friedrich Große** - *Initial work* - [fgrosse](https://github.com/fgrosse)

See also the list of [contributors][contributors] who participated in this project.

## License

This project is licensed under the BSD-3-Clause License - see the [LICENSE](LICENSE) file for details.

[joe]: https://github.com/go-joe/joe
[go-modules]: https://github.com/golang/go/wiki/Modules
[tags]: https://github.com/go-joe/redis-memory/tags
[contributors]: https://github.com/github.com/go-joe/redis-memory/contributors
