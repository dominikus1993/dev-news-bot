using System.Threading.Tasks;
using Akka.Actor;
using DevNews.Application.HackerNews.Actors;

namespace DevNews.Application.HackerNews.UseCases
{
    public class ParseHackerNewsMainPageAndNotifyUsersUseCase
    {
        private IActorRef _hackerNewsActor;

        public ParseHackerNewsMainPageAndNotifyUsersUseCase(Framework.Akka.SystemActors systemActors)
        {
            _hackerNewsActor = systemActors.HackerNewsActor;
        }

        public ValueTask Execute()
        {
            _hackerNewsActor.Tell(HackerNewsParserActor.ParseNewHackerNewsArticlesAndNotifyUsers.Instance);
            return new ValueTask();
        }
    }
}