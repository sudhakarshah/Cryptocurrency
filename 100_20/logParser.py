import numpy as np
import glob
from subprocess import Popen, PIPE, STDOUT
from matplotlib import pyplot as plt


ACTION = 0
TIMESTAMP = 1
TYPE = 2
SIZE = 3
TRANSMISSION_TIME = 4
TID = 3
MEMBER_SIZE = 4
HASHTABLE_SIZE = 5

RECIEVED = "RECIEVED"
SEND = "SEND"
UPDATE = "UPDATE"


def totalSends(log):
    return

    
def tokenize(log_line):
    return log_line.split(' ')

if __name__ == "__main__":



    logFiles = glob.glob("./node*.log")
    contents = []
    for name in logFiles:
        lines = open(name, "r").read().split('\n')
        lines = list(map(tokenize, lines))
        contents.append(lines)
    print("There are %d log files"%(len(contents)))


    haveFollowing = {}
    for log in contents:
        stuff = {}
        for l in log:
            if l[0] == "#" or len(l) < 2:
                continue
            if l[ACTION] == UPDATE:
                if l[TID] not in haveFollowing:
                    haveFollowing[l[TID]] = 0
                haveFollowing[l[TID]] += 1
                stuff[l[TID]] = 1
        print(len(stuff))
    print("Total:%d"%(len(haveFollowing)))

    transaction_delay = {}
    for log in contents:
        for l in log:
            if l[0] == "#" or len(l) < 2:
                continue
            if l[ACTION] == SEND and l[TYPE] == "TRANSACTION":
                tid = l[-1]
                duration = l[TRANSMISSION_TIME]
                if tid not in transaction_delay:
                    transaction_delay[tid] = []
                transaction_delay[tid].append(int(duration))
    delay_stat = {}
    minArray = []
    medArray = []
    maxArray = []
    for k, v in transaction_delay.items():
        median = np.median(np.asarray(v))
        minimum = min(v)
        maximum = max(v)
        minArray.append(minimum)
        maxArray.append(maximum)
        medArray.append(median)
        delay_stat[k] = (minimum, median, maximum)

    x = np.arange(len(minArray))
    plt.plot(x, minArray)
    plt.plot(x, medArray)
    plt.plot(x, maxArray)
    plt.legend(['min','med','max'], loc='upper left')
    plt.savefig("prop.png")
