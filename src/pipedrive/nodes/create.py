from robomotion.node import Node
from robomotion.decorators import *
from robomotion.runtime import Runtime
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from nodes.common import GetClient
import pipedrive
@node_decorator(name='Robomotion.Pipedrive.Create', title='Create', color='#000000', icon='M19,1L17.74,3.75L15,5L17.74,6.26L19,9L20.25,6.26L23,5L20.25,3.75M9,4L6.5,9.5L1,12L6.5,14.5L9,20L11.5,14.5L17,12L11.5,9.5M19,15L17.74,17.74L15,19L17.74,20.25L19,23L20.25,20.25L23,19L20.25,17.74')
class Create(Node):
    def __init__(self):
        super().__init__()
        
        # Input
        self.inConnectionId = InVariable(title='Connection Id', type='string', scope='Message', name='connection_id', customScope=True, messageScope=True)
        self.inData = InVariable(title='Data', type='object', scope='Message', name='data', messageScope=True)

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
        
        data = self.inData.get(ctx)
        if data is None:
            raise ValueError("Data can not be empty")
        
        type = self.optType
        if type == "_" or type == "":
            raise ValueError("Type must be selected")

        def switch(type, data, client):
            if type == "activity":
                return client.activities.create_activity(data)
            elif type == "deal":
                return client.deals.create_deal(data)
            elif type == "note":
                return client.notes.create_note(data)
            elif type == "organization":
                return client.organizations.create_organization(data)
            elif type == "person":
                return client.persons.create_person(data)
            elif type == "lead":
                return client.leads.create_lead(data)

        client = GetClient(connectionId)
        response = switch(type=type,data=data,client=client)
        self.outResult.set(ctx, response)

    def on_close(self):
        return

