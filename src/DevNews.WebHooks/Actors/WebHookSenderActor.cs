using System;
using Akka.Actor;
using DevNews.Shared.Messages;
using DevNews.WebHooks.Application.Services;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.WebHooks.Actors
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