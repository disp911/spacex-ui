package docsite

// A kind is what the project of a site does: the words its pages are written
// from. Texts use the markup of render and the fillings of expand; besides
// the kind's own words every text can use {name}, {bin}, {env}, {conf},
// {ext}, {domain} and {ver}.
type kind struct {
	id string
	// what completes "{name} is ...", with no full stop.
	what     []string
	taglines []string
	// words are name stems that suit the kind.
	words    []string
	features []feature
	commands []command
	options  []option
	guides   []guide
	faq      []qa
	changes  []string
	quick    string
	// flow is the diagram on the front and introduction pages: where data
	// comes from, what the project does with it, where it goes.
	flow [3]string
	// demo is the terminal on the front page: a command and what it prints.
	demo []string
	// daemon is the command the service unit runs; without one the project
	// runs from a timer.
	daemon string
	// snips are pieces of configuration the texts show, in YAML and TOML;
	// those named in example are added to the example file.
	snips   map[string][2]string
	example []string
}

type feature struct{ title, text, icon string }

type command struct {
	name, args, summary, text string
	flags                     []flag
}

type flag struct{ name, value, def, desc string }

// An option of the configuration file. ex is the value shown in the example
// file, written so that it reads the same in YAML and TOML; options without
// one are left out of it.
type option struct{ key, typ, def, ex, desc string }

type guide struct{ slug, title, nav, body string }

type qa struct{ q, a string }

var kinds = []*kind{backupKind, logsKind, jobsKind, monitorKind}

