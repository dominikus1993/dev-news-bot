using System;
using Akka.Actor;
using DevNews.Application.HackerNews.Servies;
using DevNews.Shared.Routing;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.FSharp.Core;

namespace DevNews.Application.HackerNews.Actors
{
    public partial class HackerNewsParserActor
    {
        public class ParseNewHackerNewsArticles
        {
            public static ParseNewHackerNewsArticles Instance = new ParseNewHackerNewsArticles();
        }    
    }
    
    public partial class HackerNewsParserActor : ReceiveActor
    {
        private IHackerNewsParser _hackerNewsParser;
        public HackerNewsParserActor(IServiceProvider sp)
        {
            _hackerNewsParser = sp.GetService<IHackerNewsParser>();
            Ready();
        }

        private void Ready()
        {
            ReceiveAsync<ParseNewHackerNewsArticles>(async msg =>
            {
                var result = _hackerNewsParser.Parse();
            });
        }

        public static ActorMetaData HackerNewsParserActorPath = ActorMetaDataModule.CreateTopLevel("hacker_news");

        public static Props Create(IServiceProvider sp) => Props.Create(() => new HackerNewsParserActor(sp));

    }
}