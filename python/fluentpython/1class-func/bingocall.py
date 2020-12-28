import random


class BingoCage:

    def __init__(self, items):
        self._items = list(items)
        random.shuffle(self._items)

    def pick(self):
        try:
            return self._items.pop()
        except IndexError:
            raise LookupError('pick from empty bingocage')

    def __call__(self):
        return self.pick()

    def tag(name, *content, cls=None, **attrs):
        if cls is not None:
            attrs['class'] = cls
        if attrs:
            attr_str = ''.join('%s=%s' % (attr, value)
                               for attr, value
                               in sorted(attrs.items()))
        else:
            attr_str = ''
        if content:
            return '\n'.join('<%s%s>%s</%s>' %
                             (name, attr_str, c, name) for c in content)
        else:
            return '<%s%s />' % (name, attr_str)