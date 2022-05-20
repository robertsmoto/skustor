from dataclasses import dataclass
from abc import ABC, abstractmethod


@dataclass
class Person(ABC):
    id: int = 0
    firstname: str = ''
    lastname: str = ''
    nickname: str = ''
    phone_number: str = ''
    mobile_number: str = ''
    email: str = ''

    @abstractmethod
    def create(self):
        pass


class Customer(Person):

    def create(self):
        return self


class Contact(Person):

    def create(self):
        return self
