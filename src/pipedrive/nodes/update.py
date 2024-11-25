from robomotion.node import Node
from robomotion.decorators import *
from robomotion.runtime import Runtime
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from nodes.common import GetClient
import pipedrive
@node_decorator(name='Robomotion.Pipedrive.Update', title='Update', color='#000000', icon='M21,10.12H14.22L16.96,7.3C14.23,4.6 9.81,4.5 7.08,7.2C4.35,9.91 4.35,14.28 7.08,17C9.81,19.7 14.23,19.7 16.96,17C18.32,15.65 19,14.08 19,12.1H21C21,14.08 20.12,16.65 18.36,18.39C14.85,21.87 9.15,21.87 5.64,18.39C2.14,14.92 2.11,9.28 5.62,5.81C9.13,2.34 14.76,2.34 18.27,5.81L21,3V10.12M12.5,8V12.25L16,14.33L15.28,15.54L11,13V8H12.5Z')
class Update(Node):
    def __init__(self):
        super().__init__()
        
        # Input
        self.inConnectionId = InVariable(title='Connection Id', type='string', scope='Message', name='connection_id', customScope=True, messageScope=True)
        self.inId = InVariable(title='Object Id', type='string', scope='Custom', name='', customScope=True, messageScope=True)
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
        
        id = self.inId.get(ctx)
        if id == "":
            raise ValueError("Id can not be empty")

        data = self.inData.get(ctx)
        if data is None:
            raise ValueError("Data can not be empty")
        
        type = self.optType
        if type == "_" or type == "":
            raise ValueError("Type must be selected")
        
        def switch(type, id, data, client):
            if type == "activity":
                return client.activities.update_activity(id, data)
            elif type == "deal":
                return client.deals.update_deal(id, data)
            elif type == "note":
                return client.notes.update_note(id, data)
            elif type == "organization":
                return client.organizations.update_organization(id, data)
            elif type == "person":
                return client.persons.update_person(id, data)
            elif type == "lead":
                return client.leads.update_lead(id, data)

        client = GetClient(connectionId)
        response = switch(type=type,id=id,data=data,client=client)
        self.outResult.set(ctx, response)

    def on_close(self):
        return

