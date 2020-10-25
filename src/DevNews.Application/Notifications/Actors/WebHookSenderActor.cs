using System;
using Akka.Actor;
using Akka.Event;
using DevNews.Akka.Messages;
using DevNews.Shared.Routing;

namespace DevNews.Application.Notifications.Actors
{
    public class WebHookSenderActor : ReceiveActor
    {
        public WebHookSenderActor(IServiceProvider sp)
        {
            Ready();
        }

        private void Ready()
        {
            Receive<SendArticles>(msg =>
            {
                Context.GetLogger().Info("Message send {msg}", msg);
            });
        }
        
        
        public static ActorMetaData HackerNewsParserActorPath = ActorMetaDataModule.CreateTopLevel("notifier");

        public static Props Create(IServiceProvider sp) => Props.Create(() => new WebHookSenderActor(sp));
    }
}