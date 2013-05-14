=====
 reg
=====

------------------
 System REGulator
------------------

:Date: May 2013
:Manual section: 1

SYNOPSIS
========

::

   reg [OPTION]...

DESCRIPTION
===========

The ``reg``-ulator monitors a resource metric of a system and signals
the system when the metric exceeds a resource budget decided
dynamically.


Operation
---------

``reg`` reads *resource supplies* on its input stream, and produces
*resource usage reports* on its output stream. Asynchronously, it
monitors *resource usage metrics* on the system. Whenever the resource
supply is exhausted, it informs an *actuator* to regulate the system.

See `RESOURCE MANAGEMENT`_ below.


Options
-------

``-i``, ``--input <FILE>``
    Use ``<FILE>`` as input stream. By default, the standard input is
    used. See `INPUT LANGUAGE`_ below.

``-o``, ``--output <FILE>``
    Use ``<FILE>`` as output stream. By default, the standard output is used.
    See `OUTPUT FORMAT`_ below.

``-t``, ``--ticks <SPEC>``
    Use ``<SPEC>`` as a time discretization function. See `TIME
    DISCRETIZATION`_ below.

``-g``, ``--granularity <N>``
    Force the granularity of the tick generator to ``<N>``. See
    `GRANULARITY`_ below.

``-s``, ``--steps <SPEC>``
    Use ``<SPEC>`` as a progress indicator function. See `PROGRESS
    INDICATOR`_ below.

``-m``, ``--monitor <SPEC>``
    Use ``<SPEC>`` as monitor for resource usage.
    See `MONITOR SPECIFICATION`_ below.

``-a``, ``--actuator <SPEC>``
    Use ``<SPEC>`` as actuator to regulate the system. See
    `ACTUATION`_ below.

``-p``, ``--periodic-output <SPEC>``
    Configure how often status records are generated on the output stream. See
    `OUTPUT RATE`_ below.

``-h``, ``--help``
    Display usage information and exit.

``--version``
    Output version information and exit.

RESOURCE MANAGEMENT
===================

Resource management is modeled as follows:

- process execution is discretized over a function *t*, typically
  time, expected to increase monotonically. Other options are
  available, see `TIME DISCRETIZATION`_ below.  The unit for *t* is
  *ticks*.

- *progress* is measured using a function *p(t)*, for example "user
  time", which increases monotonically while the monitoring system is
  active. Other options are available, see `PROGRESS INDICATOR`_
  below. The unit for *p(t)* values is *steps*.

- the *current level* of resource usage *f(t)* is observed at every
  tick. The unit for values of *f(t)* is *stuff*.  (For example bytes
  for memory, bytes/sec for channels, etc.)

Using this model, ``reg`` regulates resource usage as follows:

- ``reg`` maintains a *supply*, expressed in
  *stuff.steps* (amount of stuff, times amount of steps).

