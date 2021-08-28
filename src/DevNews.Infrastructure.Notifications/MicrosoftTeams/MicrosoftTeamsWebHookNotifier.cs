using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using DevNews.Infrastructure.Notifications.MicrosoftTeams.Serialization;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams
{
    internal class MicrosoftTeamsWebHookNotifier : INotifier
    {
        private readonly MicrosoftTeamsClient _client;
        public MicrosoftTeamsWebHookNotifier(MicrosoftTeamsClient client)
        {
            _client = client;
        }

        public async Task Notify(IEnumerable<Article> articles, CancellationToken cancellationToken = default)
        {
            var msg = MicrosoftTeamsMessage.From(articles);
            await _client.Notify(msg, cancellationToken);
        }
    }
}
