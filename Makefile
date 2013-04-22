all: reg.1 reg.pdf

.SUFFIXES: .1 .rst .dvi .pdf

.rst.1:
	rst2man $< >$@

.1.dvi:
	groff -t -Tdvi -m man $< >$@

.dvi.pdf:
	dvipdf $< $@
