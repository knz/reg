======================
 Implementing ``reg``
======================

:Author: kena
:Date: April 2013

Overview
========

The following diagram shows the overall structure of ``reg``'s process
network:

.. image:: design1.png
   :align: center
   :width: 20%

The following processes are defined:

``input``
   Reads lines of text from the input stream and provides them to ``parse``.

``output``
   Tests availability of the output stream, outputs records.

``actuator``
   Acts upon the monitored system.

``ticksource``, ``stepsource``, ``sampler``
   Generate tick, step and measurement events for ``integrate``. Optionally
   generates tick and step events for ``throttle`` if ``-p`` is used.

``parse``
   Analyses commands from the input stream. Depending on the command,
   generates either:

   - supply events to ``integrate`` (commands ``+``, ``aon``, ``aoff``)

   - report events to ``outmgt`` (command ``?``)

   - tick events to ``mergeticks`` (command ``.``)

``outmgt``
   Formats status reports for ``output``, by querying ``integrate``
   for the current status of the supply bin(s).

``integrate``
  Consumes tick, step and measurement events and updates the supply bin(s).
  Generates action events for ``actuator`` and answers status requests
  from ``outmgt``.
