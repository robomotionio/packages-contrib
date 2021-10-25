using System;
using Robomotion;
using Octokit;

namespace Github
{
    [Color("#003b57"), Icon(Icons.mdiGithubCircle), Inputs(1), Outputs(1)]
    [NodeName("Robomotion.Github.Disconnect"), Title("Disconnect")]
    public class Disconnect : Node
    {


        [Default("Message", "conn_id"), MessageScope, Title("Connection Id"), Output]
        public InVariable<string> ConnectionId { get; set; }

        public override void OnCreate()
        {
            
        }

        public override void OnMessage(ref Context ctx)
        {
        GitHubClient github = null;

        string connId = ConnectionId.Get(ctx);
        if (string.IsNullOrEmpty(connId)) {
            throw new Error("ErrInvalidArg", "Connection id can not be empty");
        }

            Github.mutex.WaitOne();

        if (!Github.connections.TryGetValue(connId, out github))
        {
            Github.mutex.ReleaseMutex();
            throw new Error("ErrInvalidArg", "Connection id not found");
        }
        Github.connections.Remove(connId);
        
        }

        public override void OnClose()
        {

        }
    }
}