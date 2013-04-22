#! /usr/bin/env python3

import fcntl
import sys
import os
import re

################################################################################
# Utility functions

re_dec = r'[+-]?(?:\d+\.\d*|\d*\.\d+|\d+)(?:[eE][+-]?\d+)?'
re_hex = r'[+-]?0x(?:H+\.H*|H*\.H+|H+)(?:[pP][+-]?H+)?'.replace('H', '[0-9a-fA-F]')
re_si = r'[YZEPTGMK]i?|[khdcmnpfazy]|da|mu|µ'
si_multipliers = {
    'Y': 1e24,            'Yi': float(2**80),
    'Z': 1e21,            'Zi': float(2**70),
    'E': 1e18,            'Ei': float(2**60),
    'P': 1e15,            'Pi': float(2**50),
    'T': 1e12,            'Ti': float(2**40),
    'G': 1e9,             'Gi': float(2**30),
    'M': 1e6,             'Mi': float(2**20),
    'K': 1e3, 'k': 1e3,   'Ki': float(2**10),
    'h': 100.,
    'da': 10.,
    'd': 1e-1,
    'c': 1e-2,
    'm': 1e-3,
    'mu': 1e-6, 'µ': 1e-6,
    'n': 1e-9,
    'p': 1e-12,
    'f': 1e-15,
    'a': 1e-18,
    'z': 1e-21,
    'y': 1e-24
}

re_mult = r'(?:(?P<hex>%s)|(?P<dec>%s))?(?P<si>%s)?' % (re_hex, re_dec, re_si)
def parse_multiplier(match):
    d = match.groupdict()
    mult = 1.
    if d['hex'] is not None:
        mult = float.fromhex(d['hex'])
    elif d['dec'] is not None:
        mult = float(d['dec'])

    if d['si'] is not None:
        mult *= si_multipliers[d['si']]

    return mult

_re_number = re.compile(re_mult + '$')
def parse_number(domain, value):
    m = _re_number.match(value)
    if m is None:
        domain.argparser.error('invalid number: %s' % value)
    return parse_multiplier(m)

re_function = r'(?:%s\.)?(?P<fun>.*)' % re_mult
_re_function = re.compile(re_function)
def decompose_function(domain, value):
    m = _re_function.match(value)
    if m is None:
        domain.argparser.error('invalid function specification: %s' % value)
    return (parse_multiplier(m), m.group('fun'))

def get_tick_function(mult, spec):
    return (mult, spec)
def get_step_function(mult, spec):
    return (mult, spec)


################################################################################
# Domains

class Domain(object):
    def __init__(self, argparser, label):
        self.argparser = argparser
        self.label = label
        self._ticks = None
        self._steps = None
        self._resource = []
        self.input = None

        self.output = 1 # stdout
        self.granularity = 1 # unit unchanged
        self.outputrate = "none" # don't output automatically

    @property
    def outputrate(self):
        return self._outputrate

    @outputrate.setter
    def outputrate(self, value):
        if value is "none":
            self._outputrate = -1
        else:
            self._outputrate = parse_number(self, value)

    @property
    def resource(self):
        return self._resource
    @resource.setter
    def resource(self, value):
        if isinstance(value, str):
            if ':' not in value:
                self.argparser.error('invalid resource format: %s (missing label?)' % value)
            label, value = value.split(':', 1)
            value = (label, value)
        self._resource.append(value)

    @property
    def granularity(self):
        return self._granularity
    @granularity.setter
    def granularity(self, value):
        if isinstance(value, str):
            value = parse_number(self, value)
        self._granularity = value

    @property
    def input(self):
        return self._input
    @input.setter
    def input(self, value):
        if value is '-':
            self._input = 0
        else:
            self._input = value

    @property
    def output(self):
        return self._output
    @output.setter
    def output(self, value):
        if value is '-':
            self._output = 1
        else:
            self._output = value

    @property
    def ticks(self):
        return self._ticks
    @ticks.setter
    def ticks(self, value):
        if isinstance(value, str):
            value = get_tick_function(*decompose_function(self, value))
        self._ticks = value

    @property
    def steps(self):
        return self._steps
    @steps.setter
    def steps(self, value):
        if isinstance(value, str):
            value = get_step_function(*decompose_function(self, value))
        self._steps = value

    def validate(self):
        for k in ('ticks', 'steps', 'input'):
            if getattr(self, k) is None:
                self.argparser.error('%s not specified for domain %s (missing -%s?)' % (k, self.label, k[0]))
        if len(self.resource) == 0:
            self.argparser.error('no resources specified for domain %s (missing -r?)' % self.label)

        self.inputfile = open(self.input, 'rb', buffering=0)
        self.outputfile = open(self.output, 'wb', buffering=0)
        fcntl.fcntl(self.inputfile, fcntl.F_SETFL, fcntl.fcntl(self.inputfile, fcntl.F_GETFL) | os.O_NDELAY)
        fcntl.fcntl(self.outputfile, fcntl.F_SETFL, fcntl.fcntl(self.outputfile, fcntl.F_GETFL) | os.O_NDELAY)

