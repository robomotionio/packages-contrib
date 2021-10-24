from robomotion.node import Node
from robomotion.decorators import *
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum
from robomotion.message import Context, Message
import Levenshtein as lev
from fuzzywuzzy import fuzz

@node_decorator(name='Robomotion.FuzzyMatcher.FuzzyMatcher', title='FuzzyMatcher', color='#F56040', icon='M18 8c0-3.31-2.69-6-6-6S6 4.69 6 8c0 4.5 6 11 6 11s6-6.5 6-11zm-8 0c0-1.1.9-2 2-2s2 .9 2 2-.89 2-2 2c-1.1 0-2-.9-2-2zM5 20v2h14v-2H5z')
class FuzzyMatcher(Node):
    def __init__(self):
        super().__init__()
        self.inRealData = InVariable(title='Real Data', type='string', scope='Custom', name='', customScope=True, messageScope=True)
        self.inTestData = InVariable(title='Test Data', type='string', scope='Custom', name='', customScope=True, messageScope=True)
        self.OutSimilarity = OutVariable(title='Similarity', type='string', scope='Message', name='similarity', messageOnly=True)
        
        
        

    def on_create(self):
        return

    def on_message(self, ctx: Context):
        #lower() converts I -> i , other Turkish characters are converted correctly. So for 'I' character, we translated
        #it to 'i', then have used lower() function
        turkish_lower_map = {
            ord(u'I'): u'ı'
        }
        
        inRealData = self.inRealData.get(ctx)
        inTestData = self.inTestData.get(ctx)
        
        if type(inRealData) != str:
            raise TypeError("Invalid Input. Real Data is not valid string")
        inRealData = inRealData.translate(turkish_lower_map).lower()

        if type(inTestData) != str:
            raise TypeError("Invalid Input. Test Data is not valid string")
        inTestData = inTestData.translate(turkish_lower_map).lower()
    
        levenshteinRatio = lev.ratio(inRealData,inTestData)
        levenshteinRatio *= 100
        wuzzyRatio = fuzz.ratio(inRealData, inTestData)

        similarity = (levenshteinRatio + wuzzyRatio)/2
        self.OutSimilarity.set(ctx, similarity)

    def on_close(self):
        return
