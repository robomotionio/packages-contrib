from robomotion.node import Node
from robomotion.decorators import *
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum
from robomotion.message import Context, Message
from os import path
from pydub import AudioSegment

@node_decorator(name='Robomotion.AudioProcessing.Convert', title='Convert', color='#F56040', icon='M18 8c0-3.31-2.69-6-6-6S6 4.69 6 8c0 4.5 6 11 6 11s6-6.5 6-11zm-8 0c0-1.1.9-2 2-2s2 .9 2 2-.89 2-2 2c-1.1 0-2-.9-2-2zM5 20v2h14v-2H5z')
class Convert(Node):
    def __init__(self):
        super().__init__()
        self.inSourcePath = InVariable(title='Source Path', type='string', scope='Custom', name='', customScope=True, messageScope=True)
        self.inDestionationPath = InVariable(title='Destionation Path', type='string', scope='Custom', name='', customScope=True, messageScope=True)                        

    def on_create(self):
        return

    def on_message(self, ctx: Context):
      
        
        inSourcePath = self.inSourcePath.get(ctx)
        inDestionationPath = self.inDestionationPath.get(ctx)
        
        if type(inSourcePath) != str:
            raise TypeError("Invalid Input. Source Path is not valid string")        

        if type(inDestionationPath) != str:
            raise TypeError("Invalid Input. Destionation Path is not valid string")
            
        filename, file_extension = os.path.splitext(inDestionationPath)
        
        sound = AudioSegment.from_mp3(src)        
        sound.export(dst, format=file_extension)
    def on_close(self):
        return
