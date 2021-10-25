using System;
using Robomotion;
using Octokit;
using System.Collections.Generic;

namespace Github
{
    [Color("#003b57"), Icon(Icons.mdiGithubCircle), Inputs(1), Outputs(1)]
    [NodeName("Robomotion.Github.ListRepositories"), Title("List Repositories")]
    public class ListRepositories : Node
    {
        [Default("Message", "conn_id"), MessageScope, Title("Connection Id")]
        public InVariable<string> ConnectionId { get; set; } 


        [Default("Message", "result"), MessageScope, Title("Result"), Output]
        public OutVariable<object> OutResult { get; set; }

        public struct Repo {
        public long id;
        public string language;
        public string name;
        public string description;
        public DateTimeOffset createdat;

        }
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
             Github.mutex.WaitOne();

            if (!Github.connections.TryGetValue(connId, out github))
            {
                Github.mutex.ReleaseMutex();
                throw new Error("ErrInvalidArg", "Connection id not found");
            }
            var repos = github.Repository.GetAllForCurrent().Result;
            List<Repo> repositories = new List<Repo>();
            foreach(var item in repos){
               var repoAttributes = new Repo {id = item.Id, language= item.Language, name= item.Name, description= item.Description, createdat= item.CreatedAt };
               repositories.Add(repoAttributes);
            }
            Github.mutex.ReleaseMutex();

            OutResult.Set(ref ctx, repositories);

        }

        public override void OnClose()
        {

        }

    }
}