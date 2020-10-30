using System;
using System.Linq;
using Akka.Actor;
using DevNews.Akka.Messages;
using DevNews.Akka.Types;
using DevNews.Application.HackerNews.Servies;
using DevNews.Application.Notifications.Actors;
using DevNews.Shared.Routing;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.FSharp.Core;

namespace DevNews.Application.HackerNews.Actors
{
    public partial class HackerNewsParserActor : ReceiveActor
    {
        public delegate IActorRef HackerNewsActorProvider();
        
        public static ActorMetaData HackerNewsParserActorPath = ActorMetaDataModule.CreateTopLevel("hacker_news");

        public static Props Create(IServiceProvider sp) => Props.Create(() => new HackerNewsParserActor(sp));
        
        public class ParseNewHackerNewsArticlesAndNotifyUsers
        {
            public static ParseNewHackerNewsArticlesAndNotifyUsers Instance = new ParseNewHackerNewsArticlesAndNotifyUsers();
        }
    }

    public partial class HackerNewsParserActor : ReceiveActor
    {
        private IHackerNewsParser _hackerNewsParser;
        private IActorRef _parser = ActorRefs.Nobody;

        public HackerNewsParserActor(IServiceProvider sp)
        {
            _hackerNewsParser = sp.GetService<IHackerNewsParser>();
            _parser = Source
            Ready();
        }

        private void Ready()
        {
            ReceiveAsync<ParseNewHackerNewsArticlesAndNotifyUsers>(async msg =>
            {
                Become();
            });
        }
        
        
        private void Working()
        {
            ReceiveAsync<ParseNewHackerNewsArticlesAndNotifyUsers>(async msg =>
            {
                var result = await _hackerNewsParser.Parse().Select(x => new Article(x.Title, x.Link)).ToListAsync();
                var notifier = Context.ActorSelection(WebHookSenderActor.HackerNewsParserActorPath.Path);
                notifier.Tell(new SendArticles(result));
            });
        }
    }
}