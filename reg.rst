=====
 reg
=====

-------------------
 Process REGulator
-------------------

:Date: April 2013
:Manual section: 1

SYNOPSIS
========

::

   reg [OPTION]...

DESCRIPTION
===========

The ``reg``-ulator forces the execution of another process to stay
within a resource budget decided dynamically.

Operation
---------

``reg`` can be either attached to an already running process (``reg
-a <pid>``), individual thread (``reg -a thread:<pid>``) or start the process
to be regulated (``reg <cmd...>``).

Once a process is harnessed, ``reg`` reads *resource supplies* on its
standard input, and produces *resource usage reports* on its standard
output. It also regulates the harnessed processes so that their
resource consumption does not exceed the provided supplies.


Options
-------

``-a``, ``--attach <PID>``
    Attach to an already running process or thread. If ``<PID>`` is a number,
    attach to an entire process, including all its threads. If ``<PID>`` is
    of the form ``thread:PID``, attach to the single thread with that PID.

``-d``, ``--domain <LABEL>``
    Define a management domain with label ``<LABEL>``. See `MANAGEMENT
    DOMAINS`_ below.

``-f``, ``--follow <PRED>``
    Use command ``<PRED>`` to determine whether to follow child
    processes/threads. See `RECURSIVE THREAD/PROCESS CREATION`_ below.

``-g``, ``--granularity <N>``
    Set the granularity of measurements/decisions to ``<N>``. See
    `GRANULARITY`_ below.

``-i``, ``--input <FILE>``
    Use ``<FILE>`` as input stream. By default, the standard input is
    used. See `INPUT LANGUAGE`_ below.

``-o``, ``--output <FILE>``
    Use ``<FILE>`` as output stream. By default, the standard output is used.
    See `OUTPUT FORMAT`_ below.

``-p``, ``--protocol <SPEC>``
    Use ``<SPEC>`` as protocol to regulate the processes. See
    `REGULATION PROTOCOL`_ below.

``-r``, ``--resource <LABEL>:<SPEC>``
    Add the resource specified by ``<SPEC>`` under control of ``reg``
    with label ``<LABEL>``. See `RESOURCE SPECIFICATION`_ below.

``-R``, ``--rate <N>``
    Configure how often status records are generated on the output stream. See
    `OUTPUT RATE`_ below.

``-s``, ``--steps <SPEC>``
    Use ``<SPEC>`` as a progress indicator function. See `PROGRESS
    INDICATOR`_ below.

``-t``, ``--ticks <SPEC>``
    Use ``<SPEC>`` as a time discretization function. See `TIME
    DISCRETIZATION`_ below.

``-v``, ``--verbose``
    Verbose execution (details on standard error).

RESOURCE MANAGEMENT
===================

Resource management is modeled as follows:

- process execution is discretized over a function *t*, typically
  time, which increases monotonically regardless of process
  state. Other options are available, see `TIME DISCRETIZATION`_
  below.  The unit for *t* is *ticks*.

- *progress* during process execution can be measured using a function
  *d(t)*, typically "user time", which increases monotonically while a
  process is running and not stopped. Other options are available, see
  `PROGRESS INDICATOR`_ below. The unit for *d(t)* values is
  *steps*.

- for a given resource *f* the harnessed processes have a *current
  level* of resource usage *f(t)*, which can be observed at every
  tick, and for which there may be no known upper limit. The unit for
  values of *f(t)* is *stuff*.  (For example bytes for memory,
  bytes/sec for channels, etc.)

Using this model, ``reg`` controls resource usage as follows:

- ``reg`` associates to each resource *f* a *supply*, expressed in
  *stuff.steps* (amount of stuff, times amount of steps).

- ``reg`` operates iteratively; at each iteration the increase of
  steps *d(t)* is observed.  Upon each increment of *d* from *dp =
  d(tp)* (previous step counter) to *dn = d(t)* (new step counter):

  1. ``reg`` measures the *integrated resource consumption* since the
     last step, computed by *f(t)* times *(dn - dp)*.

  2. the integrated resource consumption is substracted from the supply;

  3. if the supply becomes zero or negative, the process is *stopped* until
     the supply is increased again.

For example, if the step is user time in seconds and the resource is
current power usage (watts), the supply is expressed in watts.seconds
(energy). With a supply of 1 watt.second, a process that consumes .5
watts per second will be stopped after 2 seconds, and a process that
consumes 2 watts per second will be stopped after .5 seconds.

Another example: if the resource is current memory footprint (bytes),
the supply is expressed in bytes.seconds. With a supply of
100MBytes.second, a process that allocates 10MBytes every second will
stop after 10 seconds, whereas one that allocates 1GByte in one go
will be stopped directly after this first allocation with a remaining
supply of -900 MBytes.second.

Note: regulation only occurs when *t* increases.

TIME DISCRETIZATION
===================

By default, ``reg`` discretizes over *absolute time*, where the tick
unit for *t* is *seconds*.

The following alternatives are available:

