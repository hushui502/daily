import copy
import weakref

class Bus:
    def __init__(self, passengers=None):
        if passengers is None:
            self.passengers = []
        else:
            self.passengers = list(passengers)

    def pick(self, name):
        self.passengers.append(name)

    def drop(self, name):
        self.passengers.remove(name)


bus1 = Bus(['alice', 'bill'])
bus2 = copy.copy(bus1)
bus3 = copy.deepcopy(bus1)
id(bus1)
id(bus2)
id(bus3)
bus1.drop('bill')
bus2.passengers


class TwilightBus:

    def __init__(self, passengers=None):
        if passengers is None:
            self.passengers = []
        else:
            self.passengers = list(passengers)

    def pick(self, name):
        self.passengers.append(name)

    def drop(self, name):
        self.passengers.remove(name)

class Cheese:
    def __init__(self, kind):
        self.kind = kind

    def __repr__(self):
        return 'Cheese(%r)' % self.kind

stock = weakref.WeakKeyDictionary()
catalog = [Cheese('red'), Cheese('blue')]
for cheese in catalog:
    stock[cheese.kind] = cheese

sorted(stock.keys())
