all: reg.1 reg.pdf reg design.html design1.png design2.png

.SUFFIXES: .1 .rst .dvi .pdf .dot .png .html

.dot.png:
	dot -Tpng -o $@ $<

.rst.html:
	rst2html $< >$@

.rst.1:
	rst2man $< >$@

.1.dvi:
	groff -t -Tdvi -m man $< >$@

.dvi.pdf:
	dvipdf $< $@

.PHONY: reg
reg:
	GOPATH=$$PWD go build -o reg main
