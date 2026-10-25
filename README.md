# remop

## What is this?

remop is a CLI to read 1Password secrets over a SSH tunnel.

It creates two unix domain sockets, one on the remote server and the other
on the local host.

The unix domain socket on the remote server is only readable and writable by
the remote user who logins to the server by executing `remop tunnel` command.

You can run `remop read` on the remote server to get the value of the item URL
specified in the command line argument.

Then the access request dialog of the 1Password app is prompted, and you can
get the value after you accept the request.

The config file for the local host has the allow list of item URLs, so you
can limit items which `remop read` can access. This lowers the risk you allow
the unexpected item by mistake.

## Setup

### local host
 
Put the `remop` executable on the PATH.
Copy `remop.example.json` to `remop.json` and edit it for your environment.

### remote server

Put the `remop` executable on the PATH.
Set the `REMOP_SOCKET` environment variable to the same value as
`local.socket_path` in `remop.json` on the local host.

## Usage

### local host

Run `remop tunnel` to create the SSH tunnel to the target server.

```
remop tunnel _YOUR_TARGET_HOST_HERE_
```

### remote server

Run `remop read` to get the value by the specified item URL.

```
remop read op://<_your_vault_id_here>/<_your_item_id_here_>/<_field_id_here_>
```