var backupKind = &kind{
	id: "backup",
	what: []string{
		"a {small|fast|dependable} command-line tool for {incremental|deduplicated}, encrypted backups to local disks, SFTP servers and S3-compatible storage",
		"a backup program that keeps encrypted, deduplicated snapshots of your files on {a disk, an SFTP server or any S3-compatible bucket|local disks, SFTP servers and object storage}",
	},
	taglines: []string{
		"Backups you can actually restore",
		"Encrypted backups without the ceremony",
		"Back up everything. Store each byte once.",
		"Snapshots of your data, kept safe and small",
	},
	words: []string{"Vault", "Keep", "Cache", "Hoard", "Ark", "Archive", "Burrow", "Cellar"},
	features: []feature{
		{"Deduplication", "Files are split into content-defined chunks. A chunk that is already in the repository is never {uploaded|stored} again, {even if it moved to another file|however many files share it}.", "layers"},
		{"Encrypted by default", "Everything leaves the machine encrypted with AES-256-GCM. The key is derived from your password and never reaches the storage.", "lock"},
		{"Many backends", "Local folders, SFTP and anything that speaks the S3 API, with the same commands for all of them.", "cloud"},
		{"Fast restores", "Restore a single file or a whole tree from any snapshot, or mount a snapshot and browse it like a folder.", "bolt"},
		{"Retention policies", "Keep hourly, daily, weekly and monthly snapshots and let `prune` remove {the rest|everything else}.", "clock"},
		{"Verifiable", "`check` reads the repository back and proves that every snapshot can still be restored.", "shield"},
	},
	commands: []command{
		{"init", "", "Create a new repository.", "The repository is created at the location given by `--repo` or the `repository` option. You are asked for a password unless `password_file` is set.", []flag{
			{"--repo", "URL", "", "Where to create the repository."},
			{"--compression", "MODE", "auto", "Default compression for new data: `off`, `auto` or `max`."},
		}},
		{"backup", "PATH...", "Make a snapshot of the given paths.", "Only data that is not in the repository yet is read from disk and uploaded. Files whose size and modification time did not change since the last snapshot are not read at all.", []flag{
			{"--exclude", "PATTERN", "", "Leave out paths matching the pattern. May be given more than once."},
			{"--exclude-file", "FILE", "", "Read exclude patterns from a file, one per line."},
			{"--tag", "TAG", "", "Add a tag to the snapshot."},
			{"--one-file-system", "", "false", "Do not cross into other mounted filesystems."},
			{"--dry-run", "", "false", "Show what would be backed up without writing anything."},
		}},
		{"snapshots", "", "List the snapshots in the repository.", "", []flag{
			{"--tag", "TAG", "", "Only list snapshots with this tag."},
			{"--host", "HOST", "", "Only list snapshots made on this host."},
			{"--json", "", "false", "Print the list as JSON."},
		}},
		{"restore", "SNAPSHOT", "Restore files from a snapshot.", "`SNAPSHOT` is a snapshot ID, a unique prefix of one, or `latest`.", []flag{
			{"--target", "DIR", "", "Folder to restore into. Required."},
			{"--include", "PATTERN", "", "Only restore paths matching the pattern."},
			{"--verify", "", "false", "Read the restored files back and compare them with the snapshot."},
		}},
		{"prune", "", "Remove snapshots the retention policy no longer keeps.", "Data that is no longer used by any snapshot is deleted from the repository afterwards.", []flag{
			{"--keep-daily", "N", "7", "Keep the last snapshot of each of the last N days."},
			{"--keep-weekly", "N", "4", "Keep the last snapshot of each of the last N weeks."},
			{"--keep-monthly", "N", "12", "Keep the last snapshot of each of the last N months."},
			{"--dry-run", "", "false", "Show what would be removed."},
		}},
		{"check", "", "Verify the structure and data of the repository.", "Without flags only the index and the snapshot trees are checked, which is quick. `--read-data` downloads and decrypts every chunk.", []flag{
			{"--read-data", "", "false", "Read and verify all data."},
			{"--read-data-subset", "N%", "", "Read and verify a random part of the data."},
		}},
		{"mount", "SNAPSHOT DIR", "Mount a snapshot as a read-only folder.", "Available on Linux and macOS. The folder stays mounted until you press Ctrl+C or unmount it.", []flag{
			{"--allow-other", "", "false", "Let other users read the mounted folder."},
		}},
	},
	options: []option{
		{"repository", "string", "", `"s3://backups/web-01"`, "Where snapshots are stored: a local path, `sftp://user@host/path` or `s3://bucket/prefix`."},
		{"password_file", "string", "", `"/etc/{bin}/password"`, "File holding the repository password. Takes precedence over the `{env}PASSWORD` variable."},
		{"compression", "string", "auto", "", "`off`, `auto` or `max`. With `auto`, data that does not compress is stored as it is."},
		{"exclude", "list", "[]", `["*.tmp", "/home/*/.cache"]`, "Patterns of paths to leave out, in gitignore syntax."},
		{"one_file_system", "bool", "false", "", "Do not cross into other mounted filesystems."},
		{"retention.keep_hourly", "int", "24", "", "Hourly snapshots to keep."},
		{"retention.keep_daily", "int", "7", "14", "Daily snapshots to keep."},
		{"retention.keep_weekly", "int", "4", "8", "Weekly snapshots to keep."},
		{"retention.keep_monthly", "int", "12", "", "Monthly snapshots to keep."},
		{"upload.connections", "int", "5", "", "Parallel connections to the backend."},
		{"upload.limit", "string", "", `"20MiB/s"`, "Bandwidth limit, such as `20MiB/s`. Empty means no limit."},
		{"cache_dir", "string", "~/.cache/{bin}", "", "Local cache of repository metadata."},
		{"log.level", "string", "info", "", "`debug`, `info`, `warn` or `error`."},
	},
	guides: []guide{
		{"repositories", "Repositories", "", `
A repository is where {name} keeps your snapshots. It is a folder of encrypted files, and it does not matter whether that folder is on a local disk, on an SFTP server or in an S3 bucket: {name} only ever adds files to it and, during ` + "`prune`" + `, deletes files it no longer needs.

## Choosing a location

The location is a URL. A plain path is a local folder:

` + "```shell" + `
$ {bin} init --repo /mnt/backup/{bin}
$ {bin} init --repo sftp://backup@nas.local/srv/{bin}
$ {bin} init --repo s3://my-bucket/laptop
` + "```" + `

For S3 the credentials come from the usual ` + "`AWS_ACCESS_KEY_ID`" + ` and ` + "`AWS_SECRET_ACCESS_KEY`" + ` variables. Set ` + "`AWS_ENDPOINT_URL`" + ` to use a provider other than Amazon.

## Sharing a repository

Several machines can back up into the same repository. Chunks that two machines have in common are stored once, which saves a lot of space when they run the same operating system. Every snapshot records the host it was made on, so ` + "`{bin} snapshots --host`" + ` shows the snapshots of one machine.

> Only one ` + "`prune`" + ` may run on a repository at a time. {name} takes a lock in the repository and refuses to start a second one.

## Moving a repository

A repository is just files. Copy them to the new location with any tool you like, for example ` + "`rsync -a`" + ` or ` + "`rclone copy`" + `, and point the ` + "`repository`" + ` option at the copy.
`},
		{"retention", "Snapshots and retention", "Retention", `
Every run of ` + "`{bin} backup`" + ` makes a snapshot: a complete picture of the paths you gave it at that moment. Because unchanged data is stored only once, a snapshot that differs little from the one before costs almost nothing.

## Listing snapshots

` + "```shell" + `
$ {bin} snapshots
ID        Time                 Host     Tags    Paths
4f1c9a2e  {d1} 02:00:13  web-01   daily   /etc /srv
a83b51d0  {d2} 02:00:09  web-01   daily   /etc /srv
` + "```" + `

## Keeping the right ones

Snapshots pile up, so sooner or later you want to remove old ones. The retention policy says which to keep; everything else is removed by ` + "`{bin} prune`" + `:

` + "```{ext}" + `
{retention}
` + "```" + `

A snapshot kept by any rule stays. With the policy above you can go back to any of the last two weeks day by day, and to any month of the last year.

> Run ` + "`{bin} prune --dry-run`" + ` first when you change the policy. It lists the snapshots that would go without removing anything.
`},
		{"restoring", "Restoring data", "", `
There are two ways to get data back: ` + "`restore`" + ` writes files from a snapshot to a folder, and ` + "`mount`" + ` shows snapshots as a read-only folder you can browse.

## Restoring files

` + "```shell" + `
$ {bin} restore latest --target /tmp/restore --include /etc/nginx
` + "```" + `

` + "`latest`" + ` is the newest snapshot. Any snapshot ID, or a prefix of one that is unique, works as well. Files are written with their original permissions, owners and modification times when {name} runs with enough rights to set them.

## Browsing snapshots

` + "```shell" + `
$ mkdir /mnt/snapshots
$ {bin} mount latest /mnt/snapshots
` + "```" + `

Mounting needs FUSE on Linux and macFUSE on macOS. It is the easiest way to find the one file you deleted last week.

## Checking a restore

Pass ` + "`--verify`" + ` and {name} reads every restored file back and compares it with the snapshot. It takes about twice as long and is worth it for anything you are going to rely on.
`},
		{"encryption", "How encryption works", "Encryption", `
Everything {name} writes to a repository is encrypted and authenticated. Whoever can read the storage sees only the number and size of files, never their names or contents.

## Keys

When a repository is created, {name} generates a random master key. The key is stored in the repository, encrypted with a key derived from your password with scrypt. Changing the password only re-encrypts that one file.

## Data

Chunks, trees and snapshots are encrypted with AES-256-GCM, each with a random nonce. The chunk boundaries are found with a rolling hash keyed with a secret from the master key, so their sizes give nothing away about the files.

> If you lose the password, the data cannot be recovered. Keep a copy of it somewhere other than the machine you back up.
`},
	},
	faq: []qa{
		{"Can I back up to two places at once?", "Create two repositories and run `{bin} backup` once for each. Each repository is complete on its own, so losing one does not affect the other."},
		{"What happens if a backup is interrupted?", "Nothing is lost. Data uploaded before the interruption stays in the repository and is reused by the next run, which then finishes much faster."},
		{"How much space does the local cache need?", "Roughly 1% of the size of the repository. The cache only speeds things up; you can delete it at any time."},
		{"Can I change the repository password?", "Yes, with `{bin} key passwd`. Only the key file is re-encrypted, which takes a moment however large the repository is."},
		{"Does {name} back up open databases correctly?", "It reads files as they are, like any file-level tool. Dump the database first, or back up a filesystem snapshot, to get a consistent copy."},
	},
	changes: []string{
		"Faster snapshot listing on repositories with thousands of snapshots",
		"`restore` keeps extended attributes and ACLs",
		"The S3 backend retries throttled requests with backoff",
		"Fixed a crash in `prune` when a snapshot referred to a missing tree",
		"`check --read-data-subset` accepts a percentage",
		"Lower memory use when backing up folders with millions of files",
		"New `--json` output for `snapshots`",
		"Mounting works with macFUSE 4",
		"Fixed restoring symlinks that point outside the target folder",
		"The progress line shows the upload rate",
		"`backup` skips sockets and named pipes instead of failing",
		"Fixed a race that could leave a stale lock after Ctrl+C",
		"SFTP connections are kept alive during long scans",
		"`--exclude-file` ignores comment lines",
		"Fixed the size shown for sparse files",
		"Chunk index loads in half the time",
		"`prune` reports how much space it freed",
		"Added `--host` to `snapshots` and `prune`",
	},
	quick: `
This page takes you from nothing to a first snapshot and a first restore in about five minutes. It assumes {name} is [installed](/docs/installation/).

## Create a repository

` + "```shell" + `
$ {bin} init --repo /srv/backup
Enter a password for the new repository:
Repeat the password:
Created repository 9c1e4b2a at /srv/backup
` + "```" + `

Remember the password. Without it the repository cannot be read, by you or anyone else.

## Make a snapshot

` + "```shell" + `
$ export {env}REPOSITORY=/srv/backup
$ {bin} backup /etc /home
scanned 18,402 files, 2.31 GiB in 0:41
added 1.12 GiB to the repository (2.31 GiB before compression and deduplication)
snapshot 4f1c9a2e saved
` + "```" + `

Run the same command again. This time it takes a few seconds, because nothing has changed.

## Restore a file

` + "```shell" + `
$ {bin} restore latest --target /tmp/restore --include /etc/hosts
restored 1 file, 230 B
` + "```" + `

## Next steps

- Put the repository and password into [the configuration file](/docs/configuration/) so you do not have to repeat them.
- Set up a [retention policy](/docs/retention/) and run ` + "`{bin} backup`" + ` from a timer, as shown in [Deployment](/docs/deployment/).
`,
	flow: [3]string{"Your files", "Repository", "Disk, SFTP or S3"},
	demo: []string{
		"$ {bin} backup /etc /home",
		"scanned 18,402 files, 2.31 GiB in 0:41",
		"added 1.12 GiB to the repository",
		"snapshot 4f1c9a2e saved",
	},
	snips: map[string][2]string{
		"retention": {"retention:\n  keep_daily: 14\n  keep_weekly: 8\n  keep_monthly: 12", "[retention]\nkeep_daily = 14\nkeep_weekly = 8\nkeep_monthly = 12"},
	},
}

