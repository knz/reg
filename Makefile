all: reg.1 reg.pdf reg design.html design1.png 

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

.PHONY: reg version.go

version.go:
	echo 'package main; var version = "'`git describe --all --long`" $$USER@"`uname -n``date +/%F/%T`'"' >$@

reg: version.go
	rm -rf pkg
	GOPATH=$$PWD go get
	GOPATH=$$PWD go build -o reg reg.go version.go
