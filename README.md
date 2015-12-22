<div align="center">
     <a>
        <img width=710px src="http://animais.culturamix.com/blog/wp-content/gallery/fotos-aguia-1/Fotos-Aguia-2.png">
     </a>
     <br/>
     <b>A</b>(guia) <b>P</b>(rocess) <b>M</b>(anager)
     <br/><br/>
</div>
# APM - Aguia Process Manager
APM is a lightweight process manager written in Golang for Golang applications. It helps you keep your applications alive forever, to reload and start them from the source code.

Starting an application is easy:
```bash
$ ./apm bin app-name --source="github.com/topfreegames/apm"
```

This will basically compile your project source code and start it as a daemon in the background.

## Install APM

```bash
$ go get github.com/topfreegames/apm
```

## Start APM

```bash
$ ./apm serve --config-file="config/file/path.toml"
```

## Main features

### Commands overview

```bash
$ ./apm serve --config-file="config/file/path.toml"
$ ./apm serve-stop --config-file="config/file/path.toml"

$ ./apm bin app-name --source="github.com/topfreegames/apm"   # Compile, start, daemonize and auto restart application.
$ ./apm start app-name                                        # Start, daemonize and auto restart application.
$ ./apm stop app-name                                         # Stop application.
$ ./apm delete app-name                                       # Delete application forever.

$ ./apm status                                                # Display status for each app.

### Managing process via HTTP

You can also use all of the above commands via HTTP requests. Just set the flag ```--dns``` together with ```./apm serve``` and then you can use a remote client to start, stop, delete and query status for each app. 