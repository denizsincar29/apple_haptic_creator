import unittest
from ahap import freq

class TestFreq(unittest.TestCase):
    #def setUp(self) -> None:
    
    def test_freq(self):
        for i in range(80, 231):
            try:
                freq(i, False)
            except Exception as e:
                self.fail(f'the {i} hz freq not converted., exception: {e}')

    def test_raisefreq(self):
        with self.assertRaises(ValueError):
            freq(79, False)
            freq(231, False)

if __name__=="__main__":
    unittest.main()