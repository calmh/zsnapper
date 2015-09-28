:Title: ZSNAPPER(1)
:Author: Jakob Borg
:Date: September 2015

NAME
====

zsnapper - automatically manage ZFS snapshots

SYNOSIS
=======

| zsnapper [OPTIONS]

OPTIONS
=======

-c <path>
	Path to configuration file (default **/opt/local/zsnapper/etc/zsnapper.yml**)

-v
    Enable verbose output

DESCRIPTION
===========

Zsnapper automatically creates ZFS snapshots on a specified schedule, while
also removing old snapshots as required. The behavior is specified by the
**CONFIGURATION**. Snapshots are grouped into *families*, each having a name,
containing a set of *datasets*, a *schedule*, a number of snapshots to *keep*
and an option for *recursiveness*.

Snapshots are named according to a "family-timestamp" convention. For example,
a snapshot belonging to the "test" family may be called **test-
20150928T112300Z**. Time stamps are always in the UTC timezone and formatted
in the ISO 8601 basic format.

CONFIGURATION
=============

The following is an example of a valid snapper configuration file::

    - family: quick
      datasets:
        - zones/var
        - data/delegated/*
      schedule: "0 */5 * * * *"
      keep: 12
      recursive: false

    - family: hourly
      datasets:
        - zones/$(vmadm list -Ho uuid)*
      schedule: "@hourly"
      keep: 12
      recursive: true

Two families are declared, **quick** and **hourly**. Each has a list of
datasets, containing wildcards and shell expansions - see the **Dataset
Expansions** section for the meaning and syntax of these. The **schedule** is
a cron format string with a seconds field - see the **Schedule Format**
section for the meaning of the individual fields.

Dataset Expansions
------------------

The list of datasets configured per family is interpreted as a *filter* on the
list of actually existing datasets at the time of job execution. This means
that it's valid and not an error to mention a dataset that does not exist. Two
forms of expansions are interpreted: wildcards, and shell invocations.
Wildcards follow the following syntax, which is essentially standard shell
glob patterns::

  pattern:
    { term }

  term:
    '*'         matches any sequence of non-/ characters
    '?'         matches any single non-/ character
    '[' [ '^' ] { character-range } ']'
                character class (must be non-empty)
    c           matches character c (c != '*', '?', '\\', '[')
    '\\' c      matches character c

  character-range:
    c           matches character c (c != '\\', '-', ']')
    '\\' c      matches character c
    lo '-' hi   matches character c for lo <= c <= hi

Shell invocations take the form of ``$(command)``, where the command is
executed using ``/bin/sh -c``. Each line of output from the command results in
one pattern being added to the dataset list. If multiple shell invocations are
present in the same pattern, they are expanded and all combinations are added.
As an example, consider the pattern::

  a-$(echo f1; echo f2)-b-$(echo f3; echo f4;)-*

The shell invocations are expanded, resulting the following list of patterns::

  a-f1-b-f3-*
  a-f1-b-f4-*
  a-f2-b-f3-*
  a-f2-b-f4-*

This list of patterns is then compared against the list of existing datasets
to form the list of datasets to snapshot. For a concrete use case, consider
the pattern:

  zones/$(vmadm list -Ho uuid)*

On SmartOS, the shell invocation expands into a list of virtual machine IDs.
The result is a list of dataset patterns that match virtual machine zones, and
virtual machine disk volumes, but not template images as these are not
returned by the ``vmadm list`` command.

Schedule Format
---------------

> This section is copied from https://github.com/robfig/cron/blob/master/doc.go, formatted for this man page.

A cron expression represents a set of times, using 6 space-separated fields::

  Field name   | Allowed values  | Allowed special characters
  ----------   | --------------  | --------------------------
  Seconds      | 0-59            | * / , -
  Minutes      | 0-59            | * / , -
  Hours        | 0-23            | * / , -
  Day of month | 1-31            | * / , - ?
  Month        | 1-12 or JAN-DEC | * / , -
  Day of week  | 0-6 or SUN-SAT  | * / , - ?

Note: Month and Day-of-week field values are case insensitive.  "SUN", "Sun",
and "sun" are equally accepted.

Special Characters
~~~~~~~~~~~~~~~~~~

Asterisk (``*``)
  The asterisk indicates that the cron expression will match for all values of the
  field; e.g., using an asterisk in the 5th field (month) would indicate every
  month.

Slash (``/``)
  Slashes are used to describe increments of ranges. For example 3-59/15 in the
  1st field (minutes) would indicate the 3rd minute of the hour and every 15
  minutes thereafter. The form "*\/..." is equivalent to the form "first-last/...",
  that is, an increment over the largest possible range of the field.  The form
  "N/..." is accepted as meaning "N-MAX/...", that is, starting at N, use the
  increment until the end of that specific range.  It does not wrap around.

Comma (``,``)
  Commas are used to separate items of a list. For example, using "MON,WED,FRI" in
  the 5th field (day of week) would mean Mondays, Wednesdays and Fridays.

Hyphen (``-``)
  Hyphens are used to define ranges. For example, 9-17 would indicate every
  hour between 9am and 5pm inclusive.

Question mark (``?``)
  Question mark may be used instead of '*' for leaving either day-of-month or
  day-of-week blank.

Predefined schedules
~~~~~~~~~~~~~~~~~~~~

You may use one of several pre-defined schedules in place of a cron expression:

@yearly (or @annually)
  Run once a year, midnight, Jan. 1st (``0 0 0 1 1 *``)
@monthly
  Run once a month, midnight, first of month (``0 0 0 1 * *``)
@weekly
  Run once a week, midnight on Sunday (``0 0 0 * * 0``)
@daily (or @midnight)
  Run once a day, midnight (``0 0 0 * * *``)
@hourly
  Run once an hour, beginning of hour (``0 0 * * * *``)

Intervals
~~~~~~~~~

You may also schedule a job to execute at fixed intervals.  This is supported by
formatting the cron spec like this::

    @every <duration>

where "duration" is a string accepted by time.ParseDuration
(http://golang.org/pkg/time/#ParseDuration).

For example, ``"@every 1h30m10s"`` would indicate a schedule that activates every
1 hour, 30 minutes, 10 seconds.

Note: The interval does not take the job runtime into account.  For example,
if a job takes 3 minutes to run, and it is scheduled to run every 5 minutes,
it will have only 2 minutes of idle time between each run.

Time zones
~~~~~~~~~~

All interpretation and scheduling is done in the machine's local time zone (as
provided by the Go time package (http://www.golang.org/pkg/time).

Be aware that jobs scheduled during daylight-savings leap-ahead transitions will
not be run!
