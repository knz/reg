=========================
 reg: a system REGulator
=========================

This is a tool to help regulate resource usage of a system
asynchronously.

For more information, see the manual page reg(1) 
(generated from ``reg.rst``).

Requirements
------------

To build this software the following components are needed:

- ``make``
- A Go compiler and tool chain (available from the command ``go``)
- python-docutils (reStructured text tools, for documentation)

Installation
------------

To build the ``reg`` executable::

   make reg

To build ``reg`` and the documentation::

   make

After building, ``reg`` is ready to use.

