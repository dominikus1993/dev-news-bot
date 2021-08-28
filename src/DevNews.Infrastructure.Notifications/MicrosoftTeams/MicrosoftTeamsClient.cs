using DevNews.Infrastructure.Notifications.MicrosoftTeams.Model;
using DevNews.Infrastructure.Notifications.MicrosoftTeams.Serialization;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading;
using System.Threading.Tasks;

namespace DevNews.Infrastructure.Notifications.MicrosoftTeams
{
    internal sealed class MicrosoftTeamsClient
    {
        private readonly HttpClient _client;
        private readonly ITeamsMessageSerializer _teamsMessageSerializer;

        public MicrosoftTeamsClient(HttpClient client, ITeamsMessageSerializer teamsMessageSerializer)
        {
            _client = client;
            _teamsMessageSerializer = teamsMessageSerializer;
        }
        public async Task Notify(MicrosoftTeamsMessage message, CancellationToken cancellationToken = default)
        {
            using var request = new HttpRequestMessage(HttpMethod.Post, string.Empty);
            request.Headers.Accept.Add(MediaTypeWithQualityHeaderValue.Parse("application/json"));
            request.Content = new StringContent(await _teamsMessageSerializer.Serialize(message), Encoding.UTF8, "application/json");
            using var response = await _client.SendAsync(request, HttpCompletionOption.ResponseHeadersRead, cancellationToken);
            response.EnsureSuccessStatusCode();
        }
    }

}
