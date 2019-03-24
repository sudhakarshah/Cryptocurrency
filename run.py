import os
from time import sleep
from subprocess import Popen

intro_port = 9999
intro_ip = "0.0.0.0"
rate = 1
nodes = 100


commands = []
dir_path = os.path.dirname(os.path.realpath(__file__))

for i in range(nodes):
    cs = "%s/mp2 %s %d %d node%d > %s/node%i.log" % (dir_path, intro_ip, intro_port, i+8000, i, dir_path, i)
    commands.append(cs)

procs = []
for i in commands:
    procs.append(Popen(i,shell=True))
    sleep(0.05)
for p in procs:
    p.wait()
    print(p)

