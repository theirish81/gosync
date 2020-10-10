# GoSync

A simple program to easily update remote directories and files using plain HTTP.

When run in **client mode**, GoSync will monitor a given directory and when a file is created/updated/deleted,
the event and the content of the file is forwarded to a GoSync running in server mode.

When run in **server mode**, GoSync will bind to a specific directory and listen for updates from a GoSync instance
running in client mode.

## Configuration

Two files are necessary:

* **config.yml**: the configuration file. It needs to be placed in the same directory the program runs from

* **password file**: passwords. The file can be anywhere as long as config.yml references it correctly

The files need to be configured in different ways, whether you're running in client or server mode.

### Server

#### config.yml

*example:*

```yaml
server:
    port: 9999
    address: http://localhost
fs:
    root_dir: files/
    password_file: passwd
```

* `server.port`: the port the web server will bind to
* `server.address`: irrelevant
* `fs.root_dir`: the directory the server will bind to
* `fs.password_file`: the path to the password file

#### passwd

A file adhering to the **Apache htpasswd** format.

### Client

#### config.yml

*example:*

```yaml
server:
    port: 9999
    address: http://localhost
fs:
    root_dir: files/
    password_file: passwd
```

* `server.port`: the port of GoSync server the client will connect to
* `server.address`: the address of the GoSync server the client will connect to
* `fs.root_dir`: the directory the client will monitor
* `fs.password_file`: the path to the password file

#### passwd

A YAML file containing username and password:

*example:*

```yaml
username: foo
password: bar
```

## Running

Once the configuration is properly done, you can run GoSync, providing the mode in the command
line arguments.

The arguments can be:

* `-mode server`: runs GoSync in server mode
* `-mode client`: runs GoSync in client mode

**Nix example:*

```sh
./async -mode server
./async -mode server
```
