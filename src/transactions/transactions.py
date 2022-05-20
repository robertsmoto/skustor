from dataclasses import dataclass
from abc import ABC, abstractmethod
from inventory.inventory import Inventory
from locations.locations import Location


@dataclass
class Transaction(ABC):
    number: str = ''
    sender: Location = None
    receiver: Location = None
    date: str = ''
    title: str = ''
    description: str = ''

    @abstractmethod
    def create(self):
        pass


@dataclass
class TransactionInventoryJoin:
    quantity: int = 0
    inventory: Inventory = None
    transaction: Transaction = None

    @abstractmethod
    def create(self):
        pass


class AdvancedShippingNotice(Transaction):

    def create(self):
        return self


class Order(Transaction):

    def create(self):
        return self


class Return(Transaction):

    def create(self):
        return self


class ManualEntry(Transaction):

    def create(self):
        return self


class Transfer(ABC):

    def create(self):
        return self
