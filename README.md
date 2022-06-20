# Prune

> WARNING: work in progress

*prune* is a command line utility to search a directory for files and/or directories to be pruned based on a retention policy. It is mainly aimed at pruning old backups created by using timestamped directories.

Currently *prune* does not have the goal to delete the files/directories to prune by itself but to create output that can be passed to tools dedicated to do this job (e.g. ` | xargs -0 rm -rf`)

## Example

Consider the following backup repository:

    /backup-repo/
    ├── 2000-01-01T00-00-00.000Z
    ├── 2000-02-01T00-00-00.000Z
    ├── 2000-03-01T00-00-00.000Z
    ├── 2000-04-01T00-00-00.000Z
    ├── 2000-05-01T00-00-00.000Z
    ├── 2000-06-01T00-00-00.000Z
    ├── 2000-07-01T00-00-00.000Z
    ├── 2000-08-01T00-00-00.000Z
    ├── 2000-09-01T00-00-00.000Z

Running prune and instructing it to keep the latest 6 monthly backups

    prune --keep-monthly 6 /backup-repo/

would result in the following output:

    /backup-repo/2000-04-01T00-00-00.000Z
    /backup-repo/2000-05-01T00-00-00.000Z
    /backup-repo/2000-06-01T00-00-00.000Z
    /backup-repo/2000-07-01T00-00-00.000Z
    /backup-repo/2000-08-01T00-00-00.000Z
    /backup-repo/2000-09-01T00-00-00.000Z

or with the `--verbose` option:

    /backup-repo/2000-01-01T00-00-00.000Z: prune
    /backup-repo/2000-02-01T00-00-00.000Z: prune
    /backup-repo/2000-03-01T00-00-00.000Z: prune
    /backup-repo/2000-04-01T00-00-00.000Z: keep
    /backup-repo/2000-05-01T00-00-00.000Z: keep
    /backup-repo/2000-06-01T00-00-00.000Z: keep
    /backup-repo/2000-07-01T00-00-00.000Z: keep
    /backup-repo/2000-08-01T00-00-00.000Z: keep
    /backup-repo/2000-09-01T00-00-00.000Z: keep


## Usage

### xargs

List files/directories:

    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory

Prune files/directories:

    prune -0 --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory | xargs rm -rf

- [✗] spaces
- [?] globs

### xargs with NUL

List files/directories:

    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory

Prune files/directories:

    prune -0 --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory | xargs -0 rm

- [✓] spaces
- [?] globs

based on https://stackoverflow.com/a/16758699/548020

### Shell Script

    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory > prune.list
    while IFS= read -r file ; do rm -rf -- "$file" ; done < prune.list

- [✓] spaces
- [✓] globs

based on https://stackoverflow.com/a/21848934/548020


## Testing

### Test-Repo

Create a test repo:

    go build ./cmd/test-repo && ./test-repo ~/temp/test-repo


## Roadmap

### v0.1

- Support basic `--keep-*` flags
- Print files to prune to *stdout*

### v0.2

- Support `-0|--null` flag to write null terminated list of files to *stdout*

### v0.3

- Support `-p|--pattern` flag?


## Ideas

- Allow the pattern of the timestamped directories to be defined using a CLI option like `--pattern "YYYY-MM-DDThh:mm:ss.sssZ"`
- Allow outputting JSON for better scripting support by introducing a `--json` flag
- Introduces option to write a list of files to prune to a file, so it can be reviewed and used as input to actually delete the files/directories. A suitable format is yet to be discovered.
- Perform the actual delete operation (https://pkg.go.dev/os#RemoveAll), but also introduce a `--dry-run` flag. Not sure if this is the way to go though



## Inspiration

This project was heavily inspired by [borg prune](https://borgbackup.readthedocs.io/en/stable/usage/prune.html)
