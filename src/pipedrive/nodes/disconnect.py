from robomotion.node import Node
from robomotion.decorators import *
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from nodes.common import DeleteClient

@node_decorator(name='Robomotion.Pipedrive.Disconnect', title='Disconnect', color='#000000', icon='M16,17V14H9V10H16V7L21,12L16,17M14,2A2,2 0 0,1 16,4V6H14V4H5V20H14V18H16V20A2,2 0 0,1 14,22H5A2,2 0 0,1 3,20V4A2,2 0 0,1 5,2H14Z')
class Disconnect(Node):
    def __init__(self):
        super().__init__()

        # Input
        self.inConnectionId = InVariable(title='Connection Id', type='string', scope='Message', name='connection_id', customScope=True, messageScope=True)

    def on_create(self):
        return

    def on_message(self, ctx: Context):
        connectionId = self.inConnectionId.get(ctx)
        if connectionId == "":
            raise ValueError("Connection Id can not be empty")

        DeleteClient(connectionId)
        

    def on_close(self):
        return