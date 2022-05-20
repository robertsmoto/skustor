from dataclasses import dataclass
from abc import ABC, abstractmethod
from people.people import Person


@dataclass
class Location(ABC):
    id: int = 0
    name: str = ''
    description: str = ''
    phone_number: str = ''
    email: str = ''
    website: str = ''

    @abstractmethod
    def create(self):
        pass


@dataclass
class LocationPersonJoin:
    location: Location = None
    person: Person = None


class Website(Location):

    def __init__(self, domain):
        self.domain = domain

    def create(self):
        return self


class Warehouse(Location):

    def create(self):
        return self


class Store(Location):

    def create(self):
        return self


class Company(Location):

    def create(self):
        return self
