# webr
Recompile and rerun Go webapps when needed.

It acts as a proxy and blocks incoming requests while the program compiles in the background to support a recompile on reload kind of workflow.

Fetch the binary with go get like so: `go get github.com/ksev/webr`

The only requirement is that your target app has a configurable bind string for the http server that can be passed in through the command line.
The webr binary only has three command line arguments -bind, -target-bind-arg-name and -h.

The default values explicitly passed would look something like this: 

`webr -bind :8080 -target-bind-arg-name "-bind" -h false`

You can also passtrough command line arguments to the target binary by using `--`:

`webr -bind :8080 -- -arg value` 
