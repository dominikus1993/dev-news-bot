using Akka.Actor;

namespace DevNews.Application.Framework.Akka
{
    public class Actors
    {
        public Actors(ActorSystem actorSystem, IActorRef hackerNewsActor, IActorRef webHookNotifierActor)
        {
            HackerNewsActor = hackerNewsActor;
            WebHookNotifierActor = webHookNotifierActor;
            ActorSystem = actorSystem;
        }

        public IActorRef HackerNewsActor { get; }
        public IActorRef WebHookNotifierActor { get; }
        public ActorSystem ActorSystem { get; }
    }
}