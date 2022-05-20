from dataclasses import dataclass
from abc import ABC, abstractmethod


@dataclass
class Inventory(ABC):
    sku: str  # use as id
    name: str = ''
    description: str = ''
    keywords: str = ''  # comma-separated list

    cost: int = 0
    unit: str = ''

    quantity: int = 0

    length: int = 0
    width: int = 0
    height: int = 0

    identifier: str = ''
    gtin: str = ''
    isbn: str = ''

    @abstractmethod
    def create(self):
        pass


class RawMaterial(Inventory):

    def create(self):
        return self


class Part(Inventory):

    def create(self):
        return self


class Product(Inventory):

    price: int = 0

    def create(self):
        return self
