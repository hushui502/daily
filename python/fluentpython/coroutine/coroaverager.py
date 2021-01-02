from collections import namedtuple

Result = namedtuple('Result', 'count average')

def averager():
    total = 0.0
    count = 0
    average = None
    while True:
        term = yield
        if term is None:
            break
        total += term
        count += 1
        average = total/count

    return Result(count, average)

if __name__ == '__main__':
    coro_avg = averager()
    next(coro_avg)
    try:
        coro_avg.send(None)
    except StopIteration as exc:
        result = exc.value

