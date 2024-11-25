from robomotion.node import Node
from robomotion.decorators import *
from robomotion.runtime import Runtime
from robomotion.variable import Variable, InVariable, OutVariable, OptVariable, Credentials, ECategory, _Enum, _DefVal
from robomotion.message import Context, Message
from nodes.common import GetClient
import pipedrive
@node_decorator(name='Robomotion.Pipedrive.GetAll', title='Get All', color='#000000', icon='M7,13V11H21V13H7M7,19V17H21V19H7M7,7V5H21V7H7M3,8V5H2V4H4V8H3M2,17V16H5V20H2V19H4V18.5H3V17.5H4V17H2M4.25,10A0.75,0.75 0 0,1 5,10.75C5,10.95 4.92,11.14 4.79,11.27L3.12,13H5V14H2V13.08L4,11H2V10H4.25Z')
class GetAll(Node):
    def __init__(self):
        super().__init__()
        
        # Input
        self.inConnectionId = InVariable(title='Connection Id', type='string', scope='Message', name='connection_id', customScope=True, messageScope=True)

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
        
        type = self.optType
        if type == "_" or type == "":
            raise ValueError("Type must be selected")

        def switch(type, client):
            if type == "activity":
                return client.activities.get_all_activities()
            elif type == "deal":
                return client.deals.get_all_deals()
            elif type == "note":
                return client.notes.get_all_notes()
            elif type == "organization":
                return client.organizations.get_all_organizations()
            elif type == "person":
                return client.persons.get_all_persons()
            elif type == "lead":
                return client.leads.get_all_leads()

        client = GetClient(connectionId)
        response = switch(type=type,client=client)
        self.outResult.set(ctx, response)

    def on_close(self):
        return

