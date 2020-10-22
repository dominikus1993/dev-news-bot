using System;
using Akka.Actor;
using DevNews.WebHooks.Application.Services;
using Microsoft.Extensions.DependencyInjection;

namespace DevNews.WebHooks.Actors
{
    public class WebHookSenderActor : ReceiveActor
    {
        private IWebHookNotifier _webHookNotifier;
        public WebHookSenderActor(IServiceProvider sp)
        {
            _webHookNotifier = sp.GetService<IWebHookNotifier>();
        }
    }
}