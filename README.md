# go-pivot-ssh

Simple Go-based tool to proxy local TCP connections
over an SSH tunnel.

Uses host entries found in `~/.ssh/config`
for authentication. If no `IdentifyFile` is
defined for a given host, user will be
prompted to enter password.

## How to use

Install (requires Go):
```
git clone https://github.com/jdolitsky/go-pivot-ssh.git 2>/dev/null && \
  (cd go-pivot-ssh && go build -o /usr/local/bin/pivot-ssh .) && \
  rm -rf go-pivot-ssh/
```

Usage:
```
pivot-ssh <remote_host> <local_listener> <forward_to>
```

## Example

For example, access service at `10.1.1.55:80`, which is only
reachable via `10.11.1.123` (which you have SSH access to).

In your `~/.ssh/config`, add info about the host:
```
$ cat ~/.ssh/config
Host 10.11.1.123
	User myuser
    Port 2222
```

Then start a local listener at `127.0.0.1:8080` which
will forward traffic to `10.1.1.55:80` over SSH tunnel:

```
$ pivot-ssh 10.11.1.123 127.0.0.1:8080 10.1.1.55:80
Enter SSH password: 
2021/11/11 09:46:43.239008 listening for new connections...
2021/11/11 09:46:45.711348 accepted connection
2021/11/11 09:46:45.711387 listening for new connections...
2021/11/11 09:46:46.248858 connected to 10.11.1.123:2222 (1 of 2)
2021/11/11 09:46:46.525407 connected to 10.1.1.55:80 (2 of 2)
```

## Credits

Powered by the following Go libraries:

- https://github.com/elliotchance/sshtunnel
- https://github.com/kevinburke/ssh_config
