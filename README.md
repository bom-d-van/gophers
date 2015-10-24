# gophers

Gophers runs scripts or commands concurrently and kills them when gophers was killed.

# Usage

```sh
gophers gulp "find . -iname '*.go' | entr -rd grm"
```

What is grm: https://github.com/bom-d-van/bin/blob/master/grm