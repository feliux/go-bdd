# Behaviour Driven Development

```bash
$ go install github.com/cucumber/godog/cmd/godog@latest
$ export PATH="$PATH:$(go env GOPATH)/bin"

$ cd <folder>
$ godog

# Ejecuta los test para el tag puntuación
$ godog -t=@puntuacion
# Ejecuta escenarios de manera aleatoria
$ godog --random # útil para detectar dependencias
# Ejecuta escenarios de manera concurrente
$ godog --concurrency=2
$ godog --concurrency 2 --format progress
```

### Good practices

Do this...

1. Describe the behaviour, not the implementation. Be concise.
2. Use the context to share the state.

### Antipatterns

Do not do this...

1. Alternate order Given-When-Then

## References

[godog](https://github.com/cucumber/godog)

[cucumber-html-reporter](https://github.com/gkushang/cucumber-html-reporter)