var logsKind = &kind{
	id: "logs",
	what: []string{
		"a {lightweight|small} log shipper that tails files and the systemd journal, parses each line and sends it on to {HTTP endpoints, files and other collectors|whatever stores your logs}",
		"a {fast|single-binary} agent that collects logs from files and journald, turns them into structured events and ships them where you need them",
	},
	taglines: []string{
		"Ship logs, not headaches",
		"Every log line, structured and on its way",
		"A log shipper that stays out of the way",
		"From log files to structured events",
	},
	words: []string{"Tail", "Stream", "Creek", "Brook", "Courier", "Relay", "Runner", "Ferry"},
	features: []feature{
		{"Tiny footprint", "A single static binary that idles at a few megabytes of memory, so it fits on the smallest {VM|machine}.", "bolt"},
		{"Files and journald", "Tails files with glob patterns, follows rotation and reads the systemd journal with its fields.", "terminal"},
		{"Structured events", "Parses JSON, logfmt and regular expressions into fields you can filter and route on.", "filter"},
		{"Disk buffer", "When an output is down, events wait in a disk buffer and are sent in order once it is back.", "database"},
		{"Many outputs", "HTTP with batching and gzip, files, syslog and other {name} instances, all at the same time.", "plug"},
		{"At-least-once", "Positions are saved only after an output confirms delivery. A restart never loses a line.", "shield"},
	},
	commands: []command{
		{"run", "", "Start collecting and shipping logs.", "This is what the service unit runs. It reads the configuration file, starts all inputs and outputs, and keeps running until it receives SIGTERM.", []flag{
			{"--once", "", "false", "Read what is in the inputs now, ship it and exit."},
		}},
		{"validate", "", "Check the configuration file and exit.", "Every input, parser and output is checked as if {name} were starting, without opening any files or connections.", nil},
		{"tail", "INPUT", "Print the events of one input as they arrive.", "Useful to see what a parser makes of your lines before they are shipped anywhere.", []flag{
			{"--raw", "", "false", "Print the lines as read, before parsing."},
			{"--lines", "N", "10", "Start with the last N lines."},
		}},
		{"replay", "FILE", "Ship a file from the beginning.", "The file is read once and its lines go through the same parsers and outputs as the input they would match.", []flag{
			{"--input", "NAME", "", "Treat the file as part of this input."},
		}},
		{"positions", "", "Show how far each input has been read.", "", []flag{
			{"--reset", "INPUT", "", "Forget the position of an input, so it is read again from the start."},
		}},
	},
	options: []option{
		{"inputs.files", "list", "[]", `["/var/log/nginx/*.log", "/var/log/app/*.log"]`, "Glob patterns of files to tail."},
		{"inputs.journald", "bool", "false", "true", "Read the systemd journal."},
		{"inputs.start_at", "string", "end", "", "Where to begin with a file never seen before: `beginning` or `end`."},
		{"parse.format", "string", "auto", `"json"`, "`json`, `logfmt`, `regex` or `auto`, which tries JSON and logfmt."},
		{"parse.pattern", "string", "", "", "Regular expression with named groups, used when the format is `regex`."},
		{"outputs.http.url", "string", "", `"https://logs.example.com/ingest"`, "Endpoint to POST batches of events to."},
		{"outputs.http.batch_size", "int", "500", "", "Events per request."},
		{"outputs.http.compress", "bool", "true", "", "Compress requests with gzip."},
		{"outputs.file.path", "string", "", "", "Write events to this file as JSON lines."},
		{"buffer.path", "string", "/var/lib/{bin}/buffer", "", "Folder of the disk buffer."},
		{"buffer.max_size", "string", "256MiB", `"1GiB"`, "Largest size the buffer may grow to before the oldest events are dropped."},
		{"flush_interval", "duration", "1s", "", "Longest time an event waits before its batch is sent."},
		{"labels", "map", "{}", "", "Fields added to every event, such as the environment or the host role."},
		{"log.level", "string", "info", "", "`debug`, `info`, `warn` or `error`."},
	},
	guides: []guide{
		{"inputs", "Inputs", "", `
Inputs are where events come from. {name} can tail files and read the systemd journal, and it can do both at once.

## Files

` + "```{ext}" + `
{files}
` + "```" + `

Patterns are expanded every few seconds, so files created later are picked up without a restart. A file is followed through rotation whether it is renamed or truncated: {name} keeps reading the old file until it has reached its end, then switches to the new one.

## The journal

With ` + "`inputs.journald`" + ` set, every journal entry becomes an event. The message is the ` + "`message`" + ` field and the journal's own fields, such as ` + "`_SYSTEMD_UNIT`" + ` and ` + "`PRIORITY`" + `, are kept in lower case.

> Reading the journal needs membership in the ` + "`systemd-journal`" + ` group. The [service unit](/docs/deployment/) adds it.

## Where reading starts

A file {name} has never seen is read from its end, so starting it on a busy machine does not resend old logs. Set ` + "`inputs.start_at`" + ` to ` + "`beginning`" + ` to read whole files instead. Positions are kept in ` + "`/var/lib/{bin}`" + ` and survive restarts.
`},
		{"parsing", "Parsing lines", "Parsing", `
Every line is turned into an event with fields. How depends on ` + "`parse.format`" + `.

## JSON and logfmt

Lines that are JSON objects become events with the same fields. Logfmt lines, such as ` + "`level=info msg=\"started\" took=12ms`" + `, become fields in the same way. The default format ` + "`auto`" + ` tries both and keeps the line as a plain message when neither fits.

## Regular expressions

For anything else, write a pattern with named groups:

` + "```{ext}" + `
{pattern}
` + "```" + `

Each named group becomes a field. Lines that do not match are shipped as a plain message with the field ` + "`parse_error`" + ` set, so nothing disappears silently.

## Trying a parser

` + "`{bin} tail`" + ` prints the events of an input as they are parsed:

` + "```shell" + `
$ {bin} tail nginx --lines 2
{"time":"{d1}T10:02:11Z","status":200,"path":"/","took_ms":3}
{"time":"{d1}T10:02:12Z","status":404,"path":"/favicon.ico","took_ms":1}
` + "```" + `
`},
		{"outputs", "Outputs and buffering", "Outputs", `
Outputs are where events go. Every event goes to every output, and each output has its own buffer, so a slow output never holds up the others.

## HTTP

The HTTP output sends events in batches as a JSON array, compressed with gzip. A batch is sent when it is full or when ` + "`flush_interval`" + ` has passed, whichever comes first. Any answer in the 2xx range counts as delivered; 429 and 5xx answers are retried with backoff.

## The disk buffer

While an output cannot deliver, its events are written to the disk buffer. When the output is back they are sent in the order they arrived. The buffer is limited by ` + "`buffer.max_size`" + `; when it is full, the oldest events are dropped and the count of dropped events is logged once a minute.

> Positions in the inputs are saved only after an event has reached every output or the disk buffer. After a crash some events may be sent twice, but none are lost.
`},
		{"labels", "Labels and routing", "Labels", `
Labels are fields {name} adds to every event. They tell you later where an event came from without having to put that into every application's log format.

` + "```{ext}" + `
{labels}
` + "```" + `

## Labels per input

Labels can also be set on a single input, and they are added after the global ones, so they can override them. A common use is to mark which application wrote a file.

## Dropping events

Events can be dropped before they reach the outputs by matching on any field. Health checks and debug logs from busy services are the usual candidates; dropping them early saves bandwidth and storage further down the line.
`},
	},
	faq: []qa{
		{"How much CPU does {name} use?", "On a typical web server shipping a few hundred lines a second, well under 1% of one core. Parsing with regular expressions costs the most; JSON is cheapest."},
		{"Can it read logs from Docker containers?", "Yes. Point an input at `/var/lib/docker/containers/*/*.log` and set the format to `json`; each line carries the container's output in the `log` field."},
		{"What happens when the disk buffer is full?", "The oldest events are dropped to make room for new ones, and the number of dropped events is logged once a minute."},
		{"Does it support multiline messages such as stack traces?", "Yes. Give an input a `multiline` pattern that matches the first line of a message, and the lines that follow are joined to it."},
		{"Is Windows supported?", "Not yet. {name} runs on Linux and FreeBSD; macOS works for development."},
	},
	changes: []string{
		"Files deleted while being read are released right away",
		"The HTTP output honours `Retry-After`",
		"Fixed a position being saved before delivery when an output restarted",
		"Journal fields are now lower case, as documented",
		"New `positions --reset` subcommand",
		"Multiline messages are flushed after a timeout",
		"Lower CPU use for the logfmt parser",
		"Fixed reading files with lines longer than 64 KiB",
		"The disk buffer is compacted on start",
		"Added `labels` per input",
		"Fixed glob patterns with `**` on FreeBSD",
		"`validate` checks regular expressions for named groups",
		"The HTTP output can send to several URLs in turn",
		"Dropped events are counted per output",
		"`tail` prints events with colours when writing to a terminal",
		"Fixed a leak of file handles after log rotation by copytruncate",
	},
	quick: `
This page sets up {name} to ship an nginx access log to an HTTP endpoint. It assumes {name} is [installed](/docs/installation/).

## Write a configuration

Save this as ` + "`{conf}`" + `:

` + "```{ext}" + `
{example}
` + "```" + `

## Check it

` + "```shell" + `
$ {bin} validate
configuration ok: 2 inputs, 1 output
` + "```" + `

## Watch the events

` + "```shell" + `
$ {bin} tail nginx --lines 1
{"time":"{d1}T10:02:11Z","status":200,"path":"/","took_ms":3}
` + "```" + `

## Start shipping

` + "```shell" + `
$ sudo systemctl enable --now {bin}
` + "```" + `

## Next steps

- Learn how to [parse](/docs/parsing/) your own log formats.
- Add [labels](/docs/labels/) so you can tell your servers apart.
`,
	flow: [3]string{"Files & journal", "Parse & buffer", "HTTP, files, syslog"},
	demo: []string{
		"$ {bin} tail nginx --lines 2",
		`{"status":200,"path":"/","took_ms":3}`,
		`{"status":404,"path":"/x","took_ms":1}`,
	},
	daemon: "run",
	snips: map[string][2]string{
		"files": {"inputs:\n  files:\n    - /var/log/nginx/access.log\n    - /var/log/app/*.log\n  journald: true", "[inputs]\nfiles = [\"/var/log/nginx/access.log\", \"/var/log/app/*.log\"]\njournald = true"},
		"pattern": {"parse:\n  format: regex\n  pattern: '^(?P<ip>\\S+) \\S+ \\S+ \\[(?P<time>[^\\]]+)\\] \"(?P<method>\\S+) (?P<path>\\S+)'",
			"[parse]\nformat = \"regex\"\npattern = '^(?P<ip>\\S+) \\S+ \\S+ \\[(?P<time>[^\\]]+)\\] \"(?P<method>\\S+) (?P<path>\\S+)'"},
		"labels": {"labels:\n  env: production\n  role: web", "[labels]\nenv = \"production\"\nrole = \"web\""},
	},
}

