using System;
using Akka.Actor;
using DevNews.Akka.Messages;
using DevNews.Application.Notifications.Services;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.Application.Notifications.Actors
{
    public class WebHookSenderActor : ReceiveActor
    {
        private INotifier _notifier;
        public WebHookSenderActor(IServiceProvider sp)
        {
            _notifier = sp.GetService<INotifier>();
            Ready();
        }

        private void Ready()
        {
            ReceiveAsync<SendArticles>(async msg =>
            {
                await _notifier.Notify(msg.Articles);
            });
        }
    }
}