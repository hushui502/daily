import re
import reprlib

RE_WROD = re.compile('\w+')

class Sentence:

    def __init__(self, text):
        self.text = text

    def __repr__(self):
        return 'Sentence(%s)' % reprlib.repr(self.text)

    # lazy
    # def __iter__(self):
    #     for match in RE_WROD.finditer(self.text):
    #         yield match.group()

    def __iter__(self):
        return (match.group for match in RE_WROD.finditer(self.text))