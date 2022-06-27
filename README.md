# Prune

[![Go Report Card](https://goreportcard.com/badge/github.com/codezombiech/prune)](https://goreportcard.com/report/github.com/codezombiech/prune)

> WARNING: work in progress

*prune* is a command line utility to search a directory for files and/or directories to be pruned based on a retention policy. It is mainly aimed at pruning old backups created by using timestamped directories.

Currently *prune* does not have the goal to delete the files/directories to prune by itself but to create output that can be passed to tools dedicated to do this job (e.g. ` | xargs -0 rm -rf`)

## Example

Consider a backup strategy where you backup your files as compressed tar archives. Each backup isstored in a directory using the following pattern

    /backups/<timestamp>/

where `<timestamp>` represents the date the backup was created using the `%Y-%m-%dT%H-%M-%S%z` format

Running this strategy every day from 2000-01-01 to 2000-01-03 would result in the following directory tree:

    /backups
    ├── 2000-01-01T00-00-00Z
    │   ├── backup-2000-01-01T00-00-00Z.tar.gz
    │   └── backup-2000-01-01T00-00-00Z.tar.gz.sha256sum
    ├── 2000-01-02T00-00-00Z
    │   ├── backup-2000-01-02T00-00-00Z.tar.gz
    │   └── backup-2000-01-02T00-00-00Z.tar.gz.sha256sum
    └── 2000-01-03T00-00-00Z
        ├── backup-2000-01-03T00-00-00Z.tar.gz
        └── backup-2000-01-03T00-00-00Z.tar.gz.sha256sum

Running `prune` and instructing it to keep the latest 2 daily backups

    prune --pattern '%Y-%m-%dT%H-%M-%S%z' --keep-daily 2 /backups

would result in the following output, reporting paths to directories that should be pruned:

    /backups/2000-01-01T00-00-00Z

or with the `--verbose` option:

    /backups/2000-01-01T00-00-00Z: prune
    /backups/2000-01-02T00-00-00Z: keep
    /backups/2000-01-03T00-00-00Z: keep
    Total count: keep: 2, prune: 1


## Usage

### Prune

    prune [--verbose|-v] [--pattern <pattern>]
        [--keep-daily|-d <keep-count>] [--keep-monthly|-m <keep-count>] [--keep-yearly|-y <keep-count>]
        <directory>

where
- `<pattern>`: pattern to use to parse the date/time from the directory name
- `<keep-count>`: number of directories to keep
- `<directory>`: path to directory to scan for directories to prune

Without the `--verbose|-v` flag, this will list all directories to be pruned.
Setting the `--verbose|-v` flag will list all directories indicating if they would be kept/deleted and basic statistics


### Prune and Delete

This section describes strategies how the output of `prune` can be used to eventually  delete files/directories to be pruned.

#### xargs

List files/directories:

    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory

Prune files/directories:

    prune -0 --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory | xargs rm -rf

Works with:
- [✗] spaces
- [?] globs

#### xargs with NUL

List files/directories:

    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory

Prune files/directories:

    prune -0 --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory | xargs -0 rm

Works with:
- [✓] spaces
- [?] globs

based on https://stackoverflow.com/a/16758699/548020

#### Shell Script (preferred)

    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 /path/to/directory > prune.list
    while IFS= read -r file ; do rm -rf -- "$file" ; done < prune.list

Works with:
- [✓] spaces
- [✓] globs

based on https://stackoverflow.com/a/21848934/548020


## Testing

### Test-Repo

Create a test repo:

    go run ./cmd/test-repo --pattern '%Y-%m-%dT%H-%M-%SZ' ~/temp/test-repo 2000-12-01...2001-01-31

Run prune against test repo

    go run . -d 3 -m 2 -y 1 ~/temp/test-repo


## Roadmap

### v0.1

- Support basic `--keep-*` options
- Print directories to prune to *stdout*
- Print all directories with keep/status to *stdout* when `--verbose` flag is set
- Support `--pattern` option, defining the pattern used to parse the date of the timestamped directories

### v0.2

- Evaluate if `<pattern>` should be an option or a mandatory argument. If option, what would be a sane default?
- Evaluate what is more clear, `--pattern` or `--format`

### v0.3

- Support `-0|--null` flag to write null terminated list of files to *stdout*


## Ideas

- Allow the pattern of the timestamped directories to be defined using a CLI option like `--pattern "YYYY-MM-DDThh:mm:ss.sssZ"`
- Allow outputting JSON for better scripting support by introducing a `--json` flag
- Introduces option to write a list of files to prune to a file, so it can be reviewed and used as input to actually delete the files/directories. A suitable format is yet to be discovered.
- Perform the actual delete operation (https://pkg.go.dev/os#RemoveAll), but also introduce a `--dry-run` flag. Not sure if this is the way to go though
- Allow passing multiple directories, which all feed into a union set of backups. This would allow pruning a single type of backups being created using different backup strategies (versions of a backup script)



## Inspiration

This project was heavily inspired by [borg prune](https://borgbackup.readthedocs.io/en/stable/usage/prune.html)
