using System;
using System.Linq;
using Akka.Actor;
using Akka.Streams;
using Akka.Streams.Dsl;
using DevNews.Akka.Messages;
using DevNews.Akka.Types;
using DevNews.Application.HackerNews.Servies;
using DevNews.Application.Notifications.Actors;
using DevNews.Domain.HackerNews;
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
        private IHackerNewsRepository _hackerNewsRepository;
        private IActorRef _parser = ActorRefs.Nobody;

        public HackerNewsParserActor(IServiceProvider sp)
        {
            _hackerNewsParser = sp.GetService<IHackerNewsParser>();
            _hackerNewsRepository = sp.GetService<IHackerNewsRepository>();
            _parser = Source.ActorRef<ParseNewHackerNewsArticlesAndNotifyUsers>(50, OverflowStrategy.DropNew)
                .SelectAsync(Environment.ProcessorCount, async msg =>
                {
                    return await _hackerNewsParser.Parse().ToListAsync();
                }).SelectAsync(1, async articles =>
                {
                    return await _hackerNewsRepository.Exists(articles.Select(x =>
                        new Domain.Model.Article(x.Title, x.Link))).ToListAsync();
                } ).To(Sink.ActorRef<>())
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