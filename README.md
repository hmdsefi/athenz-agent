[![Go Report Card](https://goreportcard.com/badge/github.com/hamed-yousefi/athenz-agent)](https://goreportcard.com/report/github.com/hamed-yousefi/athenz-agent)
### Athenz Sidecar Agent
You're a Yahoo Athenz fan, so this app is created for you. `athenz-agent` contains athenz ZPE and ZPU utilities in Go language.
ZPU will download the domains' policy files and store
them into the filesystem. In other side, ZPE will use that policy files, and it will cache them into memory to use them as
fast as possible.

![flow](https://github.com/hamed-yousefi/athenz-agent/blob/master/docs/images/auth_flow.png)

Athenz agent exposes two APIs:
- CheckAccessWithToken
- GetServiceToken

**CheckAccessWithToken:** Accepts three arguments including client service RoleToken, provider service access,
and provider service resource.

**GetServiceToken:** Has no input argument. It returns RoleToken az result.


### How to install
For using Makefile you must edit this file and change some variables as you want.
* Set `BUILDPATH` to any directory that you want build this project into:
``` 
BUILDPATH=/home/athenz/sidecar 
```

* Set ZMS url to `URL` variable:
```
URL=https://localhost:4443/
``` 

* You need to use `athenz-conf` utility to get `athenz.conf` file from ZMS server. To do that first you need to add ZMS
  server cert file to your local machine `/etc/ssl/certs` folder and run `sudo c_rehash`. Now you can set `athenz-conf` executable file directory
  to `ATHENZCONF` variable:
```
ATHENZCONF=/home/athenz/athenz-zms-1.8.10/bin/linux
```
* Use `make sync` command to get all project dependencies.
* Now you can build the project with `make build` command in command line.
* Run project with this command:
```bash
cd $(BUILDPATH) && ./athenz-agent
```

### Configuration
All default configuration files placed in `build/config` path.

## License
MIT License, please see [LICENSE](https://github.com/hamed-yousefi/athenz-agent/blob/master/LICENSE) for details.