======================= ===================================== =================
Option                  Description                           Unit
======================= ===================================== =================
``-t <n>.realseconds``  Absolute time                         seconds . ``<n>``
``-t <n>.spentjoules``  Energy spent                          joules . ``<n>``
``-t <n>.controlled``   Explicit messages on ``reg``'s input  t-ticks . ``<n>``
``-t <n>.[METHOD]``     Custom (cf below)                     (cf below)
======================= ===================================== =================

The first numeric argument is optional and specifies a multiplier. For
example, ``-t 3600 realseconds`` uses hours as tick unit, and ``-t
1e-12 spentjoules`` uses picojoules as tick unit. If it is not
specified, it defaults to 1. For example, ``-t 1 realseconds`` is
equivalent to ``-t realseconds``.

The multiplier can also be specified using a SI multiplier: ``k`` for
1000, ``m`` for 0.001, etc. For example ``-t p spentjoules`` is
equivalent to ``-t 1e-12 spentjoules``.

Custom discretizations can be defined using the following options:

``-t <n>.re:<path>:<regex>``
  Use a regular expression match on the specified file and use the
  first match group (if any) as tick counter.

``-t <n>.cg:<subsystem>:<regex>``
  Use a regular expression match on the specified control file of the
  selected cgroup subsystem and use the first match group (if any) as
  tick counter.

When using custom time discretizations, beware to use a function that
increases even when the harnessed process is stopped. Otherwise,
deadlock would ensue: ``reg`` would stop regulating and never wake up
the harnessed process again.

PROGRESS INDICATOR
==================

By default, ``reg`` measures process progress using *user time*, where
the step unit for *d(t)* is *seconds*. The following alternatives are
available:

======================== ========================= =================
Option                   Description               Unit
======================== ========================= =================
``-s <n>.userseconds``   User time                 seconds . ``<n>``
``-s <n>.jiffies``       Scheduler time slices     jiffies . ``<n>``
``-s <n>.instructions``  Instructions executed     instructions . ``<n>``
``-s <n>.[METHOD]``      Custom                    (depends on method)
======================== ========================= =================

Custom progress functions can be configured with ``-s`` as for ``-t`` above.


RESOURCE SPECIFICATION
======================

A resource function and supply bin can be defined with the option
``-r <LABEL>:<FUNCTION>``. ``-r`` can be used multiple times with
different labels to define multiple supply bins.

The following functions are available:

=============== ============================== ===================
Function        Description                    Unit
=============== ============================== ===================
``steps``       Current step counter           (same as step unit)
``threads``     Number of threads harnessed    threads
``load``        Average CPU load               load
``vsize``       Virtual memory size            bytes
``rsize``       Resident memory size           bytes
``[METHOD]``    Custom                         (depends on method)
=============== ============================== ===================

All special progress functions (``userseconds``, ``jiffies``, etc) are
also valid resource functions.

Custom resource functions can be computed with ``-r`` as for ``-s``
and ``-t`` above.


INPUT LANGUAGE
==============

``reg`` accepts the following newline-terminated commands on its
input stream:

``. <ticks>``
  If using ``-s controlled`` (see `TIME DISCRETIZATION`_ above),
  increment the discretization counter by the specified amount of
  ticks. Otherwise, do nothing.

``+ <supply> <amount>``
  Add the specified number of stuff.steps in the selected resource
  supply(ies). If ``<amount>`` is ``*``, add an infinite supply.

``- <supply> <amount>``
  Substract the specified number of stuff.steps from the selected
  resource supply(ies). If ``<amount>`` is ``*``, empty the entire
  supply. If the bin does not exist or its supply is already empty, the
  command has no effect.

``?``
  Emit a status record on the output stream.

The syntax of ``<supply>`` for the commands ``+`` and ``-`` can be a
shell wildcard pattern, using the syntax recognized by fnmatch(1). If
a pattern matches multiple resource labels, the operation (add or
substract) is performed on all of them.

All amounts (or ticks for ``.``) can be followed by an SI
multiplier. For example, ``. 1k`` is equivalent to ``. 1000``.


OUTPUT FORMAT
=============

Each status record ends with a newline
character, and is composed of the following space-separated columns:

- the tag from command ``?``, or ``-`` if the record is produced automatically
  from ``-R`` (cf `OUTPUT RATE`_ below)
- the label of the management domain (cf. `MANAGEMENT DOMAINS`_ below),
- the current tick,
- the tick delta (number of ticks elapsed since the last status record),
- the current step, and step delta,
- the number of resource functions defined,
- for each resource function defined:

  - the label of the function,
  - the current supply,
  - the amount of supply added/substracted on the input stream since the last status record,
  - the amount of supply substracted by the process execution since the last status record,

- the number of threads harnessed,
- for each thread harnessed:

  - the process ID of the process where the thread belongs (TGID),
  - the process ID of the thread itself.

OUTPUT RATE
===========

By default, ``reg`` produces status records after each explicit ``?``
command on the input stream.

Additionally, the option ``-R <N>.steps`` and ``-R <N>.ticks``
instructs ``reg`` to emit records periodically, with the period
specified (either steps or ticks). The number can be followed by an SI
multiplier.