var jobsKind = &kind{
	id: "jobs",
	what: []string{
		"a background job runner that uses {your existing|plain} PostgreSQL as its queue, with retries, schedules and {no extra infrastructure|nothing else to operate}",
		"a {simple|reliable} job queue built on PostgreSQL: enqueue work in the same transaction as your data and let workers pick it up",
	},
	taglines: []string{
		"Background jobs on the database you already run",
		"A job queue that is just PostgreSQL",
		"Reliable background jobs, no broker required",
		"Queue it, retry it, schedule it",
	},
	words: []string{"Queue", "Shift", "Crew", "Hive", "Mill", "Forge", "Loom", "Gear"},
	features: []feature{
		{"Transactional enqueue", "Jobs are rows. Insert one in the same transaction as your data and it exists exactly when the data does.", "database"},
		{"Retries with backoff", "Failed jobs are retried with exponential backoff and jitter, up to a limit you set per queue.", "refresh"},
		{"Schedules", "Run jobs on cron schedules. Only one worker enqueues each run, however many you start.", "clock"},
		{"Any language", "Workers call your code over a small HTTP or stdin protocol, so jobs can be written in anything.", "code"},
		{"Fast pickup", "Workers use `LISTEN/NOTIFY` to pick up new jobs within milliseconds, and `SKIP LOCKED` to never collide.", "bolt"},
		{"Metrics", "Queue depth, latency and failures as Prometheus metrics, ready for your dashboards.", "chart"},
	},
	commands: []command{
		{"migrate", "", "Create or update the tables {name} needs.", "Safe to run on every deploy: migrations that already ran are skipped, and each one runs in a transaction.", []flag{
			{"--schema", "NAME", "{bin}", "Schema to create the tables in."},
		}},
		{"worker", "", "Start a worker.", "A worker takes jobs from its queues and runs them, at most `concurrency` at a time. It finishes the jobs it is running before it exits on SIGTERM.", []flag{
			{"--queues", "LIST", "default", "Comma-separated queues to take jobs from."},
			{"--concurrency", "N", "10", "Jobs run at the same time."},
		}},
		{"enqueue", "KIND [ARGS]", "Add a job from the command line.", "`ARGS` is a JSON object. The job ID is printed on success.", []flag{
			{"--queue", "NAME", "default", "Queue to add the job to."},
			{"--at", "TIME", "", "Run the job no earlier than this time."},
			{"--priority", "N", "0", "Higher runs first."},
		}},
		{"jobs", "", "List jobs.", "", []flag{
			{"--state", "STATE", "", "`available`, `running`, `retrying`, `failed` or `done`."},
			{"--queue", "NAME", "", "Only jobs in this queue."},
			{"--limit", "N", "50", "At most N jobs."},
		}},
		{"retry", "ID...", "Run failed jobs again.", "The attempt counter starts over.", []flag{
			{"--all-failed", "", "false", "Retry every failed job."},
		}},
	},
	options: []option{
		{"database_url", "string", "", `"postgres://{bin}@localhost/app"`, "Connection string of the PostgreSQL database."},
		{"schema", "string", "{bin}", "", "Schema holding the job tables."},
		{"queues", "list", `["default"]`, `["default", "mail"]`, "Queues this worker takes jobs from."},
		{"concurrency", "int", "10", "20", "Jobs run at the same time per worker."},
		{"poll_interval", "duration", "5s", "", "How often to look for jobs when no notification arrives."},
		{"retry.max_attempts", "int", "25", "", "Attempts before a job is marked failed."},
		{"retry.backoff", "duration", "10s", "", "Wait before the first retry. Doubles with every attempt."},
		{"timeout", "duration", "30m", `"5m"`, "A job running longer is cancelled and retried."},
		{"handler.url", "string", "", `"http://127.0.0.1:8080/jobs"`, "Your application's endpoint that runs jobs."},
		{"metrics.listen", "string", "", `"127.0.0.1:9187"`, "Address for the Prometheus metrics. Empty turns them off."},
		{"log.level", "string", "info", "", "`debug`, `info`, `warn` or `error`."},
	},
	guides: []guide{
		{"queues-and-workers", "Queues and workers", "Queues", `
A job is a row in the ` + "`{bin}.jobs`" + ` table: a kind, arguments as JSON, a queue and a time it may run at. Workers take available jobs, run them and mark them done.

## Queues

Queues are names; there is nothing to create. Put jobs that must not wait behind others, such as password reset emails, into a queue of their own and give it its own workers:

` + "```shell" + `
$ {bin} worker --queues mail --concurrency 5
$ {bin} worker --queues default,reports
` + "```" + `

## How a worker takes a job

Workers select jobs with ` + "`FOR UPDATE SKIP LOCKED`" + `, so any number of them can work on the same queue without taking the same job twice. New jobs send a notification on insert, and an idle worker picks them up within milliseconds instead of waiting for the next poll.

## Running your code

A worker hands each job to your application: it POSTs the job as JSON to ` + "`handler.url`" + ` and treats a 2xx answer as success. Anything else, or no answer within ` + "`timeout`" + `, is a failure and the job is [retried](/docs/retries/).
`},
		{"retries", "Retries", "", `
Jobs fail: networks drop, services restart, bugs happen. {name} retries every failed job until it succeeds or runs out of attempts.

## Backoff

The wait before a retry doubles with every attempt, starting at ` + "`retry.backoff`" + `, and a random jitter of up to 20% keeps retries of many jobs from arriving at once. With the defaults, a job is tried for about three days before it is given up.

| Attempt | Wait before it |
| --- | --- |
| 2 | 10 seconds |
| 3 | 20 seconds |
| 5 | 1 minute 20 seconds |
| 10 | about 43 minutes |
| 25 | about 3 days in total |

## Failed jobs

A job that used up its attempts is marked ` + "`failed`" + ` and stays in the table with the error of its last attempt. List and retry them with:

` + "```shell" + `
$ {bin} jobs --state failed
$ {bin} retry --all-failed
` + "```" + `

> Make your jobs safe to run twice. A worker that crashes in the middle of a job cannot tell whether its work got done, so the job will run again.
`},
		{"schedules", "Scheduled jobs", "Schedules", `
Jobs can be enqueued on a schedule written as a cron expression. Every worker knows the schedules, but only one of them enqueues each run: they agree through an advisory lock in the database.

` + "```{ext}" + `
{schedule}
` + "```" + `

Times are in UTC. A run missed because no worker was up is not made up later; the next one is enqueued as usual.

## One-off delays

To run a single job later, give it a time when you enqueue it:

` + "```shell" + `
$ {bin} enqueue send_reminder '{"user": 42}' --at "{d1} 09:00"
` + "```" + `
`},
		{"metrics", "Metrics and monitoring", "Metrics", `
With ` + "`metrics.listen`" + ` set, every worker serves Prometheus metrics. The ones worth an alert are:

| Metric | What it tells |
| --- | --- |
| ` + "`{bin}_queue_depth`" + ` | Jobs waiting, per queue |
| ` + "`{bin}_queue_latency_seconds`" + ` | How long the oldest waiting job has waited |
| ` + "`{bin}_jobs_failed_total`" + ` | Jobs that used up their attempts |
| ` + "`{bin}_jobs_duration_seconds`" + ` | How long jobs take, as a histogram |

A latency that keeps growing means the workers cannot keep up; add workers or raise ` + "`concurrency`" + `.

## Cleaning up

Done jobs are deleted after seven days by the workers themselves, a few thousand at a time, so the table stays small without a cron job of yours.
`},
	},
	faq: []qa{
		{"Will this put too much load on my database?", "An idle worker runs one cheap query per poll interval. Under load each job costs a select, an update and a delete; PostgreSQL handles thousands of these a second on modest hardware."},
		{"Which PostgreSQL versions are supported?", "PostgreSQL 13 and newer. `SKIP LOCKED` came in 9.5, but some of the indexes need 13."},
		{"Can I enqueue jobs from my application without {name}'s CLI?", "Yes. Insert a row into `{bin}.jobs`; the columns are documented in [Queues and workers](/docs/queues-and-workers/). Client libraries only make that more convenient."},
		{"Are jobs run exactly once?", "At least once. A worker that dies mid-job leaves the job to be retried, so jobs should be safe to run twice."},
		{"Can I pause a queue?", "Stop the workers of that queue. Jobs keep being accepted and wait until workers are back."},
	},
	changes: []string{
		"Workers finish running jobs before exiting on SIGTERM",
		"New `{bin}_queue_latency_seconds` metric",
		"Fixed a schedule being enqueued twice around a daylight saving change",
		"`jobs` can filter by kind",
		"Done jobs are cleaned up in smaller batches",
		"The HTTP handler receives the attempt number in a header",
		"Fixed `retry --all-failed` skipping jobs with a priority",
		"Lower latency when many workers wait on the same queue",
		"`migrate` takes an advisory lock, so parallel deploys are safe",
		"Jobs can be enqueued with a uniqueness key",
		"Added `--at` to `enqueue`",
		"Fixed a timeout not cancelling the request to the handler",
		"Connection pool size follows `concurrency`",
		"The error of every attempt is kept, not only the last",
		"Schedules accept `@hourly` and `@daily`",
	},
	quick: `
This page sets up {name} with a local PostgreSQL and runs a first job. It assumes {name} is [installed](/docs/installation/).

## Create the tables

` + "```shell" + `
$ export {env}DATABASE_URL=postgres://localhost/app
$ {bin} migrate
applied 6 migrations to schema {bin}
` + "```" + `

## Start a worker

Your application needs an endpoint that runs jobs. For a first try, any server that answers 200 will do:

` + "```shell" + `
$ {bin} worker --concurrency 2
worker started: queues=default concurrency=2
` + "```" + `

## Enqueue a job

In another terminal:

` + "```shell" + `
$ {bin} enqueue send_welcome '{"user": 42}'
enqueued job 1
` + "```" + `

The worker picks the job up right away and logs the result.

## Next steps

- Read how [queues and workers](/docs/queues-and-workers/) fit together.
- Set up [schedules](/docs/schedules/) for recurring work.
`,
	flow: [3]string{"Your app", "PostgreSQL queue", "Workers"},
	demo: []string{
		"$ {bin} enqueue send_welcome '{\"user\": 42}'",
		"enqueued job 1842",
		"$ {bin} jobs --state done --limit 1",
		"1842  send_welcome  done  0.21s",
	},
	daemon: "worker",
	snips: map[string][2]string{
		"schedule": {"schedules:\n  - kind: send_digest\n    cron: \"0 7 * * *\"\n  - kind: cleanup_sessions\n    cron: \"*/15 * * * *\"\n    queue: maintenance",
			"[[schedules]]\nkind = \"send_digest\"\ncron = \"0 7 * * *\"\n\n[[schedules]]\nkind = \"cleanup_sessions\"\ncron = \"*/15 * * * *\"\nqueue = \"maintenance\""},
	},
}

