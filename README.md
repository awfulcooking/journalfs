# journalfs

[FUSE filesystem](https://en.wikipedia.org/wiki/Filesystem_in_Userspace) presenting the systemd journal as directories of .log files.

## Build

1. `apt install libsystemd-dev golang fuse3`
1. `echo "user_allow_other" | sudo tee -a /etc/fuse.conf`
1. `make`

## Install

1. `make install`

2. Create a system user and group

```bash
sudo useradd --system --user-group journalfs
sudo usermod -aG systemd-journal journalfs
sudo usermod -aG journalfs $USER
```

  Note, you'll have to start a new shell or desktop session to activate $USER's group membership.

3. Create a mount point, set ownership and permissions:

```
sudo mkdir /journal
sudo chown journalfs:journalfs /journal
sudo chmod 750 /journal
```

4. Install the service:

```
sudo cp contrib/journalfs.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now journalfs
```

### Directories

```
$ ls -1p /journal
automount/
device/
mount/
path/
scope/
service/
slice/
socket/
swap/
target/
timer/
```

Directories will contain .log files.

## Contributions welcome

Help flesh this out, or fix bugs. See Issues for ideas.

## License

GPLv3
