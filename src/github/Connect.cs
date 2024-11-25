using System;
using Robomotion;
using Octokit;

namespace Github
{
    [Color("#003b57"), Icon(Icons.mdiGithubCircle), Inputs(1), Outputs(1)]
    [NodeName("Robomotion.Github.Connect"), Title("Connect")]
    public class Connect : Node
    {
        [Default("Custom", ""), MessageScope, CustomScope, Title("Application Name")]
        public InVariable<string> AppName { get; set; }

        [Category(CategoryAttribute.ECategory.Token),MessageScope, CustomScope, Title("Access Token"), Option]
        public Credential OptCredentials { get; set; }

        [Default("Message", "conn_id"), MessageScope, Title("Connection Id"), Output]
        public OutVariable<string> OutConnectionId { get; set; }

        public override void OnCreate()
        {
            
        }

        public override void OnMessage(ref Context ctx)
        {
            string appName = AppName.Get(ctx);
            if (string.IsNullOrEmpty(appName)) {
                throw new Error("ErrInvalidArg", "Application Name can not be empty");
            }
            var creds = OptCredentials.Get(ctx);
            var token = creds["value"];
            
            Guid guid = System.Guid.NewGuid();

            Github.mutex.WaitOne();

            var github = new GitHubClient(new ProductHeaderValue(appName));
            var tokenAuth = new Credentials(token.ToString());
            github.Credentials = tokenAuth;
            Github.connections.Add(guid.ToString(), github);

            Github.mutex.ReleaseMutex();

            OutConnectionId.Set(ref ctx, guid.ToString());
        }

        public override void OnClose()
        {

        }
    }
}
