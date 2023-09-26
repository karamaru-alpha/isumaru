## Isumaru

### Requirements

- [go](https://github.com/golang/go)
- [alp](https://github.com/tkuchiki/alp)
- [slp](https://github.com/tkuchiki/slp)


### Tutorial in local

```sh
$ make run-web
$ make run-agent
$ make run-api
```

#### Show AccessLog

cf. http://localhost:5173/accesslog/1695447044/isu1

<img width="500" src="https://github.com/karamaru-alpha/isumaru/assets/38310693/7aeaa88d-f035-4e0f-b211-234ca94a48cb">


#### Show SlowQueryLog

cf. http://localhost:5173/slowquerylog/1695447044/isu1

<img width="500" src="https://github.com/karamaru-alpha/isumaru/assets/38310693/00b99df3-267f-4272-90b1-a4c85f3144f3">

#### Collect Access/SlowQuery Log

- Push collect button.
- Insert mock access/slowQuery log within a duration(default: 15s).
  - `$ make access`

cf. http://localhost:5173/

<img width="500" src="https://github.com/karamaru-alpha/isumaru/assets/38310693/bae57213-b807-4a21-8405-19bd6b371fef">


#### Edit Config

- You can edit target/alp/slp config.

cf. http://localhost:5173/setting

<img width="500" src="https://github.com/karamaru-alpha/isumaru/assets/38310693/5fd89809-07b1-49c9-a2e1-0d44e698bd09">
