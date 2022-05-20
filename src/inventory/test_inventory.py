import unittest
from . import inventory


class TestRawMaterial(unittest.TestCase):

    def test_instantiate_raw_material(self):
        raw_material = inventory.RawMaterial(sku='RM001')
        raw_material.name = 'Raw Material 001'
        self.assertEqual(raw_material.sku, 'RM001')
        print('raw material', raw_material)

    def test_instantiate_product(self):
        product = inventory.Product(sku='PR001')
        product.price = 21_00
        self.assertEqual(product.price, 2100)
        print("product", product)


if __name__ == '__main__':
    unittest.main()
