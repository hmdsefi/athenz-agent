### Athenz Sidecar Agent
`athenz-agent` contains athenz ZPE and ZPU utilities in Go language. ZPU will download the domains' policy files and store 
them into the filesystem. In other side, ZPE will use that policy files, and it will cache them into memory to use them as
fast as possible.

### How to install
For using Makefile you must edit this file and change some of variables as you want. 
* Set `BUILDPATH` to any directory that you want build this project into: 
``` 
BUILDPATH=/home/athenz/sidecar 
```

* Set ZMS url to `URL` variable:
```
URL=https://localhost:4443/
``` 

* We need to use `athenz-conf` utility to get `athenz.conf` file from ZMS server. To do that first you need to add ZMS 
server cert file to your local machine `/etc/ssl/certs` folder and run `sudo c_rehash`. Now you can set `athenz-conf` executable file directory
to `ATHENZCONF` variable:
```
ATHENZCONF=/home/athenz/athenz-zms-1.8.10/bin/linux
```
* Use `make sync` command to get all of project dependencies.
* Now you can build the project with `make build` command in command line.
* Run project with this command:
```bash
cd $(BUILDPATH) && ./athenz-agent
```

## Configuration
There are three different config file in this project:
* `agent.json`: You can use this config file to set root directory of build path. This file must always be beside of the generated executable file. Also, you can set config directory and zpe and zpu
 config file names in this file.
* `zpe.conf`: This file contains zpe configuration like caching duration, service name, domain name, and etc.
* `zpu.con`: This file contains zpu configuration like domain names that zpu must download them.