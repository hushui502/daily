import time


def clock(func):
    def clocked(*args):
        t0 = time.time()
        result = func(*args)
        elapsed = time.time() - t0
        name = func.__name__
        arg_str = ', '.join(repr(arg) for arg in args)
        print('[%0.8f] %s(%s) -> %r' % (elapsed, name, arg_str, result))
        return result
    return clocked

@clock
def snnoze(seconds):
    time.sleep(seconds)

@clock
def factorial(n):
    return 1 if n < 2 else n * factorial(n-1)

if __name__ == '__main__':
    print('*' * 40, 'Calling snooze(.123)')
    snnoze(.123)
    print('*' * 40, 'Calling factorial(6)')
    print('6 != ', factorial(6))