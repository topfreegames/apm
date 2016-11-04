<div align="center">
     <a>
        <img height="400px" width="400px" src="https://cloud.githubusercontent.com/assets/5413599/12247882/2fcb1ca6-b89d-11e5-933e-efade26acf13.jpg">
     </a>
     <br/>
     <b>A</b>(guia) <b>P</b>(rocess) <b>M</b>(anager)
     <br/><br/>
</div>
# APM - Aguia Process Manager
APM is a lightweight process manager written in Golang for Golang applications. It helps you keep your applications alive forever, to reload and start them from the source code.

[![ReportCard](http://goreportcard.com/badge/topfreegames/apm)](http://goreportcard.com/badge/topfreegames/apm)
[![GoDoc](https://godoc.org/github.com/topfreegames/apm?status.svg)](https://godoc.org/github.com/topfreegames/apm)

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

In order to properly use APM, you always need to start a server. This will be changed in the next version, but in the meantime you need to run the command bellow to start using APM.
```bash
$ apm serve
```
If not config file is provided, it will default to a folder '.apmenv' where apm is first started.

## Stop APM

```bash
$ apm serve-stop
```

## Starting a new application
If it's the first time you are starting a new golang application, you need to tell APM to first build its binary. Then you need to first run:
```bash
$ apm bin app-name --source="github.com/yourproject/project"
```

This will automatically compile, start and daemonize your application. If you need to later on, stop, restart or delete your app from APM, you can just run normal commands using the app-name you specified. Example:
```bash
$ apm stop app-name
$ apm restart app-name
$ apm delete app-name
```

## Main features

### Commands overview

```bash
$ apm serve --config-file="config/file/path.toml"
$ apm serve-stop --config-file="config/file/path.toml"

$ apm bin app-name --source="github.com/topfreegames/apm"   # Compile, start, daemonize and auto restart application.
$ apm start app-name                                        # Start, daemonize and auto restart application.
$ apm restart app-name                                      # Restart a previously saved process
$ apm stop app-name                                         # Stop application.
$ apm delete app-name                                       # Delete application forever.

$ apm save                                                  # Save current process list
$ apm resurrect                                             # Restore previously saved processes

$ apm status                                                # Display status for each app.

### Managing process via HTTP

You can also use all of the above commands via HTTP requests. Just set the flag ```--dns``` together with ```./apm serve``` and then you can use a remote client to start, stop, delete and query status for each app. 
