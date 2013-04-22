======================
 Implementing ``reg``
======================

:Author: kena
:Date: April 2013

Overview
========

The following diagram shows the overall structure of ``reg``'s process
network:

.. image:: design2.png
   :align: center
   :width: 20%

The following processes are defined:

``input``
   Reads lines of text from the input stream and provides them to ``parse``.

``output``
   Tests availability of the output stream, outputs records.

``protocol``
   Acts upon the harnessed process(es).

``measure``
   Generates tick, step and measurement events for ``integrate``. Optionally
   generates tick and step events for ``report`` if ``-R`` is used.

``parse``
   Analyses commands from the input stream. Depending on the command,
   generates either:

   - supply events to ``integrate`` (commands ``+`` and ``-``),

   - report events to ``report`` (command ``?``).

   - tick events to ``measure`` if flag ``-t controlled`` is used (command ``.``).

``report``
   Formats status reports for ``output``, by querying ``integrate``
   for the current status of the supply bin(s).

``integrate``
  Consumes tick, step and measurement events and updates the supply bin(s).
  Generates action events for ``protocol`` and answers status requests
  from ``report``.

Detailed network
================

.. image:: design1.png
   :align: center
   :width: 40%
