using System;
using Robomotion;
using Octokit;


namespace Github
{
    [Color("#003b57"), Icon(Icons.mdiGithubCircle), Inputs(1), Outputs(1)]
    [NodeName("Robomotion.Github.CreateRepository"), Title("Create Repository")]
    public class CreateRepository : Node
    {
        [Default("Message", "conn_id"), MessageScope, Title("Connection Id")]
        public InVariable<string> ConnectionId { get; set; } 

        [Default("Custom", ""), MessageScope, CustomScope, Title("Repository Name")]
        public InVariable<string> RepoName { get; set; }

        [Default(true), Title("Auto Init"), Option]
        public bool OptAutoInit { get; set; }

        [Default(true), Title("Private"), Option]
        public bool OptPrivate { get; set; }

        [Default("Message", "result"), MessageScope, Title("Result"), Output]
        public OutVariable<object> OutResult { get; set; }

       
        public override void OnCreate()
        {
        }

        public override void OnMessage(ref Context ctx)
        {
            GitHubClient github = null;
            string connId = ConnectionId.Get(ctx);
            if (string.IsNullOrEmpty(connId))
            {
                throw new Error("ErrInvalidArg", "Connection id can not be empty");
            }
            string repoName = RepoName.Get(ctx);
            if (string.IsNullOrEmpty(repoName)) {
                throw new Error("ErrInvalidArg", "Repository Name can not be empty");
            }
             Github.mutex.WaitOne();

            if (!Github.connections.TryGetValue(connId, out github))
            {
                Github.mutex.ReleaseMutex();
                throw new Error("ErrInvalidArg", "Connection id not found");
            }
            var repos = github.Repository.GetAllForCurrent().Result;
           
            Github.mutex.ReleaseMutex();
            
            var result = github.Repository.Create(new NewRepository(repoName) {AutoInit = OptAutoInit, Private = OptPrivate }).Result;
            OutResult.Set(ref ctx, result);

        }

        public override void OnClose()
        {

        }

    }
}