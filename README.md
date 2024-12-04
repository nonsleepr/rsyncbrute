# rsyncbrute

This Go project implements a tool to perform brute-force attacks against rsync servers. It attempts to authenticate with different username/password combinations against a specified rsync share.

The project is inspired by and based on the following Nmap script:

[https://github.com/nmap/nmap/blob/master/scripts/rsync-brute.nse](https://github.com/nmap/nmap/blob/master/scripts/rsync-brute.nse)

The rsync protocol description can be found here:

[https://github.com/RsyncProject/rsync/blob/master/csprotocol.txt](https://github.com/RsyncProject/rsync/blob/master/csprotocol.txt)

## Installation

```bash
go install github.com/nonsleepr/rsyncbrute@latest
```

## Usage

```
rsyncbrute -host <host> -port <port> -share <share> -usernames <usernames_file> -passwords <passwords_file> [-max-connections <max_connections>]
```

| Flag               | Description                                                                      | Default |
| ------------------ | -------------------------------------------------------------------------------- | ------- |
| `-host`            | Rsync server hostname or IP address.                                             |         |
| `-port`            | Rsync server port.                                                               | 873     |
| `-share`           | Rsync share name.                                                                |         |
| `-usernames`       | Path to the file containing usernames (one username per line).                   |         |
| `-passwords`       | Path to the file containing passwords (one password per line).                   |         |
| `-max-connections` | Maximum number of parallel connections.                                          | 200     |


## Performance

The tool achieves a performance of approximately 2,500,000 attempts per hour with 5000 maximum connections. Using IP addresses instead of hostnames can further improve performance.
