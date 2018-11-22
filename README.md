# dumper
A tiny utility to dump rootnroll demo containers logs to files for further analysis.

## Usage

Dump all rootnroll containers created more than `LIFETIME` minutes ago to `BASE_LOGS_DIR` directory:
```
./dumper BASE_LOGS_DIR LIFETIME
```

These are the containers that have the label `rootnroll.demo.name` applied to them.
