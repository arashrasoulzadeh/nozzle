### Nozzle sdk 
never miss io!

## what is it?
A daemon-based I/O provider for Go applications. It uses the inbox pattern to ensure files are written in order and only when their contents have changed.

## how to use it?
```go 
	n, e := app.StartDaemon("temp")
	assert.NoError(t, e)
	n.Write("/tmp/arash.txt", []byte("test"))
```