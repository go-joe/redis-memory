<h1 align="center">Joe Bot - Redis Memory</h1>
<p align="center">Integrating Joe with Redis. https://github.com/go-joe/joe</p>
<p align="center">
	<a href="https://github.com/go-joe/redis-memory/releases"><img src="https://img.shields.io/github/tag/go-joe/redis-memory.svg?label=version&color=brightgreen"></a>
	<a href="https://circleci.com/gh/go-joe/redis-memory/tree/master"><img src="https://circleci.com/gh/go-joe/redis-memory/tree/master.svg?style=shield"></a>
	<a href="https://goreportcard.com/report/github.com/go-joe/redis-memory"><img src="https://goreportcard.com/badge/github.com/go-joe/redis-memory"></a>
	<a href="https://codecov.io/gh/go-joe/redis-memory"><img src="https://codecov.io/gh/go-joe/redis-memory/branch/master/graph/badge.svg"/></a>
	<a href="https://godoc.org/github.com/go-joe/redis-memory"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?color=blue"></a>
	<a href="https://github.com/go-joe/redis-memory/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-BSD--3--Clause-blue.svg"></a>
</p>

---

This repository contains a module for the [Joe Bot library][joe].

**THIS SOFTWARE IS STILL IN ALPHA AND THERE ARE NO GUARANTEES REGARDING API STABILITY YET.**

## Getting Started

This library is packaged using the new [Go modules][go-modules]. You can get it via:

```
go get github.com/go-joe/redis-memory
```

### Example usage

In order to connect your bot to Redis you can simply pass it as module when
creating a new bot:

[embedmd]:# (_examples/main.go)
```go
package main

import (
	"github.com/go-joe/joe"
	"github.com/go-joe/redis-memory"
	"github.com/pkg/errors"
)

type ExampleBot struct {
	*joe.Bot
}

func main() {
	b := &ExampleBot{
		Bot: joe.New("example",
			redis.Memory("localhost:6379"),
		),
	}

	b.Respond("remember (.+) is (.+)", b.Remember)
	b.Respond("what is (.+)", b.WhatIs)
	b.Respond("show keys", b.ShowKeys)

	err := b.Run()
	if err != nil {
		b.Logger.Fatal(err.Error())
	}
}

func (b *ExampleBot) Remember(msg joe.Message) error {
	key, value := msg.Matches[0], msg.Matches[1]
	msg.Respond("OK, I'll remember %s is %s", key, value)
	return b.Store.Set(key, value)
}

func (b *ExampleBot) WhatIs(msg joe.Message) error {
	key := msg.Matches[0]
	var value string
	ok, err := b.Store.Get(key, &value)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve key %q from brain", key)
	}

	if ok {
		msg.Respond("%s is %s", key, value)
	} else {
		msg.Respond("I do not remember %q", key)
	}

	return nil
}

func (b *ExampleBot) ShowKeys(msg joe.Message) error {
	keys, err := b.Store.Keys()
	if err != nil {
		return err
	}

	msg.Respond("I got %d keys:", len(keys))
	for i, k := range keys {
		msg.Respond("%d) %q", i+1, k)
	}
	return nil
}
```

## Built With

* [go-redis](https://github.com/go-redis/redis) - redis client in Go
* [redimock](https://github.com/fzerorubigd/redimock) - redis mock library in tcp level
* [pkg/errors](https://github.com/pkg/errors) - Simple error handling primitives
* [testify](https://github.com/stretchr/testify) - A simple unit test library
* [zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging in Go

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

## Acknowledgments

- [embedmd][embedmd] for a cool tool to embed source code in markdown files

[joe]: https://github.com/go-joe/joe
[go-modules]: https://github.com/golang/go/wiki/Modules
[tags]: https://github.com/go-joe/redis-memory/tags
[contributors]: https://github.com/github.com/go-joe/redis-memory/contributors