``reg`` does not block on output: if the output stream is blocked, the
deltas accumulate until ``reg`` becomes able to output records again. If
more than one ``?`` input commands are received on the input, or periods
of ``-R`` are elapsed while the output stream is blocked, they are
ignored and only one status record is emitted on the output stream
when it becomes unblocked.

With option ``-R 0`` (flood), as many status records are generated as
possible when the output stream is unblocked. The consumer process is
then in charge of controlling the rate by throttling its input.

With ``-R none`` the automatic output is disabled and records are only
output when ``?`` is received on the input.  (this is the default).


GRANULARITY
===========

The rate at which ``reg`` monitors ``t`` and makes regulation decisions
is determined by the *granularity* parameter, selected with option
``-g <value>``.

The granularity is the multiple of the unit of the time discretization
function that ``reg`` attempts to track. For example, with time
measured in seconds and ``-g 0.001``, ``reg`` will attempt to keep
track of resource usage every millisecond.

By default, the granularity is 1.


RECURSIVE THREAD/PROCESS CREATION
==================================

By default, all threads and processes recursively created by
the regulated program are collectively regulated by the same
``reg`` instance.

If the option ``-f <pred>`` is specified, ``reg`` will run the command
``<pred>`` upon the creation of each new thread or process to decide
whether to keep the child thread/process regulated.

If the ``<pred>`` exits with status 0, the created thread/process
stays regulated. If ``<pred>`` exits with a non-zero status, the
created thread/process is removed from ``reg``'s control. Three
command line arguments are provided to ``<pred>``:

- the parent ID (PPID),
- the thread group leader ID of the newly created thread (TGID)
- the process ID of the newly created thread (PID).

(If TGID = PID, a new process was created. Otherwise, a new thread was
created in the process identified by the TGID.)

The default behavior is thus equivalent to ``-f true``.

REGULATION PROTOCOL
===================

By default, ``reg`` uses Linux cgroups' "freeze" subsystem to regulate
processes: the processes are frozen if a resource supply is exhausted,
and thawed when the supply becomes available again.

The protocol can be specified as follows:

================== ================================================
Option             Description
================== ================================================
``-p freeze``      Use cgroups/freeze as regulation mechanism (default).
``-p stop``        Use SIGSTOP/SIGCONT as regulation mechanism.
``-p out:<FILE>``  Send commands through ``<FILE>``.
``-p run:<CMD>``   Use the external program ``<CMD>``.
================== ================================================

With ``-p fd``, the following commands are sent to the specified file:

``overflow <RES> <SUPPLY> <DELTA> <DOM> <PIDs...>``

    Signal an overflow. The fields are as follows:

    ============== =================================
    Field          Description
    ============== =================================
    ``<RES>``      Resource label causing the overflow, as configured by ``-r``.
    ``<SUPPLY>``   Current supply for the resource.
    ``<DELTA>``    Last amount substracted by the process.
    ``<DOM>``      cf. `MANAGEMENT DOMAINS`_ below.
    ``<PIDs...>``  Current list of harnessed processes.
    ============== =================================

``ok <PIDs...>``
    Signal that all supplies are zero or positive.

With ``-p run``, the specified command is invoked as follows:

``<CMD> overflow <RES> <SUPPLY> <DELTA> <DOM> <PIDs...>``

or

``<CMD> ok <PIDs...>``

(same argument meanings as ``-p fd`` above)

Note: the effect of an overflow command should be to stop the progress
function *d(t)* (make it constant), so that its integrated resource
consumption stays zero until the supply is increased and the process
is restarted.

MANAGEMENT DOMAINS
==================

In the current implementation, a given thread can be harnessed by at
most one ``reg`` instance. Therefore, each ``reg`` instance can
monitor multiple time discretization, progress and resource usage
functions simultaneously.

This is supported as follows:

- ``reg`` defines one or more *management domains*; the first is
  always defined and is named ``default``. More domains are declared
  with option ``-d``.

- each management domain must define:

  - exactly one time discretization function,
  - exactly one progress function,
  - one or more resource functions,
  - an input and output stream.

- the parameters ``-t``, ``-s``, ``-g``, ``-R``, ``-r``, ``-i`` and ``-o``
  described above set the corresponding parameter of the domain
  ``default``. If either ``-t`` or ``-s`` are not used, ``default``
  uses real time and user time, respectively. If either ``-i`` or
  ``-o`` are not used, ``default`` uses the standard input and output,
  respectively.

- to set parameters in a domain ``DOM``, the options ``-t DOM=<arg>``,
  ``-s DOM=<arg>``, ``-g DOM=<arg>``, ``-R DOM=<arg>``, ``-r
  DOM=<arg>``, ``-i DOM=<arg>``, ``-o DOM=<arg>`` can be
  used.


EXIT STATUS
===========

``reg`` terminates with the following exit codes:

0
   All harnessed process/thread have terminated, or both the input and
   output streams have been closed.

1
   A configuration or environment error prevents ``reg`` from starting.

2
   An invalid command was received on the input stream.

Other errors (signals, unknown situations etc) are reported with other
exit codes.
