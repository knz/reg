===========================================================
 Abstracting and controlling resource usage as event flows
===========================================================

Introduction
------------

Process models
--------------

"Processes" are an abstraction used in theoretical computer science to
model a running system which communicates with other systems. In
*process models*, (CSP, KPN, Actors, pi-calculus, etc), processes are
connected to each other via channels and exchange messages over these
channels. Some of these models also allow processes to *duplicate*
themselves conditionally and define a new behavior (program) for the
newly created process.

These abstractions are *useful* in that they simplify the
understanding of the observable behavior of processes: all that there
is to know about processes "from the outside world's perspective" can
be modeled as message protocols over the inter-process channels and
what happens upon process creation events.

Applications of process models
------------------------------

Until recently, there have been two main applications of process models:

- protocol proofs in general process networks: deadlock freedom,
  guarantee of termination, guarantee of progress, etc.

  These proofs typically exploit formal *protocol descriptions* (ie
  formal descriptions of the *type, order* and sometimes *arity* of
  messages exchanged between individual processes, and optionally of
  the conditions of process duplication) and derive analytically
  global properties over type, order and arity of messages over the
  entire network.

- performance analysis and prediction in static process networks (ie
  without dynamic process duplication): throughput/latency, critical
  path, etc.

  These methods typically assume that channels and processes are
  characterized by a performance *capacity* (typically bandwidth and
  minimum latency for channels, and start-up time and maximum
  processing rate for processes) and that messages over channels imply
  a discrete *cost* upon this capacity budget.

Resource usage in computer systems
----------------------------------

In contrast to abstract process models, the actual *processes* used in
operating systems have been primarily designed to encapsulate
*resource usage and accounting*.

From these systems perspective, processes have to deal with *primary*
resources, which are axiomatically available depending on the concrete
underlying hardware platform, and *secondary* resources which are
built via abstraction upon primary resources.

Two primary resources are assumed to be available in every computing
system with the same semantics: *storage* ("memory") and *progress in
bounded time*. Storage offers a finite set of addresses to a process
and a protocol which associates data retention semantics to each
address. Progress offers a guarantee that if a process is ready to
execute an instruction, the instruction will execute within a finite
amount of time.

A third primary resource is also available in every computing system,
but potentially with different semantics every time: *input/output
channels with the outside world*. The semantics are different
depending on which *device* happens to be present on the other
side of the channel and which protocol it supports.

On top of these three primary resources, operating systems derive
secondary resources, for example:

- *virtual storage*, defined on top of primary storage and I/O
  channels to storage devices;

- *process contexts*, defined on top of virtual storage, which
  inventorize resources used by each process and enable two or more
  programs to run side-by-side;

- *virtual file systems* with directory, files, symbolic links,
  permissions, etc., defined on top of I/O channels to storage
  devices, and optionally on top of virtual storage;

- *virtual channels* which homogeneize access to files, network
  interfaces and inter-process FIFOs, defined on top of both I/O
  channels and virtual storage;

- *synchronization facilities* that enable processes to explicitly
  wait upon events, either from the outside world or other processes,
  defined on top of virtual storage and the progress contract.

On top of the three "main" primary resources named above, *Unix* and
related operating systems for "general-purpose computing" also
typically require three additional primary resources:

- an *extended progress contract with real time guarantees*, ie a
  "real time clock" able to deliver events at *regular* time
  intervals. Unix needs this to enable programs to wait for a real
  amount of time (``usleep, alarm, setitimer``). This differs from the
  base progress contract identified above, which may deliver progress
  at irregular time intervals as long as they are finite;

- *asynchronous notifications* that trigger the execution of a
  previously agreed program upon reception of an external event,
  regardless of which step of a program is currently running. Unix
  needs this to offer *progress independence* between processes, ie
  the progress of one process cannot be impeded by lack of cooperation
  (or an error) in another process;

- *virtual addressing*, ie a translation facility in hardware between
  logical primary storage addresses used by programs to either
  physical addresses or asynchronous notifications, where the
  translation mapping is not configurable by the program itself. Unix
  needs this to offer the following services:

  - *process isolation*;

  - *unification of virtual file systems and virtual storage* by enabling
    a process to map (parts of) a file as virtual storage addresses;

  - predictable *address space layout* to every program, with "data and
    text" addresses starting from 0 and "stack" addresses starting
    from the maximum address;

  - *virtual storage over-commit*, ie the ability to reserve more
    storage addresses than there are storage cells available.

Note that these resources are primary because they cannot be emulated
on top of the other primary resources named above.

Events for Unix processes
-------------------------


test
