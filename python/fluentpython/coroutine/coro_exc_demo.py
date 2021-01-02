# BEGIN EX_CORO_EXC
class DemoException(Exception):
    """An exception type for the demonstration."""

def demo_exc_handling():
    print('-> coroutine started')
    try:
        while True:
            try:
                x = yield
            except DemoException:  # <1>
                print('*** DemoException handled. Continuing...')
            else:  # <2>
                print('-> coroutine received: {!r}'.format(x))
        raise RuntimeError('This line should never run.')  # <3>
    finally:
        print('-> coroutine ending')
# END EX_CORO_EXC

if __name__ == '__main__':
    exc_coro = demo_exc_handling()
    next(exc_coro)
