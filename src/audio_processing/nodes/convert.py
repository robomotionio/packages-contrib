from robomotion.node import Node
from robomotion.decorators import *
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum
from robomotion.message import Context, Message
import os
from os import path
from pydub import AudioSegment

@node_decorator(name='Robomotion.AudioProcessing.Convert', title='Convert', color='#000000', icon='M8 4v10.184C7.686 14.072 7.353 14 7 14c-1.657 0-3 1.343-3 3s1.343 3 3 3 3-1.343 3-3V7h7v4.184c-.314-.112-.647-.184-1-.184-1.657 0-3 1.343-3 3s1.343 3 3 3 3-1.343 3-3V4H8z')
class Convert(Node):
    def __init__(self):
        super().__init__()
        self.inSourcePath = InVariable(title='Source Path', type='string', scope='Custom', name='', customScope=True, messageScope=True)
        self.inDestinationPath = InVariable(title='Destination Path', type='string', scope='Custom', name='', customScope=True, messageScope=True)                        

    def on_create(self):
        return

    def on_message(self, ctx: Context):
      
        
        inSourcePath = self.inSourcePath.get(ctx)
        inDestinationPath = self.inDestinationPath.get(ctx)
        
        if type(inSourcePath) != str:
            raise TypeError("Invalid Input. Source Path is not valid string")        

        if type(inDestinationPath) != str:
            raise TypeError("Invalid Input. Destination Path is not valid string")
            
        filename, file_extension = os.path.splitext(inDestinationPath)
        file_extension = file_extension[1:] #first character is '.', so it is removed

        sound = AudioSegment.from_mp3(inSourcePath)        
        sound.export(inDestinationPath, format=file_extension)
    def on_close(self):
        return
