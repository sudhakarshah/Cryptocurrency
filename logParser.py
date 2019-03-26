import glob
from subprocess import Popen, PIPE, STDOUT


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
                haveFollowing[l[TID]] = 1
                stuff[l[TID]] = 1
        print(len(stuff))
    print("Total:%d"%(len(haveFollowing)))

    exit()
    
    # fetch service log
    slog = open("service.log", "r").read().split('\n')
    slog = list(map(tokenize, slog))

    total = 0
    hashtable = {}
    for l in slog:
        if len(l) < 2:
            continue
        if l[1] == "transaction":
            hashtable[l[3]] = 1
    print(len(hashtable))