###
# bl



################################################################################
# Domain registry

domains = {}

def print_domains(file):
    columns = ('granularity', 'outputrate', 'ticks', 'steps', 'input', 'output')
    print('# label\t' + '\t'.join(columns), file=file)
    for dom in domains.values():
        print('"' + dom.label + '": {\t' + '\t'.join((repr(getattr(dom, x)) + ', ' for x in columns)) + '},', file=file)

def print_config(args, file):
    print('{ "protocol" : "%s",\n' \
          '  "follow" : "%s",\n' \
          '  "domains" : {' % (args.protocol, args.follow), file=file)
    print_domains(file)
    print("} }", file=file)

################################################################################
# Main program

if __name__ == "__main__":
    import functools
    import argparse
    import re

    class NewDomain(argparse.Action):
        def __call__(self, parser, namespace, label, option_string = None):
            global domains
            if label in domains:
                parser.error('domain %s already defined' % label)
            domains[label] = Domain(parser, label)

    re_labelprefix = re.compile(r'\w+=')
    class SetDomainProperty(argparse.Action):
        def __call__(self, parser, namespace, value, option_string = None):
            global domains
            label = 'default'
            if re_labelprefix.match(value) is not None:
                label, value = value.split('=', 1)
            if label not in domains:
                parser.error('domain %s not defined (missing -d?)' % label)
            dom = domains[label]

            # set the property
            setattr(domains[label], self.dest, value)

    parser = argparse.ArgumentParser(description='Process REGulator.',
                                     epilog='For more information, see reg(1).')

    defdom = Domain(parser, 'default')
    defdom.ticks = 'realseconds'
    defdom.steps = 'userseconds'
    defdom.input = '-'
    domains['default'] = defdom

    add_arg = parser.add_argument

    add_arg('cmd', metavar='CMD', type=str, nargs='*', help='command to run')

    add_arg('-a', '--attach', metavar='PID',
            help='attach to an already running process or thread')
    add_arg('-f', '--follow', metavar='PRED',
            help='ask %(metavar)s whether to follow child processes')
    add_arg('-v', '--verbose', default=0, action='count',
            help='verbose execution')
    add_arg('-p', '--protocol', metavar='SPEC', default='freeze',
            help='use %(metavar)s as protocol to regulate the process')

    add_arg('-d', '--domain', metavar='LABEL', action=NewDomain,
            help="define a new management domain")
    add_arg('-g', '--granularity', metavar='N', action=SetDomainProperty,
            help="set the regulation granularity to %(metavar)s")
    add_arg('-i', '--input', metavar='FILE', action=SetDomainProperty,
            help='use %(metavar)s as input stream')
    add_arg('-o', '--output', metavar='FILE', action=SetDomainProperty,
            help='use %(metavar)s as output stream')
    add_arg('-r', '--resource', metavar='LABEL:SPEC', action=SetDomainProperty,
            help='define a resource labeled LABEL with function SPEC')
    add_arg('-R', '--outputrate', metavar='N', action=SetDomainProperty,
            help='set the output rate for status records')
    add_arg('-s', '--steps', metavar='SPEC', action=SetDomainProperty,
            help='use %(metavar)s as progress indicator function')
    add_arg('-t', '--ticks', metavar='SPEC', action=SetDomainProperty,
            help='use %(metavar)s as time discretization function')

    args = parser.parse_args()

    # consistency check
    if args.attach is not None and len(args.cmd) > 0:
        parser.error('cannot specify both -a and command to execute')
    for label in domains:
        domains[label].validate()

    if args.verbose:
        print_config(args, sys.stderr)
