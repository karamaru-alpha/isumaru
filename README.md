## Isumaru

### Requirements

- [go](https://github.com/golang/go)
- [alp](https://github.com/tkuchiki/alp)
- [slp](https://github.com/tkuchiki/slp)


### Tutorial

```sh
$ make run-web
$ make run-agent
$ make run
```

#### Show AccessLog

cf. http://localhost:5173/accesslog/1695447044/isu1


#### Show SlowQueryLog

cf. http://localhost:5173/slowquerylog/1695447044/isu1

#### Collect Access/SlowQuery Log

- Push collect button.
- Insert mock access/slowQuery log `$ make access`

cf. http://localhost:5173/
