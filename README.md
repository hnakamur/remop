# remop

## What is this?

remop is a CLI tool to read 1Password secrets over an SSH tunnel.

It creates two Unix domain sockets: one on the remote server and the other
on the local host.

The Unix domain socket on the remote server is readable and writable only by
the remote user who logs in to the server by executing the `remop tunnel`
command.

You can run `remop read` on the remote server to get the value for the item
URL specified as a command-line argument.

Then an access request dialog appears in the 1Password app, and you can get
the value after you approve the request.

The config file on the local host has an allow list of item URLs, so you can
restrict which items `remop read` can access. This lowers the risk of
accidentally allowing an unexpected item.

## Setup

### Local host

Put the `remop` executable on your `PATH`.
Copy `remop.example.json` to `remop.json` and edit it for your environment.

### Remote server

Put the `remop` executable on your `PATH`.
Set the `REMOP_SOCKET` environment variable to the same value as
`local.socket_path` in the local host's `remop.json`.

## Usage

### Local host

Run `remop tunnel` to create an SSH tunnel to the target server.

```
remop tunnel <your_target_host>
```

### Remote server

Run `remop read` to get the value for the specified item URL.

```
remop read op://<your_vault_id>/<your_item_id>/<your_field_id>
```
