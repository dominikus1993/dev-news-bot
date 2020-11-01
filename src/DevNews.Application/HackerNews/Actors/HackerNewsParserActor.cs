using System;
using System.Linq;
using Akka.Actor;
using Akka.Streams;
using Akka.Streams.Dsl;
using DevNews.Application.HackerNews.Servies;
using DevNews.Application.Notifications.Actors;
using DevNews.Core.HackerNews;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.FSharp.Core;
using DevNews.Akka.Routing;
using DevNews.Akka.Messages;
using DevNews.Core.Model;

namespace DevNews.Application.HackerNews.Actors
{
    public partial class HackerNewsParserActor : ReceiveActor
    {
        public delegate IActorRef HackerNewsActorProvider();

        public static ActorMetaData HackerNewsParserActorPath = ActorMetaDataModule.CreateTopLevel("hacker_news");

        public static Props Create(IServiceProvider sp) => Props.Create(() => new HackerNewsParserActor(sp));

        public class ParseNewHackerNewsArticlesAndNotifyUsers
        {
            public static ParseNewHackerNewsArticlesAndNotifyUsers Instance =
                new ParseNewHackerNewsArticlesAndNotifyUsers();
        }

        public class CompleteParsingArticles
        {
            public static CompleteParsingArticles Instance = new CompleteParsingArticles();
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
            var parser = Source.ActorRef<ParseNewHackerNewsArticlesAndNotifyUsers>(50, OverflowStrategy.DropNew)
                .SelectAsync(Environment.ProcessorCount,
                    async msg => { return await _hackerNewsParser.Parse().ToListAsync(); }).SelectAsync(1,
                    async articles =>
                    {
                        return await _hackerNewsRepository.Exists(articles.Select(x =>
                                new Article(x.Title, x.Link)))
                            .Where(x => !x.Item2)
                            .Select(x => x.Item1)
                            .ToListAsync();
                    })
                .Select(articles => new SendArticles(articles))
                .To(Sink.ActorRef<SendArticles>(Self, CompleteParsingArticles.Instance));
                
            Ready();
        }

        private void Ready()
        {
            ReceiveAsync<ParseNewHackerNewsArticlesAndNotifyUsers>(async msg => { });
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