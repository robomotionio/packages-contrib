import uuid

Clients = {}
def AddClient (Client):
    id = str(uuid.uuid1())
    Clients[id] = Client
    return id

def GetClient (id):
    return Clients[id]

def DeleteClient (id):
    del Clients[id]