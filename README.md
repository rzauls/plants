# Plants API

Generic-ish http API, originally to test `swaggo/swag` openAPI spec generation, but eventually refactored into a testpiece 
for checking out the new `http.ServeMux` router and `slog` package from std lib.

## Test watcher
There is a `Justfile` that contains a `nodemon` command for watching tests while writing more code. These are not required to run the API in any way, just for my personal DX.

To use it you need 2 applications:

    - `just` a command runner with saner defaults than `make` <https://github.com/casey/just>

    - `nodemon` a file watcher that runs a command after files get modified <https://nodemon.io/>
