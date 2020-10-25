using System;
using Akka.Actor;
using DevNews.Application.HackerNews.Servies;

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
        public HackerNewsParserActor(IServiceProvider sp, IHackerNewsParser hackerNewsParser)
        {
            _hackerNewsParser = hackerNewsParser;
            Ready();
        }

        private void Ready()
        {
            ReceiveAsync<ParseNewHackerNewsArticles>(async msg =>
            {
                var result = _hackerNewsParser.Parse();
            });
        }
    }
}