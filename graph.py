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
#
BLOCK_HASH = 2
BLOCK_TRANS = -1

RECIEVED = "RECIEVED"
SEND = "SEND"
UPDATE = "UPDATE"
ACCEPTED = "ACCEPTED"


def totalSends(log):
    return


def tokenize(log_line):
    return log_line.split(' ')


if __name__ == "__main__":

    logFiles = glob.glob("./block_100_20/node*.log")
    contents = []
    for name in logFiles:
        lines = open(name, "r").read().split('\n')
        lines = list(map(tokenize, lines))
        contents.append(lines)
    print("There are %d log files" % (len(contents)))

    transaction_first_occ = {}
    block_first_occ = {}
    block_to_trans = {}
    block_last_occ = {}
    split_events = []
    for log in contents:
        for l in log:
            if l[0] == "#" or len(l) < 2 or not l[0][:1].isalpha():
                continue
            if l[ACTION] == UPDATE:
                if l[TID] not in transaction_first_occ:
                    transaction_first_occ[l[TID]] = float("inf")
                if float(l[TIMESTAMP]) < transaction_first_occ[l[TID]]:
                    transaction_first_occ[l[TID]] = float(l[TIMESTAMP])
            if l[ACTION] == ACCEPTED:
                # First occ
                if l[BLOCK_HASH] not in block_first_occ:
                    block_first_occ[l[BLOCK_HASH]] = float("inf")
                if float(l[TIMESTAMP]) < block_first_occ[l[BLOCK_HASH]]:
                    block_first_occ[l[BLOCK_HASH]] = float(l[TIMESTAMP])
                # Last occ
                if l[BLOCK_HASH] not in block_last_occ:
                    block_last_occ[l[BLOCK_HASH]] = 0.0
                if float(l[TIMESTAMP]) > block_last_occ[l[BLOCK_HASH]]:
                    block_last_occ[l[BLOCK_HASH]] = float(l[TIMESTAMP])
                # Map block to array of transactions
                transactions = set(l[BLOCK_TRANS].split(","))
                block_to_trans[l[BLOCK_HASH]] = transactions
            if l[ACTION] == "CHAIN_SPLIT":
                split_events.append(float(l[TIMESTAMP]))

    block_ordering = list(block_first_occ.keys())

    haveFollowing = {}
    for log in contents:
        for l in log:
            if l[0] == "#" or len(l) < 2 or not l[0].isalpha():
                continue
            if l[ACTION] == ACCEPTED:
                if l[BLOCK_HASH] not in haveFollowing:
                    haveFollowing[l[BLOCK_HASH]] = 0
                haveFollowing[l[BLOCK_HASH]] += 1
    print("Total:%d" % (len(haveFollowing)))
    print(haveFollowing)

    block_occ_diff = []
    for k in block_ordering:
        diff = block_last_occ[k] - block_first_occ[k]
        block_occ_diff.append(diff)

    transaction_appear = []
    ordered_trans = [(k, transaction_first_occ[k]) for k in sorted(transaction_first_occ, key=transaction_first_occ.get)]
    for k, v in ordered_trans:
        shortest = float("inf")
        for bh, ti in block_first_occ.items():
            transactions = block_to_trans[bh]
            if k in transactions:
                diff = ti - v
                shortest = min(shortest, diff)
        transaction_appear.append(shortest)

    x = np.arange(len(block_occ_diff))
    plt.figure(0)
    plt.plot(x, block_occ_diff)
    plt.savefig("block_prop.png")

    plt.figure(1)
    x = np.arange(len(haveFollowing))
    y = []
    for k in block_ordering:
        y.append(haveFollowing[k])
    plt.scatter(x, y)
    plt.savefig("block_count.png")

    plt.figure(2)
    x = np.arange(len(transaction_appear))
    plt.scatter(x, transaction_appear)
    plt.savefig("trans_to_block_prop.png")

    plt.figure(3)
    plt.hist(split_events, bins=200)
    plt.savefig("split_frequency.png")
