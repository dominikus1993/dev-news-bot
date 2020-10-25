using System.Threading.Tasks;
using Akka.Actor;
using DevNews.Application.HackerNews.Actors;

namespace DevNews.Application.HackerNews.UseCases
{
    public class ParseHackerNewsMainPageAndNotifyUsersUseCase
    {
        private IActorRef _hackerNewsActor;

        public ParseHackerNewsMainPageAndNotifyUsersUseCase(Framework.Akka.Actors actors)
        {
            _hackerNewsActor = actors.HackerNewsActor;
        }

        public ValueTask Execute()
        {
            _hackerNewsActor.Tell(HackerNewsParserActor.ParseNewHackerNewsArticlesAndNotifyUsers.Instance);
            return new ValueTask();
        }
    }
}