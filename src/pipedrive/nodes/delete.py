from robomotion.node import Node
from robomotion.decorators import *
from robomotion.runtime import Runtime
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from nodes.common import GetClient
import pipedrive
@node_decorator(name='Robomotion.Pipedrive.Delete', title='Delete', color='#000000', icon='M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z')
class Delete(Node):
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

        def switch(type, id, client):
            if type == "activity":
                return client.activities.delete_activity(id)
            elif type == "deal":
                return client.deals.delete_deal(id)
            elif type == "note":
                return client.notes.delete_note(id)
            elif type == "organization":
                return client.organizations.delete_organization(id)
            elif type == "person":
                return client.persons.delete_person(id)
            elif type == "lead":
                return client.leads.delete_lead(id)

        client = GetClient(connectionId)
        response = switch(type=type,id=id,client=client)
        self.outResult.set(ctx, response)

    def on_close(self):
        return

