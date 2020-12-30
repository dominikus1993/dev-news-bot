using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading.Tasks;
using DevNews.Core.Abstractions;
using DevNews.Core.Model;
using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using Newtonsoft.Json;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams
{
    internal class MicrosoftTeamsWebHookNotifier : INotifier
    {
        public const string MicrosoftTeamsWebHookNotifierApi = nameof(MicrosoftTeamsWebHookNotifierApi);
        
        private readonly HttpClient _client;
        private readonly ITeamsMessageSerializer _teamsMessageSerializer;

        public MicrosoftTeamsWebHookNotifier(IHttpClientFactory clientFactory, ITeamsMessageSerializer teamsMessageSerializer)
        {
            _client = clientFactory.CreateClient(MicrosoftTeamsWebHookNotifierApi);
            _teamsMessageSerializer = teamsMessageSerializer;
        }

        public async Task Notify(IEnumerable<Article> articles)
        {
            var msg = CreateMicrosoftTeamsMessage(articles);
            using var request = new HttpRequestMessage(HttpMethod.Post, string.Empty);
            request.Headers.Accept.Add(MediaTypeWithQualityHeaderValue.Parse("application/json"));
            var json = await _teamsMessageSerializer.Serialize(msg);
            request.Content = new StringContent(json, Encoding.UTF8, "application/json");
            var response = await _client.SendAsync(request);
            response.EnsureSuccessStatusCode();
        }

        private static MicrosoftTeamsMessage CreateMicrosoftTeamsMessage(IEnumerable<Article> articles)
        {
            var text = new MicrosoftTeamsTextBlock("Witam serdecznie, oto nowe newsy");
            var links = articles.Select(article =>
                new MicrosoftTeamsAction(Guid.NewGuid().ToString(), article.Title, article.Link)).ToList();
            var action = new MicrosoftTeamsActionSet(Guid.NewGuid().ToString(), links);
            var content = new MicrosoftTeamsMessageContent(new IMicrosoftTeamsMessageContentBody[] {text, action});
            var att = new MicrosoftTeamsAttachment(content);
            return new MicrosoftTeamsMessage(new[] {att});
        }
    }
}