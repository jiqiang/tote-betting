# A tote betting calculator

git clone https://github.com/jiqiang/tote-betting.git

## run calculator

```sh
$ go run ./main.go
```

## feed from file

```sh
$ cat bets | go run ./main.go
```

## generate executable

```sh
$ go build -o tbc
```

## run executable

```sh
$ ./tbc
```
or
```sh
$ cat bets | ./tbc
```