var monitorKind = &kind{
	id: "monitor",
	what: []string{
		"a {small|self-hosted} uptime monitor that checks your websites, APIs and TLS certificates and tells you when something breaks",
		"a {single-binary|lightweight} monitor for HTTP endpoints, TCP ports and certificate expiry, with alerts and a status page",
	},
	taglines: []string{
		"Know before your users do",
		"Uptime checks you run yourself",
		"Watch your sites. Sleep at night.",
		"Simple monitoring for websites and certificates",
	},
	words: []string{"Watch", "Beacon", "Sentry", "Lookout", "Signal", "Pulse", "Lantern", "Heron"},
	features: []feature{
		{"HTTP and TCP checks", "Check status codes, response times and body content, or simply whether a port answers.", "globe"},
		{"Certificate expiry", "Every HTTPS check also looks at the certificate and warns weeks before it runs out.", "lock"},
		{"Alerts that matter", "Webhooks, email and chat, sent only after a failure is confirmed from a second attempt.", "bell"},
		{"Status page", "A plain, fast status page with the history of every check, served by {name} itself.", "chart"},
		{"Config as code", "All checks live in one file you can keep in git. No database to back up.", "code"},
		{"Light on resources", "Hundreds of checks on the smallest VM, with state kept in a single file.", "bolt"},
	},
	commands: []command{
		{"run", "", "Run all checks on their schedule.", "Also serves the status page when `http.listen` is set.", nil},
		{"check", "URL", "Run one check now and print the result.", "Handy for trying out a check before adding it to the configuration.", []flag{
			{"--expect", "CODE", "200", "Status code that counts as up."},
			{"--contains", "TEXT", "", "Text the body must contain."},
			{"--timeout", "DURATION", "10s", "Give up after this long."},
		}},
		{"status", "", "Show the current state of every check.", "", []flag{
			{"--json", "", "false", "Print the state as JSON."},
		}},
		{"test-alert", "CHANNEL", "Send a test alert to a channel.", "Use it to check webhook URLs and mail settings.", nil},
		{"certs", "", "List certificates and when they expire.", "", []flag{
			{"--within", "DAYS", "", "Only certificates expiring within this many days."},
		}},
	},
	options: []option{
		{"interval", "duration", "1m", "", "Time between two runs of a check."},
		{"timeout", "duration", "10s", "", "A check that takes longer has failed."},
		{"retries", "int", "1", "", "Extra attempts before a failure is reported."},
		{"checks", "list", "[]", "", "The checks; see [Checks](/docs/checks/)."},
		{"alerts.webhook", "string", "", `"https://hooks.example.com/ops"`, "URL to POST alerts to as JSON."},
		{"alerts.email", "string", "", `"ops@example.com"`, "Address to mail alerts to."},
		{"alerts.smtp", "string", "", `"smtp://mail.example.com:587"`, "Mail server for email alerts."},
		{"tls.warn_days", "int", "21", "14", "Warn when a certificate expires within this many days."},
		{"state_file", "string", "/var/lib/{bin}/state.db", "", "Where results and history are kept."},
		{"http.listen", "string", "", `"127.0.0.1:8090"`, "Address of the status page. Empty turns it off."},
		{"log.level", "string", "info", "", "`debug`, `info`, `warn` or `error`."},
	},
	guides: []guide{
		{"checks", "Checks", "", `
A check is something {name} looks at regularly and judges as up or down. Each one has a name, a target and what counts as success.

` + "```{ext}" + `
{checks}
` + "```" + `

## HTTP checks

An HTTP check is up when the server answers with the expected status code, within the timeout, and, if ` + "`contains`" + ` is set, with that text in the body. Redirects are followed up to five times.

## TCP checks

A TCP check is up when a connection to the port can be opened. Nothing is sent.

## Retries

A check that fails is tried again right away, ` + "`retries`" + ` times, before it counts as down. One lost packet should not wake anybody up.
`},
		{"alerts", "Alerts", "", `
{name} sends an alert when a check goes down and another when it comes back up. Nothing is sent while the state stays the same.

## Webhooks

A webhook alert is a POST with a JSON body:

` + "```json" + `
{"check": "website", "state": "down", "since": "{d1}T03:12:40Z", "error": "status 502"}
` + "```" + `

Most chat services accept incoming webhooks, often with a small adapter. Test the URL with ` + "`{bin} test-alert webhook`" + `.

## Email

Set ` + "`alerts.email`" + ` and ` + "`alerts.smtp`" + `. The SMTP URL may contain a user and password; ` + "`smtps://`" + ` connects with TLS from the start, ` + "`smtp://`" + ` upgrades with STARTTLS.

> Alerts that cannot be delivered are retried for an hour and then logged.
`},
		{"certificates", "Certificate expiry", "Certificates", `
Every HTTPS check also records the certificate the server presented: who issued it, which names it covers and when it expires.

` + "```shell" + `
$ {bin} certs --within 30
website     {domain}       expires in 12 days   R11
api         api.example.org     expires in 27 days   E6
` + "```" + `

A warning goes out once when a certificate enters the ` + "`tls.warn_days`" + ` window, and again a week before it expires. Certificates that renew automatically never come close, so a warning usually means the renewal is failing.
`},
		{"status-page", "The status page", "Status page", `
With ` + "`http.listen`" + ` set, {name} serves a status page: each check with its current state and a bar of the last 90 days.

The page is plain HTML with no scripts and no external requests, so it loads anywhere. Put it behind your reverse proxy to give it a name and a certificate, or keep it on localhost and reach it over SSH.

## JSON

The same information is available at ` + "`/status.json`" + ` for scripts and other monitors.
`},
	},
	faq: []qa{
		{"How many checks can one instance run?", "Thousands. A check costs one request per interval; the default of a minute makes a thousand checks about seventeen requests a second."},
		{"Can I check from several locations?", "Run {name} in each location. Every instance checks on its own and alerts on its own."},
		{"Does it keep history?", "Yes, ninety days of results per check, in the state file. Older results are summarised per day."},
		{"Can I silence alerts during maintenance?", "Yes: `{bin} status --pause website 2h` pauses a check for two hours."},
		{"Does the status page need a database?", "No. Everything lives in the state file next to the configuration."},
	},
	changes: []string{
		"Certificate warnings include the issuer",
		"TCP checks report the time to connect",
		"Fixed an alert being sent twice after a restart",
		"The status page shows 90 days instead of 30",
		"New `certs --within` flag",
		"Webhook alerts carry the error of the last attempt",
		"Fixed checks of URLs with non-ASCII host names",
		"`check` prints the response time",
		"Email alerts work with servers that require STARTTLS",
		"The state file is compacted daily",
		"Added `contains` for HTTP checks",
		"Fixed redirects losing the query string",
		"Checks spread their first run over the interval",
		"New `/status.json` endpoint",
	},
	quick: `
This page sets up {name} to watch one website and alert a webhook. It assumes {name} is [installed](/docs/installation/).

## Try a check

` + "```shell" + `
$ {bin} check https://{domain}/
up: 200 OK in 84 ms, certificate valid for 61 days
` + "```" + `

## Write a configuration

Save this as ` + "`{conf}`" + `:

` + "```{ext}" + `
{example}
` + "```" + `

## Start it

` + "```shell" + `
$ sudo systemctl enable --now {bin}
$ {bin} status
website   up   84 ms   checked 12s ago
` + "```" + `

## Next steps

- Add more [checks](/docs/checks/), including TCP ports.
- Turn on the [status page](/docs/status-page/).
`,
	flow: [3]string{"Your sites", "Checks", "Alerts & status page"},
	demo: []string{
		"$ {bin} check https://example.org/",
		"up: 200 OK in 84 ms",
		"certificate valid for 61 days",
	},
	daemon: "run",
	snips: map[string][2]string{
		"checks": {"checks:\n  - name: website\n    url: https://{domain}/\n    contains: Welcome\n  - name: api\n    url: https://api.example.org/health\n    expect: 204\n  - name: database\n    tcp: db.internal:5432",
			"[[checks]]\nname = \"website\"\nurl = \"https://{domain}/\"\ncontains = \"Welcome\"\n\n[[checks]]\nname = \"api\"\nurl = \"https://api.example.org/health\"\nexpect = 204\n\n[[checks]]\nname = \"database\"\ntcp = \"db.internal:5432\""},
	},
	example: []string{"checks"},
}
