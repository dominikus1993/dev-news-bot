using System;
using System.Collections.Generic;
using Akka.Actor;
using Akka.Event;
using DevNews.Akka.Routing;
using DevNews.Akka.Messages;
using DevNews.Application.Notifications.Services;
using DevNews.Core.Model;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Application.Notifications.Actors
{
    public partial class WebHookSenderActor
    {
        private class NotifyUser
        {
            public NotifyUser(IEnumerable<Article> article, INotifier notifier)
            {
                Articles = article;
                Notifier = notifier;
            }

            public IEnumerable<Article> Articles { get; }
            public INotifier Notifier { get;  }
            
        }
    }
    public partial class WebHookSenderActor : ReceiveActor
    {
        private readonly IEnumerable<INotifier> _notifiers;
        
        public WebHookSenderActor(IServiceProvider sp)
        {
            _notifiers = sp.GetServices<INotifier>();
            Ready();
        }

        private void Ready()
        {
            Receive<SendArticles>(msg =>
            {
                foreach (var notifier in _notifiers)
                {
                    Context.Self.Tell(new NotifyUser(msg.Articles, notifier));
                }
            });
            
            ReceiveAsync<NotifyUser>(async msg =>
            {
                await msg.Notifier.Notify(msg.Articles);
            });
        }
        
        
        public static ActorMetaData HackerNewsParserActorPath = ActorMetaDataModule.CreateTopLevel("notifier");

        public static Props Create(IServiceProvider sp) => Props.Create(() => new WebHookSenderActor(sp));
    }
}