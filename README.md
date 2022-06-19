# Prune

prune or golang-prune is a command line utility to search a directory for files and/or directories to be pruned based on their relative age.

## Example

    prune --keep-daily 7 --dry-run .
    prune --keep-daily 14 --keep-monthly 6 --keep-yearly 1 --dry-run .
    prune --pattern "YYYY-MM-DDThh:mm:ss.sssZ" --keep-daily 7 --dry-run .
    prune --pattern "YYYY-MM-DDThh:mm:ss.sssZ" --keep-daily 7 --dry-run --json .

Alternative:

    prune -0 --keep-daily 14 --keep-monthly 6 --keep-yearly 1 . | xargs -0 rm

## Usage

tbd



This project was heavily inspired by [borg prune](https://borgbackup.readthedocs.io/en/stable/usage/prune.html)
