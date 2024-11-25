from robomotion.node import Node
from robomotion.decorators import *
from robomotion.runtime import Runtime
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from nodes.common import GetClient
import pipedrive
@node_decorator(name='Robomotion.Pipedrive.Get', title='Get', color='#000000', icon='M5,3H19A2,2 0 0,1 21,5V19A2,2 0 0,1 19,21H5A2,2 0 0,1 3,19V5A2,2 0 0,1 5,3M12,17L17,12H14V8H10V12H7L12,17Z')
class Get(Node):
    def __init__(self):
        super().__init__()
        
        # Input
        self.inConnectionId = InVariable(title='Connection Id', type='string', scope='Message', name='connection_id', customScope=True, messageScope=True)
        self.inId = InVariable(title='Object Id', type='string', scope='Custom', name='', customScope=True, messageScope=True)

        # Output
        self.outResult = OutVariable(title='Result', type='object', scope='Message', name='result', messageOnly=True)

        #Options
        self.optType = Variable(title='Type', type='string', enum=_Enum(enums= ["activity","deal","note","organization","person","lead"], enumNames=["Activity","Deal","Note","Organization","Person","Lead"]), default="activity", option=True) 

    def on_create(self):
        return

    def on_message(self, ctx: Context):
        connectionId = self.inConnectionId.get(ctx)
        if connectionId == "":
            raise ValueError("Connection Id can not be empty")
        
        id = self.inId.get(ctx)
        if id == "":
            raise ValueError("Id can not be empty")
        
        type = self.optType
        if type == "_" or type == "":
            raise ValueError("Type must be selected")

        def switch(type, data, client):
            if type == "activity":
                return client.activities.get_activity(data)
            elif type == "deal":
                return client.deals.get_deal(data)
            elif type == "note":
                return client.notes.get_note(data)
            elif type == "organization":
                return client.organizations.get_organization(data)
            elif type == "person":
                return client.persons.get_person(data)
            elif type == "lead":
                return client.leads.get_lead(data)

        client = GetClient(connectionId)
        response = switch(type=type,data=id,client=client)
        self.outResult.set(ctx, response)

    def on_close(self):
        return