- at each tick event, the increase of steps *p(t)* is observed.  Upon
  each increment of *p* from *p = p(tp)* (previous step counter) to
  *p' = p(t)* (new step counter):

  1. ``reg`` measures the *integral resource consumption* since the
     *last step*, computed by *f(t)* times *(p - p')*.

  2. the integral resource consumption is substracted from the supply;

  3. if the supply becomes zero or negative, the actuator is triggered
     at each subsequent *t* event, until the supply is
     increased again.

It is expected that the actuator alters the system to slow down the
progress function, so that the integral resource consumption
is reduced.

Consider for example an instance of ``reg`` whose actuators *stops* a
process upon supply exhaustion, and *restarts* the process when the
supply becomes positive again. The step function is the user time used
by the process in seconds.

If the resource is current power usage (watts), the supply is
expressed in joules (watts.seconds = energy). With a supply of 1 joule, a
process that consumes .5 joules per second will be stopped after 2
seconds, and a process that consumes 2 joules per second will be
stopped after .5 seconds.

If the resource is current channel throughput (bytes/second), the
supply is expressed in bytes (bytes/second . seconds). With a supply
of 100MBytes, a process that consumes 10MB/s every second will stop
after 10 seconds, whereas one that uses 1GB/s will be stopped
immediately after the first tick with a remaining supply of -900
MBytes.

Note: regulation only occurs when *t* increases.

TIME DISCRETIZATION
===================

By default, ``reg`` discretizes over *absolute time*, where the tick
unit for *t* is *seconds*.

The argument ``-t`` specifies a function that generates tick
events. The general syntax is ``-t <TYPE>/<FLAGS>:<ARG>``, where
``<TYPE>`` indicates the type of function, ``<FLAGS>`` indicate how
the function's events are modified, and ``<ARG>`` is an argument to
the function. See `Tick functions`_ and `Tick function flags`_ below.

When using custom time discretizations, beware to use a function that
increases monotonically, even when the monitored system is
idle. Otherwise, deadlock would ensue: ``reg`` would stop regulating
and never trigger the actuator again.

Tick functions
--------------

The following predefined functions are available:

======================= =====================================
Function                Description
======================= =====================================
``time``                Absolute time (in seconds)
``ptime``               Absolute time (in periods)
``cmd``                 Shell command (individual calls)
``proc``                Shell command (interactive)
``instant``             Generates a single event
======================= =====================================


With both functions ``time`` and ``ptime``, the argument
specifies the time period between tick events. The difference between
``time`` and ``ptime`` is that the value of each ``time`` event
reports the number of *seconds* actually elapsed since the last event,
whereas the value of each ``ptime`` event reports the number of
*periods* elapsed.

For example, ``-t time:2s`` would generate events 2, 4, 6, 8... at
2-second intervals, whereas ``-t ptime:2s`` would generate events 1, 2,
3, 4..., also at 2-second intervals.

With function ``cmd``, the command given as argument is run
repeatedly. A tick event is generated every time the command
terminates, using the value reported on its standard output.

With the function ``proc``, the command given as argument is run in
the background. A tick event is generated every time the command
outputs a line of text on its standard output.

With the function ``instant``, a single tick event is generated, whose
value is determined by the argument to the function (default 0). This
feature was originally implemented for debugging ``reg``.

Tick function flags
-------------------

The optional ``<FLAGS>`` indicate how the function's values are
translated to tick events.

``z`` (force origin zero)
   Force the sequence of tick events to have origin value 0, even if
   the underlying function has a different origin.

``d`` (deltas, applies to ``cmd`` and ``proc``)
   Each output from the command reports the additional
   number of ticks elapsed since the last output.

``o`` (self-determined origin, applies to ``cmd`` and ``proc``)
   The first output from the command indicates the origin of
   the tick function.

``m`` (monotonic, applies to ``cmd`` and ``proc``)
   The command reports monotonically increasing values, from a common
   origin. Implies ``o``.

Examples
--------

All the following examples cause a tick event to be generated
every 3 seconds, reporting a +3 tick increase at each event.

The following specifications use ``reg``'s start time as origin:

``-t time:3s``

``-t proc/do:"date +%s; while sleep 3; do echo 3; done"``

``-t proc/m:"while sleep 3; do date +%s; done"``

The following specifications force origin 0:

``-t time/z:3s``

``-t cmd/d:"sleep 3; echo 3"``

``-t proc/doz:"date +%s; while sleep 3; do echo 3; done"``

``-t proc/d:"while sleep 3; do echo 3; done"``

``-t proc/mz:"while sleep 3; do date +%s; done"``


PROGRESS INDICATOR
==================

The argument ``-s`` specifies a progress indicator function, which
maps tick increases into step increases. The general syntax
is ``-s <TYPE>/<FLAGS>:<ARG>``, similarly to ``-t`` above.

Step functions
--------------

The following predefined functions are available:

======================= =====================================
Function                Description
======================= =====================================
``cmd``                 Shell command (individual calls)
``proc``                Shell command (interactive)
``const``               Report constant progress
======================= =====================================

With function ``cmd``, the command given as argument is run at each
tick event. The tick value is provided as command-line argument to the
command. The progress indicator event is generated when the command
terminates, using the value reported on its standard output.

With function ``proc``, the command given as argument is run in the
background.  At each tick event, the tick value is written on the
command's standard input. The progress indicator event is generated
when the process responds on its standard output.

With function ``const``, each tick event is mapped to a constant
number of steps. The function argument determines this number
of steps, and defaults to 0 (no progress). This
feature was originally implemented for debugging ``reg``.

Step function flags
-------------------

The optional ``<FLAGS>`` indicate how the function's values are
translated to tick events.

``z`` (force origin zero)
   Force the sequence of step events to have origin value 0, even if
   the underlying function has a different origin.

``d`` (deltas, applies to ``cmd`` and ``proc``)
   Each output from the command reports the additional
   number of steps elapsed since the last output.

``o`` (self-determined origin, applies to ``cmd`` and ``proc``)
   The origin of the tick function is provided as first input to the
   step function. The first output from the command indicates the
   origin of the step function.

``m`` (monotonic, applies to ``cmd`` and ``proc``)
   The command reports monotonically increasing values, from a common
   origin. Implies ``o``.

Example
-------

The following specification uses process 99298's CPU time as step
function::

  -t cmd/m:"ps -o cputime= -p 99298|tr ':.' '  '|awk '{print \$1*60+\$2+\$3/100. }'"

With this specification, ``reg`` runs the command at every tick
event. The ``ps`` command reports the CPU time of process 99298. The
filtering by ``tr`` and ``awk`` translates ``ps``'s CPU time
formatting into a number of seconds.

MONITOR SPECIFICATION
=====================

The argument ``-m`` specifies a resource function, which
maps tick/step increases into resource usage. The general
syntax is ``-m <TYPE>:<ARG>``.

The following functions are available:

=============== =====================================
Function        Description
=============== =====================================
``cmd``         Shell command (individual calls)
``proc``        Shell command (interactive)
``const``       Report constant resource usage
=============== =====================================

With function ``cmd``, the command given as argument is run at each
tick event. The tick and step values are provided as command-line
arguments to the command. The resource usage event is generated when
the command terminates, using the value reported on its standard
output.

With function ``proc``, the command given as argument is run in the
background.  At each tick event, the tick and step values are written
on the command's standard input, separated by a space. The resource
usage event is generated when the process responds on its standard
output.

With both ``cmd`` and ``proc``, the first input to the command is the origin of
the ticks and steps functions.

With function ``const``, each tick event is mapped to a constant
resource usage. The function argument determines the amount
in stuff units, and defaults to 0 (no resource usage). This
feature was originally implemented for debugging ``reg``.

INPUT LANGUAGE
==============

``reg`` accepts the following newline-terminated commands on its
input stream:

``. <ticks>``
  Generate a tick event with the specified amount of
  ticks. This can be combined with ``-t instant`` to
  place tick generation fully under control of the input stream.

``+ <amount>``
  Add the specified number of stuff.steps to the resource
  supply.

``aon`` / ``aoff``
  Enable/disable reporting supply exhaustion to the actuator.

``?``
  Emit a status record on the output stream.

OUTPUT FORMAT
=============

Each status record ends with a newline character, and is composed of
the following space-separated columns:

- the current tick,
- the tick delta (number of ticks elapsed since the last status record),
- the current step & step delta,
- the current supply & supply delta.

OUTPUT RATE
===========

By default, ``reg`` produces status records after each explicit ``?``
command on the input stream.

Additionally, the option ``-p steps:<N>`` and ``-p ticks:<N>``
instructs ``reg`` to emit records periodically, with the period
specified (either steps or ticks). If the period is zero, a record
is emitted for each ticks/steps event.

``reg`` does not block on output: if the output stream is blocked, the
deltas accumulate until ``reg`` becomes able to output records again. If
more than one ``?`` input commands are received on the input, or periods
of ``-p`` are elapsed while the output stream is blocked, they are
ignored and only one status record is emitted on the output stream
when it becomes unblocked.

With option ``-p flood``, as many status records are generated as
possible when the output stream is unblocked. The consumer process is
then in charge of controlling the rate by throttling its input.

With ``-p none`` the automatic output is disabled and records are only
output when ``?`` is received on the input.  (this is the default).


GRANULARITY
===========

The rate at which ``reg`` monitors ``t`` and makes regulation decisions
is determined by the *granularity* parameter, selected with option
``-g <value>``.

In other words, ``reg`` groups the tick events generated by the time
discretization function so that the minimum increment between
subsequent events is ``<value>``.  For example, with ``-t
time:300ms -g 2``, ``reg`` will coalesce approximately every 6 events
into a single +2 second event.

If ``<value>`` is 0, the granularity is not enforced (all tick events
are used). This is the default.


ACTUATION
=========

When the supply is exhausted, ``reg`` informs the actuator defined by
argument ``-a`` periodically (at every subsequent tick event) until
the supply is provisioned again.

The actuator can be defined by ``-a <TYPE>:<ARG>``. The following actuator types
are supported.

================== ================================================
Actuator           Description
================== ================================================
``print``          Print the current supply status to file.
``cmd``            Shell command (individual calls)
``proc``           Shell command (interactive)
``discard``        Do nothing
================== ================================================

With function ``print``, the current supply status and last
ticks/steps/supply update are printed to the file specified with
``<ARG>`` at each tick event when the supply is exhausted.

With ``cmd``, the shell command is run at each tick event, with the
current ticks/steps/supply update provided as command-line arguments.

With ``proc``, the shell command is run in the background, and the
current ticks/steps/supply update is provided on the command's
standard input at each tick event.

The following actuators have therefore the same effect:

``-a print:/dev/tty``

``-a cmd:'echo $@>/dev/tty'``

``-a proc:'while read a; do echo $a>/dev/tty; done'``

Note: the effect of an actuator should be to stop/throttle the
progress function *p(t)* (e.g. make it constant), so that its integral
resource consumption stays zero until the supply is increased and the
process is restarted.

EXIT STATUS
===========

``reg`` terminates with exit status 0 when its input stream is
exhausted (end-of-file is encounted while reading).

Errors, signals, unknown situations, etc. are reported with other exit
codes.

AUTHOR
======

Writen by Raphael 'kena' Poss.

REPORTING BUGS
==============

Report bugs to: https://github.com/knz/reg/issues
