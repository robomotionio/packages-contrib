using System.Collections.Generic;
using System.Threading;
using Octokit;
using System;

namespace Github
{

    static class Github
    {
        public static Mutex mutex = new Mutex();

        public static Dictionary<string, GitHubClient> connections = new Dictionary<string, GitHubClient>();

        static void Main(string[] args)
        {
            Robomotion.Main.Start(args);

        }
    }
}
