from robomotion.node import Node
from robomotion.decorators import *
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from pipedrive.client import Client
import pipedrive

from nodes.common import AddClient

@node_decorator(name='Robomotion.Pipedrive.Connect', title='Connect', color='#000000', icon='M10,17V14H3V10H10V7L15,12L10,17M10,2H19A2,2 0 0,1 21,4V20A2,2 0 0,1 19,22H10A2,2 0 0,1 8,20V18H10V20H19V4H10V6H8V4A2,2 0 0,1 10,2Z')
class Connect(Node):
    def __init__(self):
        super().__init__()

        #Input
        self.inDomain = InVariable(title='Company Domain', type='string', scope='Custom', name='', customScope=True, messageScope=True)

        # Output
        self.outConnectionId = OutVariable(title='Connection Id', type='string', scope='Message', name='connection_id', messageOnly=True)

        #Options
        self.optApiKey = Credentials(title='API Key', category=ECategory.Token)

    def on_create(self):
        return

    def on_message(self, ctx: Context):
        domain = self.inDomain.get(ctx)
        if domain == "":
            raise ValueError("Company Domain can not be empty")

        vaultApiKey = self.optApiKey.get_vault_item(ctx)
        if vaultApiKey == "":
            raise ValueError("API Key can not be empty")

        apiKey = vaultApiKey["value"]
        client = Client(domain=domain)
        client.set_api_token(apiKey)
        id = AddClient(client)
        self.outConnectionId.set(ctx, id)

    def on_close(self):
        return