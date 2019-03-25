import os
import sys
from time import sleep
from subprocess import Popen

intro_port = 9999
intro_ip = "0.0.0.0"
rate = 20
nodes = 10

prefix = ""
if len(sys.argv) >= 2:
    prefix = sys.argv[1]
else:
    print("Need an argument $(LOG_PREFIX). Got %s"%(sys.argv))
    exit()
if len(sys.argv) == 3:
    intro_ip = "sp19-cs425-g62-01.cs.illinois.edu"

commands = []
dir_path = os.path.dirname(os.path.realpath(__file__))


for i in range(nodes):
    cs = "%s/mp2 %s %d %d node%d > %s/node%s%d.log" % (dir_path, intro_ip, intro_port, i+8000, i, dir_path, prefix, i)
    commands.append(cs)

procs = []
for i in commands:
    procs.append(Popen(i,shell=True))
for p in procs:
    p.wait()
    print(p)

