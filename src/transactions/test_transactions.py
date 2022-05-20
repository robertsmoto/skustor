import unittest
from . import transactions
import locations.locations as locations
import inventory.inventory as inventory


class TestAdvancedShippingNotice(unittest.TestCase):

    def setUp(self):
        self.supplier01 = locations.Company(id=100)
        self.warehouse01 = locations.Warehouse(id=200)

    def test_instantiate_the_class(self):
        asn = transactions.AdvancedShippingNotice(
                number='ASN-001',
                )
        asn.date = '2022-05-14'
        self.assertEqual(asn.date, '2022-05-14')
        self.sender = self.supplier01
        self.receiver = self.warehouse01
        print(self.supplier01, self.warehouse01)


class TestTransactionDetails(unittest.TestCase):

    def setUp(self):
        self.supplier01 = locations.Company(id=100)
        self.warehouse01 = locations.Warehouse(id=200)
        self.asn = transactions.AdvancedShippingNotice(
                number='ASN-001',
                )
        self.part01 = inventory.Part(sku='PART-01')
        self.part02 = inventory.Part(sku='PART-02')
        self.part03 = inventory.Part(sku='PART-03')

    def test_instantiate_trans_details(self):
        detail01 = transactions.TransactionInventoryJoin(
                quantity=5,
                inventory=self.part01
                )
        self.assertEqual(detail01.quantity, 5)


if __name__ == '__main__':
    unittest.main()
