# dreamkast-weaver


```
$ docker compose up -d
$ dbmate up
```

```
$ weaver multi deploy weaver.toml 
╭───────────────────────────────────────────────────╮
│ app        : dreamkast-weaver                     │
│ deployment : 4ed57675-3aeb-46f1-a6ac-4508ec345e98 │
╰───────────────────────────────────────────────────╯
S0409 17:26:14.818395 stdout  83d577f9           ] env platform is either empty or invalid
D0409 17:26:14.822094 main    83d577f9 frontend.go:82] ENV_PLATFORM platform="local"
D0409 17:26:14.823541 main    83d577f9 frontend.go:143] Frontend available addr="[::]:12345"
S0409 17:26:14.824386 stdout  80a92c60                ] env platform is either empty or invalid
D0409 17:26:14.825077 main    80a92c60 frontend.go:82 ] ENV_PLATFORM platform="local"
D0409 17:26:14.825593 main    80a92c60 frontend.go:143] Frontend available addr="[::]:12345"
```

```
$ curl localhost:12345/vote
$ curl localhost:12345/show
```
